package trafficpolicy

import (
	"math"

	bufferv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/buffer/v3"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

const bufferFilterName = "envoy.filters.http.buffer"

type BufferIR struct {
	maxRequestBytes uint32
}

func (b *BufferIR) Equals(other *BufferIR) bool {
	if b == nil && other == nil {
		return true
	}
	if b == nil || other == nil {
		return false
	}

	return b.maxRequestBytes == other.maxRequestBytes
}

// Validate performs validation on the buffer component
// Note: buffer validation is not needed as it's a single uint32 field
func (b *BufferIR) Validate() error {
	return nil
}

// MergeInto handles merging buffer policy from p2 into p1
func (b *BufferIR) MergeInto(
	p1, p2 *TrafficPolicy,
	p2Ref *pluginsdkir.AttachedPolicyRef,
	opts policy.MergeOptions,
	mergeOrigins pluginsdkir.MergeOrigins,
) {
	if !policy.IsMergeable(p1.spec.buffer, p2.spec.buffer, opts) {
		return
	}

	switch opts.Strategy {
	case policy.AugmentedDeepMerge, policy.OverridableDeepMerge:
		if p1.spec.buffer != nil {
			return
		}
		fallthrough // can override p1 if it is unset

	case policy.AugmentedShallowMerge, policy.OverridableShallowMerge:
		p1.spec.buffer = p2.spec.buffer
		mergeOrigins.SetOne("buffer", p2Ref)

	default:
		logger.Warn("unsupported merge strategy for buffer policy", "strategy", opts.Strategy, "policy", p2Ref)
	}
}

// bufferForSpec translates the buffer spec into an envoy buffer policy and stores it in the traffic policy IR
func bufferForSpec(spec v1alpha1.TrafficPolicySpec, out *trafficPolicySpecIr) {
	if spec.Buffer == nil {
		return
	}

	out.buffer = &BufferIR{
		maxRequestBytes: uint32(spec.Buffer.MaxRequestSize.Value()),
	}
}

func (p *trafficPolicyPluginGwPass) handleBuffer(fcn string, pCtxTypedFilterConfig *ir.TypedFilterConfigMap, buffer *BufferIR) {
	if buffer == nil {
		return
	}

	// Add buffer configuration to the typed_per_filter_config for route-level override
	bufferPerRoute := &bufferv3.BufferPerRoute{
		Override: &bufferv3.BufferPerRoute_Buffer{
			Buffer: &bufferv3.Buffer{
				MaxRequestBytes: &wrapperspb.UInt32Value{Value: buffer.maxRequestBytes},
			},
		},
	}
	pCtxTypedFilterConfig.AddTypedConfig(bufferFilterName, bufferPerRoute)

	// Add a filter to the chain. When having a buffer policy for a route we need to also have a
	// globally disabled buffer filter in the chain otherwise it will be ignored.
	if p.bufferInChain == nil {
		p.bufferInChain = make(map[string]*bufferv3.Buffer)
	}
	if _, ok := p.bufferInChain[fcn]; !ok {
		p.bufferInChain[fcn] = &bufferv3.Buffer{
			MaxRequestBytes: &wrapperspb.UInt32Value{Value: math.MaxUint32},
		}
	}
}

// need to add disabled buffer to the filter chain
// enable on route
