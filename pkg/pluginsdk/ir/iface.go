package ir

import (
	"context"
	"fmt"
	"time"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/plugins"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
)

var logger = logging.New("pluginsdk/ir")

type ListenerContext struct {
	Policy            PolicyIR
	PolicyAncestorRef gwv1.ParentReference
}

type RouteConfigContext struct {
	Policy            PolicyIR
	FilterChainName   string
	TypedFilterConfig TypedFilterConfigMap
}

type VirtualHostContext struct {
	Policy            PolicyIR
	FilterChainName   string
	TypedFilterConfig TypedFilterConfigMap
}

type TypedFilterConfigMap map[string]proto.Message

// AddTypedConfig SETS the config for a given key. // TODO: consider renaming to SetTypedConfig
func (r *TypedFilterConfigMap) AddTypedConfig(key string, v proto.Message) {
	if *r == nil {
		*r = make(TypedFilterConfigMap)
	}
	(*r)[key] = v
}

func (r *TypedFilterConfigMap) GetTypedConfig(key string) proto.Message {
	if r == nil || *r == nil {
		return nil
	}
	if v, ok := (*r)[key]; ok {
		return v
	}
	return nil
}

func (r *TypedFilterConfigMap) ToAnyMap() map[string]*anypb.Any {
	typedPerFilterConfigAny := map[string]*anypb.Any{}
	for k, v := range *r {
		config, err := utils.MessageToAny(v)
		if err != nil {
			logger.Error("unexpected marshalling error", "error", err)
			continue
		}
		typedPerFilterConfigAny[k] = config
	}
	return typedPerFilterConfigAny
}

type RouteBackendContext struct {
	FilterChainName string
	Backend         *BackendObjectIR
	// TypedFilterConfig will be output on the Route or WeightedCluster level after all plugins have run
	TypedFilterConfig TypedFilterConfigMap
}

type RouteContext struct {
	FilterChainName string
	Policy          PolicyIR
	In              HttpRouteRuleMatchIR
	// TypedFilterConfig will be output on the Route level after all plugins have run
	TypedFilterConfig TypedFilterConfigMap
}

type HcmContext struct {
	Policy PolicyIR
}

// ProxyTranslationPass represents a single translation pass for a gateway. It can hold state
// for the duration of the translation.
// Each of the functions here will be called in the order they appear in the interface.
type ProxyTranslationPass interface {
	//	Name() string
	// called 1 time for each listener
	ApplyListenerPlugin(
		ctx context.Context,
		pCtx *ListenerContext,
		out *envoy_config_listener_v3.Listener,
	)
	// called 1 time per filter chain after listeners and allows tweaking HCM settings.
	ApplyHCM(ctx context.Context,
		pCtx *HcmContext,
		out *envoy_hcm.HttpConnectionManager,
	)

	// called 1 time for all the routes in a filter chain. Use this to set default PerFilterConfig
	// No policy is provided here.
	ApplyRouteConfigPlugin(
		ctx context.Context,
		pCtx *RouteConfigContext,
		out *envoy_config_route_v3.RouteConfiguration,
	)
	ApplyVhostPlugin(
		ctx context.Context,
		pCtx *VirtualHostContext,
		out *envoy_config_route_v3.VirtualHost,
	)
	// no policy applied - this is called for every backend in a route.
	// For this to work the backend needs to register itself as a policy. TODO: rethink this.
	// Note: TypedFilterConfig should be applied in the pCtx and is shared between ApplyForRoute, ApplyForBackend
	// and ApplyForRouteBacken (do not apply on the output route directly)
	ApplyForBackend(
		ctx context.Context,
		pCtx *RouteBackendContext,
		in HttpBackend,
		out *envoy_config_route_v3.Route,
	)
	// Applies a policy attached to a specific Backend (via extensionRef on the BackendRef).
	// Note: TypedFilterConfig should be applied in the pCtx and is shared between ApplyForRoute, ApplyForBackend
	// and ApplyForRouteBackend
	ApplyForRouteBackend(
		ctx context.Context,
		policy PolicyIR,
		pCtx *RouteBackendContext,
	)
	// called once per route rule if SupportsPolicyMerge returns false, otherwise this is called only
	// once on the value returned by MergePolicies.
	// Applies policy for an HTTPRoute that has a policy attached via a targetRef.
	// The output configures the envoy_config_route_v3.Route
	// Note: TypedFilterConfig should be applied in the pCtx and is shared between ApplyForRoute, ApplyForBackend
	// and ApplyForRouteBacken (do not apply on the output route directly)
	ApplyForRoute(
		ctx context.Context,
		pCtx *RouteContext,
		out *envoy_config_route_v3.Route,
	)

	// TODO: Should this return an error?
	// called 1 time per filter-chain.
	// If a plugin emits new filters, they must be with a plugin unique name.
	// filters added to impact specific routes should be disabled on the listener level, so they don't impact other routes.
	HttpFilters(ctx context.Context, fc FilterChainCommon) ([]plugins.StagedHttpFilter, error)
	// TODO: Should this return an error?
	NetworkFilters(ctx context.Context) ([]plugins.StagedNetworkFilter, error)
	// called 1 time (per envoy proxy). replaces GeneratedResources and allows adding clusters to the envoy.
	ResourcesToAdd(ctx context.Context) Resources
}

type UnimplementedProxyTranslationPass struct{}

var _ ProxyTranslationPass = UnimplementedProxyTranslationPass{}

func (s UnimplementedProxyTranslationPass) ApplyListenerPlugin(ctx context.Context, pCtx *ListenerContext, out *envoy_config_listener_v3.Listener) {
}

func (s UnimplementedProxyTranslationPass) ApplyHCM(ctx context.Context, pCtx *HcmContext, out *envoy_hcm.HttpConnectionManager) {
}

func (s UnimplementedProxyTranslationPass) ApplyForBackend(ctx context.Context, pCtx *RouteBackendContext, in HttpBackend, out *envoy_config_route_v3.Route) {
}

func (s UnimplementedProxyTranslationPass) ApplyRouteConfigPlugin(ctx context.Context, pCtx *RouteConfigContext, out *envoy_config_route_v3.RouteConfiguration) {
}

func (s UnimplementedProxyTranslationPass) ApplyVhostPlugin(ctx context.Context, pCtx *VirtualHostContext, out *envoy_config_route_v3.VirtualHost) {
}

func (s UnimplementedProxyTranslationPass) ApplyForRoute(ctx context.Context, pCtx *RouteContext, out *envoy_config_route_v3.Route) {
}

func (s UnimplementedProxyTranslationPass) ApplyForRouteBackend(ctx context.Context, policy PolicyIR, pCtx *RouteBackendContext) {
}

func (s UnimplementedProxyTranslationPass) HttpFilters(ctx context.Context, fc FilterChainCommon) ([]plugins.StagedHttpFilter, error) {
	return nil, nil
}

func (s UnimplementedProxyTranslationPass) NetworkFilters(ctx context.Context) ([]plugins.StagedNetworkFilter, error) {
	return nil, nil
}

func (s UnimplementedProxyTranslationPass) ResourcesToAdd(ctx context.Context) Resources {
	return Resources{}
}

// Resources is the set of resources to add to the envoy config.
type Resources struct {
	// Clusters is the set of clusters to add to the envoy config.
	Clusters []*envoy_config_cluster_v3.Cluster
}

// GwTranslationCtx is the context for a gateway translation pass.
type GwTranslationCtx struct{}

// PolicyIR represents the IR state between a policy and a resource.
type PolicyIR interface {
	// CreationTime is the time the policy was created.
	// in case multiple policies attached to the same resource, we sort by policy creation time.
	CreationTime() time.Time
	// Equals returns true if the policy is the same as the input.
	Equals(in any) bool
}

// PolicyWrapper wraps a policy object and its IR.
type PolicyWrapper struct {
	// ObjectSource is the source of the policy object.
	ObjectSource `json:",inline"`
	// Policy is the policy object itself. TODO: we can probably remove this
	Policy metav1.Object

	// Errors processing the policy for status.
	// note: these errors are based on policy itself, regardless of whether it's attached to a resource.
	// Errors should be formatted for users, so do not include internal lib errors.
	// Instead use a well defined error such as ErrInvalidConfig
	Errors []error

	// PolicyIR is the IR of the policy objects. ideally with structural errors removed.
	// Opaque to us other than metadata.
	PolicyIR PolicyIR

	// TargetRefs is where to attach the policy. This usually comes from the policy CRD.
	TargetRefs []PolicyRef
}

// ResourceName returns the name of the resource the policy is attached to.
func (c PolicyWrapper) ResourceName() string {
	return c.ObjectSource.ResourceName()
}

// Equals returns true if the policy is the same as the input.
func (c PolicyWrapper) Equals(in PolicyWrapper) bool {
	if c.ObjectSource != in.ObjectSource {
		return false
	}

	return versionEquals(c.Policy, in.Policy) && c.PolicyIR.Equals(in.PolicyIR)
}

func versionEquals(a, b metav1.Object) bool {
	var versionEquals bool
	if a.GetGeneration() != 0 && b.GetGeneration() != 0 {
		versionEquals = a.GetGeneration() == b.GetGeneration()
	} else {
		versionEquals = a.GetResourceVersion() == b.GetResourceVersion()
	}
	return versionEquals && a.GetUID() == b.GetUID()
}

// ErrNotAttachable is returned when a policy is not attachable to an object.
var ErrNotAttachable = fmt.Errorf("policy is not attachable to this object")

// PolicyRun is the interface for a policy run.
type PolicyRun interface {
	// NewGatewayTranslationPass allocates state for single listener+rotue translation pass.
	NewGatewayTranslationPass(ctx context.Context, tctx GwTranslationCtx, reporter reports.Reporter) ProxyTranslationPass
	// ProcessBackend processes a backend for a single listener+rotue translation pass.
	ProcessBackend(ctx context.Context, in BackendObjectIR, out *envoy_config_cluster_v3.Cluster) error
}
