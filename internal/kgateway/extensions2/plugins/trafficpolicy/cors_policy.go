package trafficpolicy

import (
	corsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	envoy_wellknown "github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/proto"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

type CorsIR struct {
	// corsConfig is the envoy cors policy
	corsConfig *corsv3.CorsPolicy
}

var _ PolicySubIR = &CorsIR{}

func (c *CorsIR) Equals(other PolicySubIR) bool {
	otherCors, ok := other.(*CorsIR)
	if !ok {
		return false
	}
	if c == nil && otherCors == nil {
		return true
	}
	if c == nil || otherCors == nil {
		return false
	}
	return proto.Equal(c.corsConfig, otherCors.corsConfig)
}

// Validate performs validation on the CORS policy component
func (c *CorsIR) Validate() error {
	if c == nil || c.corsConfig == nil {
		return nil
	}
	return c.corsConfig.ValidateAll()
}

// MergeInto handles merging CORS policy from p2 into p1
func (c *CorsIR) MergeInto(
	p1, p2 *TrafficPolicy,
	p2Ref *pluginsdkir.AttachedPolicyRef,
	opts policy.MergeOptions,
	mergeOrigins pluginsdkir.MergeOrigins,
) {
	if !policy.IsMergeable(p1.spec.cors, p2.spec.cors, opts) {
		return
	}

	switch opts.Strategy {
	case policy.AugmentedDeepMerge, policy.OverridableDeepMerge:
		if p1.spec.cors != nil {
			return
		}
		fallthrough // can override p1 if it is unset

	case policy.AugmentedShallowMerge, policy.OverridableShallowMerge:
		p1.spec.cors = p2.spec.cors
		mergeOrigins.SetOne("cors", p2Ref)

	default:
		logger.Warn("unsupported merge strategy for cors policy", "strategy", opts.Strategy, "policy", p2Ref)
	}
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
