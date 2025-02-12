package upstream

import (
	"bytes"
	"context"
	"fmt"
	"maps"
	"time"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_lambda_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/aws_lambda/v3"
	"github.com/solo-io/go-utils/contextutils"
	"google.golang.org/protobuf/types/known/anypb"
	skubeclient "istio.io/istio/pkg/config/schema/kubeclient"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/plugins"
	"github.com/kgateway-dev/kgateway/v2/pkg/client/clientset/versioned"
)

var (
	ParameterGK = schema.GroupKind{
		Group: ParameterGroup,
		Kind:  ParameterKind,
	}
)

const (
	ExtensionName = "Upstream"
)

type upstreamDestination struct {
	FunctionName string
}

func (d *upstreamDestination) CreationTime() time.Time {
	return time.Time{}
}

func (d *upstreamDestination) Equals(in any) bool {
	d2, ok := in.(*upstreamDestination)
	if !ok {
		return false
	}
	return d.FunctionName == d2.FunctionName
}

type UpstreamIr struct {
	Region    string
	AccountId string
	AwsSecret *ir.Secret
}

func (u *UpstreamIr) data() map[string][]byte {
	if u.AwsSecret == nil {
		return nil
	}
	return u.AwsSecret.Data
}

func (u *UpstreamIr) Equals(other any) bool {
	otherUpstream, ok := other.(*UpstreamIr)
	if !ok {
		return false
	}
	return maps.EqualFunc(u.data(), otherUpstream.data(), func(a, b []byte) bool {
		return bytes.Equal(a, b)
	})
}

func registerTypes(ourCli versioned.Interface) {
	skubeclient.Register[*v1alpha1.Upstream](
		v1alpha1.UpstreamGVK.GroupVersion().WithResource("upstreams"),
		v1alpha1.UpstreamGVK,
		func(c skubeclient.ClientGetter, namespace string, o metav1.ListOptions) (runtime.Object, error) {
			return ourCli.GatewayV1alpha1().Upstreams(namespace).List(context.Background(), o)
		},
		func(c skubeclient.ClientGetter, namespace string, o metav1.ListOptions) (watch.Interface, error) {
			return ourCli.GatewayV1alpha1().Upstreams(namespace).Watch(context.Background(), o)
		},
	)
}

func NewPlugin(ctx context.Context, commoncol *common.CommonCollections) extensionsplug.Plugin {
	registerTypes(commoncol.OurClient)

	col := krt.WrapClient(kclient.New[*v1alpha1.Upstream](commoncol.Client), commoncol.KrtOpts.ToOptions("Upstreams")...)

	// FIXME: hack.
	l := contextutils.LoggerFrom(ctx)

	gk := v1alpha1.UpstreamGVK.GroupKind()
	translate := buildTranslateFunc(commoncol.Secrets)

	ucol := krt.NewCollection(col, func(krtctx krt.HandlerContext, i *v1alpha1.Upstream) *ir.Upstream {
		objIr, err := translate(krtctx, i)
		if err != nil {
			l.Warnf("failed to translate upstream: %v", err)
			return nil
		}
		return &ir.Upstream{
			ObjectSource: ir.ObjectSource{
				Kind:      gk.Kind,
				Group:     gk.Group,
				Namespace: i.GetNamespace(),
				Name:      i.GetName(),
			},
			GvPrefix:          "upstream",
			CanonicalHostname: hostname(i),
			Obj:               i,
			ObjIr:             objIr,
		}
	})

	endpoints := krt.NewCollection(col, func(krtctx krt.HandlerContext, i *v1alpha1.Upstream) *ir.EndpointsForUpstream {
		return processEndpoints(i)
	})

	return extensionsplug.Plugin{
		ContributesUpstreams: map[schema.GroupKind]extensionsplug.UpstreamPlugin{
			gk: {
				UpstreamInit: ir.UpstreamInit{
					InitUpstream: func(ctx context.Context, in ir.Upstream, out *envoy_config_cluster_v3.Cluster) {
						if err := processUpstream(ctx, in, out); err != nil {
							l.Warnf("failed to process upstream: %v", err)
						}
					},
				},
				Endpoints: endpoints,
				Upstreams: ucol,
			},
		},
		ContributesPolicies: map[schema.GroupKind]extensionsplug.PolicyPlugin{
			ParameterGK: {
				Name:                      ExtensionName,
				NewGatewayTranslationPass: newPlug,
				PoliciesFetch: func(n, ns string) ir.PolicyIR {
					// The "Parameters" kind is a synthetic (i.e. virtual) kind used to relay
					// information without having to have a real policy object. In this case
					// the lambda function name is the name of the parameter.
					return &upstreamDestination{
						FunctionName: n,
					}
				},
			},
		},
	}
}

func buildTranslateFunc(secrets *krtcollections.SecretIndex) func(krtctx krt.HandlerContext, i *v1alpha1.Upstream) (*UpstreamIr, error) {
	return func(krtctx krt.HandlerContext, i *v1alpha1.Upstream) (*UpstreamIr, error) {
		if i.Spec.Type != v1alpha1.UpstreamTypeAws {
			return nil, nil
		}
		ns := i.GetNamespace()
		secretRef := gwv1.SecretObjectReference{
			Name:      gwv1.ObjectName(i.Spec.Aws.SecretRef.Name),
			Namespace: ptr.To(gwv1.Namespace(ns)),
		}
		secret, err := secrets.GetSecret(krtctx, krtcollections.From{GroupKind: v1alpha1.UpstreamGVK.GroupKind(), Namespace: ns}, secretRef)
		if err != nil {
			return nil, fmt.Errorf("failed to get the %s secret in the %s namespace: %w", i.Spec.Aws.SecretRef.Name, ns, err)
		}
		if secret == nil {
			return nil, fmt.Errorf("internal error:secret not found")
		}
		var ir UpstreamIr
		ir.AwsSecret = secret
		ir.Region = i.Spec.Aws.Region
		ir.AccountId = i.Spec.Aws.AccountId
		return &ir, nil
	}
}

func processUpstream(ctx context.Context, in ir.Upstream, out *envoy_config_cluster_v3.Cluster) error {
	up, ok := in.Obj.(*v1alpha1.Upstream)
	if !ok {
		return fmt.Errorf("internal error: upstream is not a v1alpha1.Upstream")
	}
	ir, ok := in.ObjIr.(*UpstreamIr)
	if !ok {
		return fmt.Errorf("internal error: upstreamIr is not a UpstreamIr")
	}

	upType := up.Spec.Type
	switch upType {
	case v1alpha1.UpstreamTypeStatic:
		processStatic(ctx, up.Spec.Static, out)
	case v1alpha1.UpstreamTypeAws:
		if err := processAws(ctx, up.Spec.Aws, ir, out); err != nil {
			return fmt.Errorf("failed to process aws upstream: %w", err)
		}
	default:
		return fmt.Errorf("unknown upstream type: %s", upType)
	}
	return nil
}

func hostname(in *v1alpha1.Upstream) string {
	if in.Spec.Static == nil {
		return ""
	}
	if len(in.Spec.Static.Hosts) > 0 {
		return string(in.Spec.Static.Hosts[0].Host)
	}
	return ""
}

func processEndpoints(up *v1alpha1.Upstream) *ir.EndpointsForUpstream {
	spec := up.Spec
	switch {
	case spec.Static != nil:
		return processEndpointsStatic(spec.Static)
	case spec.Aws != nil:
		return processEndpointsAws(spec.Aws)
	}
	return nil
}

type plugin2 struct {
	needFilter map[string]bool
	derived    staticSecretDerivation
	arn        string
}

func newPlug(ctx context.Context, tctx ir.GwTranslationCtx) ir.ProxyTranslationPass {
	return &plugin2{}
}

func (p *plugin2) Name() string { return ExtensionName }

// called 1 time for each listener
func (p *plugin2) ApplyListenerPlugin(ctx context.Context, pCtx *ir.ListenerContext, out *envoy_config_listener_v3.Listener) {
}

func (p *plugin2) ApplyVhostPlugin(ctx context.Context, pCtx *ir.VirtualHostContext, out *envoy_config_route_v3.VirtualHost) {
}

// called 0 or more times
func (p *plugin2) ApplyForRoute(ctx context.Context, pCtx *ir.RouteContext, outputRoute *envoy_config_route_v3.Route) error {
	return nil
}

func (p *plugin2) ApplyForRouteBackend(
	ctx context.Context,
	policy ir.PolicyIR,
	pCtx *ir.RouteBackendContext,
) error {
	pol, ok := policy.(*upstreamDestination)
	if !ok {
		return fmt.Errorf("internal error: policy is not a upstreamDestination")
	}
	if p.needFilter == nil {
		p.needFilter = make(map[string]bool)
	}
	p.needFilter[pCtx.FilterChainName] = true

	u, ok := pCtx.Upstream.ObjIr.(*UpstreamIr)
	if !ok {
		return fmt.Errorf("internal error: upstream is not a Upstream")
	}
	derived, err := deriveStaticSecret(u.AwsSecret)
	if err != nil {
		return fmt.Errorf("failed to derive static secret: %w", err)
	}
	if derived.access == "" || derived.secret == "" {
		return fmt.Errorf("required credentials not found in secret")
	}
	p.derived = derived
	p.arn = fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", u.Region, u.AccountId, pol.FunctionName)

	lambdaRouteFunc := &envoy_lambda_v3.PerRouteConfig{
		InvokeConfig: &envoy_lambda_v3.Config{
			Arn: p.arn,
			Credentials: &envoy_lambda_v3.Credentials{
				AccessKeyId:     derived.access,
				SecretAccessKey: derived.secret,
				SessionToken:    derived.session,
			},
		},
	}
	lambdaRouteFuncAny, err := anypb.New(lambdaRouteFunc)
	if err != nil {
		return err
	}
	pCtx.AddTypedConfig(FilterName, lambdaRouteFuncAny)

	return nil
}

// HttpFilters ...
// called 1 time per listener
// if a plugin emits new filters, they must be with a plugin unique name.
// any filter returned from route config must be disabled, so it doesnt impact other routes.
func (p *plugin2) HttpFilters(ctx context.Context, fc ir.FilterChainCommon) ([]plugins.StagedHttpFilter, error) {
	if !p.needFilter[fc.FilterChainName] {
		return nil, nil
	}

	l := contextutils.LoggerFrom(ctx)
	l.Infof("building lambda filter for arn: %s", p.arn)

	filterConfig := &envoy_lambda_v3.Config{
		Arn:            p.arn,
		InvocationMode: envoy_lambda_v3.Config_SYNCHRONOUS,
		Credentials: &envoy_lambda_v3.Credentials{
			AccessKeyId:     p.derived.access,
			SecretAccessKey: p.derived.secret,
			SessionToken:    p.derived.session,
		},
	}

	pluginStage := plugins.DuringStage(plugins.OutAuthStage)
	lambdaF, err := plugins.NewStagedFilter(FilterName, filterConfig, pluginStage)
	if err != nil {
		return nil, fmt.Errorf("failed to create lambda filter: %w", err)
	}

	return []plugins.StagedHttpFilter{
		lambdaF,
	}, nil
}

func (p *plugin2) UpstreamHttpFilters(ctx context.Context) ([]plugins.StagedUpstreamHttpFilter, error) {
	return nil, nil
}

func (p *plugin2) NetworkFilters(ctx context.Context) ([]plugins.StagedNetworkFilter, error) {
	return nil, nil
}

// called 1 time (per envoy proxy). replaces GeneratedResources
func (p *plugin2) ResourcesToAdd(ctx context.Context) ir.Resources {
	return ir.Resources{}
}
