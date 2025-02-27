package httplistenerpolicy

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"time"

	envoyaccesslog "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/solo-io/go-utils/contextutils"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/config/schema/kubeclient"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	extensionplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/plugins"
	"github.com/kgateway-dev/kgateway/v2/pkg/client/clientset/versioned"
)

type httpListenerPolicy struct {
	ct        time.Time
	compress  bool
	accessLog []*envoyaccesslog.AccessLog
}

func (d *httpListenerPolicy) CreationTime() time.Time {
	return d.ct
}

func (d *httpListenerPolicy) Equals(in any) bool {
	d2, ok := in.(*httpListenerPolicy)
	if !ok {
		return false
	}

	// Check the TargetRef and Compress fields
	if d.compress != d2.compress {
		return false
	}

	// Check the AccessLog slice
	if !slices.EqualFunc(d.accessLog, d2.accessLog, func(log *envoyaccesslog.AccessLog, log2 *envoyaccesslog.AccessLog) bool {
		return proto.Equal(log, log2)
	}) {
		return false
	}

	return true
}

func registerTypes(ourCli versioned.Interface) {
	kubeclient.Register[*v1alpha1.HTTPListenerPolicy](
		v1alpha1.HTTPListenerPolicyGVK.GroupVersion().WithResource("httplistenerpolicies"),
		v1alpha1.HTTPListenerPolicyGVK,
		func(c kubeclient.ClientGetter, namespace string, o metav1.ListOptions) (runtime.Object, error) {
			return ourCli.GatewayV1alpha1().HTTPListenerPolicies(namespace).List(context.Background(), o)
		},
		func(c kubeclient.ClientGetter, namespace string, o metav1.ListOptions) (watch.Interface, error) {
			return ourCli.GatewayV1alpha1().HTTPListenerPolicies(namespace).Watch(context.Background(), o)
		},
	)
}

func NewPlugin(ctx context.Context, commoncol *common.CommonCollections) extensionplug.Plugin {
	// need upstreams for backend refs
	registerTypes(commoncol.OurClient)
	col := krt.WrapClient(kclient.New[*v1alpha1.HTTPListenerPolicy](commoncol.Client), commoncol.KrtOpts.ToOptions("HTTPListenerPolicy")...)
	gk := v1alpha1.HTTPListenerPolicyGVK.GroupKind()

	translateFn := translateFn(ctx, commoncol.Backends)

	policyCol := krt.NewCollection(col, func(krtctx krt.HandlerContext, i *v1alpha1.HTTPListenerPolicy) *ir.PolicyWrapper {
		objSrc := ir.ObjectSource{
			Group:     gk.Group,
			Kind:      gk.Kind,
			Namespace: i.Namespace,
			Name:      i.Name,
		}

		polIR, err := translateFn(krtctx, i, objSrc)

		return &ir.PolicyWrapper{
			ObjectSource: objSrc,
			Policy:       i,
			PolicyIR:     polIR,
			TargetRefs:   convert(i.Spec.TargetRef),
			Errors:       []error{err},
		}
	})

	return extensionplug.Plugin{
		ContributesPolicies: map[schema.GroupKind]extensionsplug.PolicyPlugin{
			gk: {
				Name:                      "httplistenerpolicies",
				NewGatewayTranslationPass: NewGatewayTranslationPass,
				Policies:                  policyCol,
			},
		},
		ExtraHasSynced: func() bool {
			return commoncol.HasSynced()
		},
	}
}

func translateFn(
	ctx context.Context,
	backends *krtcollections.BackendIndex,
) func(krtctx krt.HandlerContext, i *v1alpha1.HTTPListenerPolicy, src ir.ObjectSource) (*httpListenerPolicy, error) {
	return func(krtctx krt.HandlerContext, i *v1alpha1.HTTPListenerPolicy, src ir.ObjectSource) (*httpListenerPolicy, error) {
		accessLog, err := convertAccessLogConfig(ctx, i, backends, krtctx, src)
		if err != nil {
			return nil, err
		}

		return &httpListenerPolicy{
			ct:        i.CreationTimestamp.Time,
			compress:  i.Spec.Compress,
			accessLog: accessLog,
		}, nil
	}
}

func convert(targetRef v1alpha1.LocalPolicyTargetReference) []ir.PolicyTargetRef {
	return []ir.PolicyTargetRef{{
		Kind:  string(targetRef.Kind),
		Name:  string(targetRef.Name),
		Group: string(targetRef.Group),
	}}
}

type httpListenerPolicyPluginGwPass struct {
}

func (p *httpListenerPolicyPluginGwPass) ApplyListenerPlugin(ctx context.Context, pCtx *ir.ListenerContext, out *envoy_config_listener_v3.Listener) {
	// no op
}

func NewGatewayTranslationPass(ctx context.Context, tctx ir.GwTranslationCtx) ir.ProxyTranslationPass {
	return &httpListenerPolicyPluginGwPass{}
}

func (p *httpListenerPolicyPluginGwPass) Name() string {
	return "httplistenerpolicies"
}

func (p *httpListenerPolicyPluginGwPass) ApplyHCM(
	ctx context.Context,
	pCtx *ir.HcmContext,
	out *envoy_hcm.HttpConnectionManager,
) error {
	policy, ok := pCtx.Policy.(*httpListenerPolicy)
	if !ok {
		return fmt.Errorf("internal error: expected httplistener policy, got %T", pCtx.Policy)
	}

	al := policy.accessLog
	if len(al) == 0 {
		contextutils.LoggerFrom(ctx).Info("tim -- no access log configs were cached for policy")
		return nil
	}
	sort.Slice(al, func(i, j int) bool {
		return al[i].Name < al[j].Name
	})

	// translate access logging configuration
	contextutils.LoggerFrom(ctx).Infof("tim -- plugin translated %d access log configs", len(al))
	out.AccessLog = append(out.GetAccessLog(), al...)
	return nil
}

func (p *httpListenerPolicyPluginGwPass) ApplyVhostPlugin(ctx context.Context, pCtx *ir.VirtualHostContext, out *envoy_config_route_v3.VirtualHost) {
}

// called 0 or more times
func (p *httpListenerPolicyPluginGwPass) ApplyForRoute(ctx context.Context, pCtx *ir.RouteContext, outputRoute *envoy_config_route_v3.Route) error {
	return nil
}

func (p *httpListenerPolicyPluginGwPass) ApplyForRouteBackend(
	ctx context.Context,
	policy ir.PolicyIR,
	pCtx *ir.RouteBackendContext,
) error {
	return nil
}

// called 1 time per listener
// if a plugin emits new filters, they must be with a plugin unique name.
// any filter returned from listener config must be disabled, so it doesnt impact other listeners.
func (p *httpListenerPolicyPluginGwPass) HttpFilters(ctx context.Context, fcc ir.FilterChainCommon) ([]plugins.StagedHttpFilter, error) {
	return nil, nil
}

func (p *httpListenerPolicyPluginGwPass) UpstreamHttpFilters(ctx context.Context) ([]plugins.StagedUpstreamHttpFilter, error) {
	return nil, nil
}

func (p *httpListenerPolicyPluginGwPass) NetworkFilters(ctx context.Context) ([]plugins.StagedNetworkFilter, error) {
	return nil, nil
}

// called 1 time (per envoy proxy). replaces GeneratedResources
func (p *httpListenerPolicyPluginGwPass) ResourcesToAdd(ctx context.Context) ir.Resources {
	return ir.Resources{}
}
