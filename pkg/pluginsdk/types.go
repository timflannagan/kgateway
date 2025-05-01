// pkg/pluginsdk/types.go
//
// Primary SDK surface: Plugin registration types, the ProxyTranslationPass
// interface, no-op helper, and opaque aliases for concrete IR objects.

package pluginsdk

import (
	"context"

	envoy_listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"google.golang.org/protobuf/proto"
)

// -----------------------------------------------------------------------------
// Top-level plugin bundle
// -----------------------------------------------------------------------------

type Plugin struct {
	ContributesPolicies  ContributesPolicies
	ContributesBackends  ContributesBackends
	ContributesResources ResourcesPlugin
}

type ContributesPolicies map[GroupKind]PolicyPlugin
type ContributesBackends map[GroupKind]BackendPlugin

type ResourcesPlugin interface {
	ResourcesToAdd(context.Context) Resources
}

type BackendPlugin struct{}                 // placeholder for future use
type GroupKind struct{ Group, Kind string } // coarse alias of schema.GroupKind

// -----------------------------------------------------------------------------
// PolicyPlugin contract
// -----------------------------------------------------------------------------

type PolicyPlugin struct {
	Name                      string
	Policies                  any // e.g. a krt.Collection handle
	NewGatewayTranslationPass func(context.Context, GwTranslationCtx) ProxyTranslationPass
	MergePolicies             func([]PolicyIR) PolicyIR
	SupportsPolicyMerge       bool
}

// -----------------------------------------------------------------------------
// Core interface every plugin pass must satisfy
// -----------------------------------------------------------------------------

type ProxyTranslationPass interface {
	ApplyListenerPlugin(context.Context, *ListenerContext, *envoy_listener.Listener)
	ApplyHCM(context.Context, *HcmContext, *envoy_hcm.HttpConnectionManager) error
	ApplyRouteConfigPlugin(context.Context, *RouteConfigContext, *envoy_route.RouteConfiguration)
	ApplyVhostPlugin(context.Context, *VirtualHostContext, *envoy_route.VirtualHost)
	ApplyForBackend(context.Context, *RouteBackendContext, HttpBackend, *envoy_route.Route) error
	ApplyForRouteBackend(context.Context, PolicyIR, *RouteBackendContext) error
	ApplyForRoute(context.Context, *RouteContext, *envoy_route.Route) error
	HttpFilters(context.Context, FilterChainCommon) ([]StagedHttpFilter, error)
	NetworkFilters(context.Context) ([]StagedNetworkFilter, error)
	ResourcesToAdd(context.Context) Resources
}

// -----------------------------------------------------------------------------
// Convenience no-op implementation
// -----------------------------------------------------------------------------

type UnimplementedProxyTranslationPass struct{}

var _ ProxyTranslationPass = UnimplementedProxyTranslationPass{}

func (UnimplementedProxyTranslationPass) ApplyListenerPlugin(context.Context, *ListenerContext, *envoy_listener.Listener) {
}
func (UnimplementedProxyTranslationPass) ApplyHCM(context.Context, *HcmContext, *envoy_hcm.HttpConnectionManager) error {
	return nil
}
func (UnimplementedProxyTranslationPass) ApplyRouteConfigPlugin(context.Context, *RouteConfigContext, *envoy_route.RouteConfiguration) {
}
func (UnimplementedProxyTranslationPass) ApplyVhostPlugin(context.Context, *VirtualHostContext, *envoy_route.VirtualHost) {
}
func (UnimplementedProxyTranslationPass) ApplyForBackend(context.Context, *RouteBackendContext, HttpBackend, *envoy_route.Route) error {
	return nil
}
func (UnimplementedProxyTranslationPass) ApplyForRouteBackend(context.Context, PolicyIR, *RouteBackendContext) error {
	return nil
}
func (UnimplementedProxyTranslationPass) ApplyForRoute(context.Context, *RouteContext, *envoy_route.Route) error {
	return nil
}
func (UnimplementedProxyTranslationPass) HttpFilters(context.Context, FilterChainCommon) ([]StagedHttpFilter, error) {
	return nil, nil
}
func (UnimplementedProxyTranslationPass) NetworkFilters(context.Context) ([]StagedNetworkFilter, error) {
	return nil, nil
}
func (UnimplementedProxyTranslationPass) ResourcesToAdd(context.Context) Resources {
	return Resources{}
}

// -----------------------------------------------------------------------------
// Opaque aliases to keep SDK independent of internal packages
// -----------------------------------------------------------------------------

type (
	BackendObjectIR      = any
	HttpRouteRuleMatchIR = any
	FilterChainCommon    = any
	StagedHttpFilter     = any
	StagedNetworkFilter  = any
	Resources            = struct{ Clusters []proto.Message }
	PolicyIR             interface {
		CreationTime() any
		Equals(other any) bool
	}
	GwTranslationCtx = struct{}
)

// ***** STABLE CONTRACT END *****
