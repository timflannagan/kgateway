// pkg/xds/bootstrap/config.go
package bootstrap

import (
	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extensions_filters_network_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoywellknown "github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"google.golang.org/protobuf/proto"
)

// ConfigBuilder helps construct a partial bootstrap config for validation.
type ConfigBuilder struct {
	configs ir.TypedFilterConfigMap
}

// New creates a new ConfigBuilder
func New() *ConfigBuilder {
	return &ConfigBuilder{
		configs: make(ir.TypedFilterConfigMap),
	}
}

// AddFilterConfig adds a filter configuration to the builder.
func (b *ConfigBuilder) AddFilterConfig(name string, config proto.Message) {
	b.configs.AddTypedConfig(name, config)
}

// Build creates a partial bootstrap config suitable for validation.
func (b *ConfigBuilder) Build() (*envoy_config_bootstrap_v3.Bootstrap, error) {
	hcmAny, err := utils.MessageToAny(&envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
		StatPrefix: "placeholder",
		RouteSpecifier: &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{
			RouteConfig: &envoy_config_route_v3.RouteConfiguration{
				VirtualHosts: []*envoy_config_route_v3.VirtualHost{{
					Name:                 "placeholder_vhost",
					Domains:              []string{"*"},
					TypedPerFilterConfig: b.configs.ToAnyMap(),
				}},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &envoy_config_bootstrap_v3.Bootstrap{
		Node: &envoy_config_core_v3.Node{
			Id:      "validation-node-id",
			Cluster: "validation-cluster",
		},
		StaticResources: &envoy_config_bootstrap_v3.Bootstrap_StaticResources{
			Listeners: []*envoy_config_listener_v3.Listener{{
				Name: "placeholder_listener",
				Address: &envoy_config_core_v3.Address{
					Address: &envoy_config_core_v3.Address_SocketAddress{
						SocketAddress: &envoy_config_core_v3.SocketAddress{
							Address:       "0.0.0.0",
							PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{PortValue: 8081},
						},
					},
				},
				FilterChains: []*envoy_config_listener_v3.FilterChain{{
					Name: "placeholder_filter_chain",
					Filters: []*envoy_config_listener_v3.Filter{{
						Name: envoywellknown.HTTPConnectionManager,
						ConfigType: &envoy_config_listener_v3.Filter_TypedConfig{
							TypedConfig: hcmAny,
						},
					}},
				}},
			}},
		},
	}, nil
}
