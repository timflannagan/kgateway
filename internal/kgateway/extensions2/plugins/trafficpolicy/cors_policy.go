package trafficpolicy

import (
	corsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	envoy_wellknown "github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/proto"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
)

type CorsIR struct {
	// corsConfig is the envoy cors policy
	corsConfig *corsv3.CorsPolicy
}

func (c *CorsIR) Equals(other *CorsIR) bool {
	if c == nil && other == nil {
		return true
	}
	if c == nil || other == nil {
		return false
	}
	return proto.Equal(c.corsConfig, other.corsConfig)
}

// Validate performs validation on the CORS policy component
func (c *CorsIR) Validate() error {
	if c == nil || c.corsConfig == nil {
		return nil
	}
	return c.corsConfig.Validate()
}

// corsForSpec translates the cors spec into an envoy cors policy and stores it in the traffic policy IR
func corsForSpec(in *v1alpha1.TrafficPolicy, out *trafficPolicySpecIr) error {
	if in.Spec.Cors == nil {
		return nil
	}
	out.cors = &CorsIR{
		corsConfig: utils.ToEnvoyCorsPolicy(in.Spec.Cors.HTTPCORSFilter),
	}
	return nil
}

func (p *trafficPolicyPluginGwPass) handleCors(fcn string, pCtxTypedFilterConfig *ir.TypedFilterConfigMap, cors *CorsIR) {
	if cors == nil || cors.corsConfig == nil {
		return
	}

	// Adds the CorsPolicy to the typed_per_filter_config.
	// Also requires Cors http_filter to be added to the filter chain.
	pCtxTypedFilterConfig.AddTypedConfig(envoy_wellknown.CORS, cors.corsConfig)

	// Add a filter to the chain. When having a cors policy for a route we need to also have a
	// globally cors http filter in the chain otherwise it will be ignored.
	if p.corsInChain == nil {
		p.corsInChain = make(map[string]*corsv3.Cors)
	}
	if _, ok := p.corsInChain[fcn]; !ok {
		p.corsInChain[fcn] = &corsv3.Cors{}
	}
}
