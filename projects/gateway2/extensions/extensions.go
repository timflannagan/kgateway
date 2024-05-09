package extensions

import (
	"context"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/translator"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	controllerruntime "sigs.k8s.io/controller-runtime"
	apiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// K8sGatewayExtensions is responsible for providing implementations for translation utilities
// which have Enterprise variants.
type K8sGatewayExtensions interface {
	// GetPluginRegistry exposes the plugins supported by this implementation.
	GetPluginRegistry() registry.PluginRegistry

	// GetTranslator allows an extnsion to provide custom translation for
	// different gateway classes.
	GetTranslator(*apiv1.Gateway) translator.K8sGwTranslator
}

// K8sGatewayExtensionsFactoryParameters contains the parameters required to start Gloo K8s Gateway Extensions (including Translator Plugins)
type K8sGatewayExtensionsFactoryParameters struct {
	Mgr                     controllerruntime.Manager
	AuthConfigClient        v1.AuthConfigClient
	RouteOptionClient       gatewayv1.RouteOptionClient
	VirtualHostOptionClient gatewayv1.VirtualHostOptionClient
	StatusReporter          reporter.StatusReporter
	KickXds                 func(ctx context.Context)
}

// K8sGatewayExtensionsFactory returns an extensions.K8sGatewayExtensions
type K8sGatewayExtensionsFactory func(
	ctx context.Context,
	params K8sGatewayExtensionsFactoryParameters,
) (K8sGatewayExtensions, error)

// NewK8sGatewayExtensions returns the Open Source implementation of K8sGatewayExtensions
func NewK8sGatewayExtensions(
	_ context.Context,
	params K8sGatewayExtensionsFactoryParameters,
) (K8sGatewayExtensions, error) {
	queries := query.NewData(
		params.Mgr.GetClient(),
		params.Mgr.GetScheme(),
	)
	plugins := registry.BuildPlugins(
		queries,
		params.Mgr.GetClient(),
		params.RouteOptionClient,
		params.VirtualHostOptionClient,
		params.StatusReporter,
	)
	pluginRegistry := registry.NewPluginRegistry(plugins)

	return &k8sGatewayExtensions{
		mgr:                     params.Mgr,
		routeOptionClient:       params.RouteOptionClient,
		virtualHostOptionClient: params.VirtualHostOptionClient,
		statusReporter:          params.StatusReporter,
		translator:              translator.NewTranslator(queries, pluginRegistry),
		pluginRegistry:          pluginRegistry,
	}, nil
}

type k8sGatewayExtensions struct {
	mgr                     controllerruntime.Manager
	routeOptionClient       gatewayv1.RouteOptionClient
	virtualHostOptionClient gatewayv1.VirtualHostOptionClient
	statusReporter          reporter.StatusReporter
	translator              translator.K8sGwTranslator
	pluginRegistry          registry.PluginRegistry
}

func (e *k8sGatewayExtensions) GetTranslator(_ *apiv1.Gateway) translator.K8sGwTranslator {
	return e.translator
}

func (e *k8sGatewayExtensions) GetPluginRegistry() registry.PluginRegistry {
	return e.pluginRegistry
}
