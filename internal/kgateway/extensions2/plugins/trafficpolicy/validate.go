package trafficpolicy

import (
	"context"

	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extensions_filters_network_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_extensions_transport_sockets_tls_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	envoywellknown "github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/validator"
)

// TODO: Does this belong on the builder?
// TODO: Break this up into smaller functions. Just for POC purposes right now. Would
// be nice to re-use the p.handle* functions, but they're under a different type.
func (p *TrafficPolicy) Validate(ctx context.Context, v validator.Validator) error {
	typedFilterConfig := make(ir.TypedFilterConfigMap)
	if p.spec.transform != nil {
		typedFilterConfig.AddTypedConfig(transformationFilterNamePrefix, p.spec.transform)
	}
	if p.spec.rustformation != nil {
		typedFilterConfig.AddTypedConfig(rustformationFilterNamePrefix, p.spec.rustformation)
	}
	if p.spec.localRateLimit != nil {
		typedFilterConfig.AddTypedConfig(localRateLimitFilterNamePrefix, p.spec.localRateLimit)
	}
	if p.spec.ExtProc != nil && p.spec.ExtProc.provider != nil && p.spec.ExtProc.provider.ExtProc != nil {
		typedFilterConfig.AddTypedConfig(extProcFilterName(p.spec.ExtProc.provider.ResourceName()), p.spec.ExtProc.provider.ExtProc)
	}
	if p.spec.extAuth != nil && p.spec.extAuth.provider != nil && p.spec.extAuth.provider.ExtAuth != nil {
		typedFilterConfig.AddTypedConfig(extAuthFilterName(p.spec.extAuth.provider.ResourceName()), p.spec.extAuth.provider.ExtAuth)
	}
	if p.spec.rateLimit != nil && p.spec.rateLimit.provider != nil && p.spec.rateLimit.provider.RateLimit != nil {
		typedFilterConfig.AddTypedConfig(getRateLimitFilterName(p.spec.rateLimit.provider.ResourceName()), p.spec.rateLimit.provider.RateLimit)
	}

	hcmAny, err := utils.MessageToAny(&envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
		StatPrefix: "placeholder",
		RouteSpecifier: &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{
			RouteConfig: &envoy_config_route_v3.RouteConfiguration{
				VirtualHosts: []*envoy_config_route_v3.VirtualHost{{
					Name:                 "placeholder_vhost",
					Domains:              []string{"*"},
					TypedPerFilterConfig: typedFilterConfig.ToAnyMap(),
				}},
			},
		},
	})
	if err != nil {
		return err
	}

	bootstrap, err := buildFakeBootstrap(&envoyResources{
		listeners: []*envoy_config_listener_v3.Listener{&envoy_config_listener_v3.Listener{
			Name: "placeholder_listener",
			Address: &envoy_config_core_v3.Address{
				Address: &envoy_config_core_v3.Address_SocketAddress{SocketAddress: &envoy_config_core_v3.SocketAddress{
					Address:       "0.0.0.0",
					PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{PortValue: 8081},
				}},
			},
			FilterChains: []*envoy_config_listener_v3.FilterChain{{
				Name: "placeholder_filter_chain",
				Filters: []*envoy_config_listener_v3.Filter{
					{
						Name: envoywellknown.HTTPConnectionManager,
						ConfigType: &envoy_config_listener_v3.Filter_TypedConfig{
							TypedConfig: hcmAny,
						},
					},
				},
			}},
		}},
	})
	if err != nil {
		return err
	}

	data, err := protojson.Marshal(bootstrap)
	if err != nil {
		return err
	}
	if err := v.Validate(ctx, string(data)); err != nil {
		return err
	}

	return nil
}

type envoyResources struct {
	listeners []*envoy_config_listener_v3.Listener
	clusters  []*envoy_config_cluster_v3.Cluster
	secrets   []*envoy_extensions_transport_sockets_tls_v3.Secret
}

func buildFakeBootstrap(resources *envoyResources) (*envoy_config_bootstrap_v3.Bootstrap, error) {
	bootstrap := &envoy_config_bootstrap_v3.Bootstrap{
		Node: &envoy_config_core_v3.Node{
			Id:      "validation-node-id",
			Cluster: "validation-cluster",
		},
		StaticResources: &envoy_config_bootstrap_v3.Bootstrap_StaticResources{
			Listeners: resources.listeners,
			Clusters:  resources.clusters,
			Secrets:   resources.secrets,
		},
	}
	return bootstrap, nil
}
