package trafficpolicy

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

type AutoHostRewriteIR struct {
	// TODO(tim): Is this the right way to do this?
	autoHostRewrite *wrapperspb.BoolValue
}

// Equals compares two AutoHostRewriteIRs for equality.
func (a *AutoHostRewriteIR) Equals(other *AutoHostRewriteIR) bool {
	if a == nil && other == nil {
		return true
	}
	if a == nil || other == nil {
		return false
	}
	return proto.Equal(a.autoHostRewrite, other.autoHostRewrite)
}

// Validate performs validation on the auto host rewrite component
func (a *AutoHostRewriteIR) Validate() error {
	return nil
}

// MergeInto handles merging auto host rewrite policy from p2 into p1
func (a *AutoHostRewriteIR) MergeInto(
	p1, p2 *TrafficPolicy,
	p2Ref *pluginsdkir.AttachedPolicyRef,
	opts policy.MergeOptions,
	mergeOrigins pluginsdkir.MergeOrigins,
) {
	if !policy.IsMergeable(p1.spec.autoHostRewrite, p2.spec.autoHostRewrite, opts) {
		return
	}

	switch opts.Strategy {
	case policy.AugmentedDeepMerge, policy.OverridableDeepMerge:
		if p1.spec.autoHostRewrite != nil {
			return
		}
		fallthrough // can override p1 if it is unset

	case policy.AugmentedShallowMerge, policy.OverridableShallowMerge:
		p1.spec.autoHostRewrite = p2.spec.autoHostRewrite
		mergeOrigins.SetOne("autoHostRewrite", p2Ref)

	default:
		logger.Warn("unsupported merge strategy for AutoHostRewrite policy", "strategy", opts.Strategy, "policy", p2Ref)
	}
}

func (a *AutoHostRewriteIR) AutoHostRewrite() *wrapperspb.BoolValue {
	if a == nil {
		return nil
	}
	return a.autoHostRewrite
}

// autoHostRewriteForSpec translates the auto host rewrite spec into an envoy auto host rewrite policy and stores it in the traffic policy IR
func autoHostRewriteForSpec(spec v1alpha1.TrafficPolicySpec, out *trafficPolicySpecIr) {
	if spec.AutoHostRewrite == nil {
		return
	}
	out.autoHostRewrite = &AutoHostRewriteIR{
		autoHostRewrite: wrapperspb.Bool(*spec.AutoHostRewrite),
	}
}
