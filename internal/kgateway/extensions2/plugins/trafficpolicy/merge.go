package trafficpolicy

import (
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

// MergeTrafficPolicies merges two traffic policies together, with the second policy taking precedence where applicable.
// The mergeOrigins parameter is used to track merge origins and should be populated with information
// about the origin policy reference for each merged field.
func MergeTrafficPolicies(
	p1, p2 *TrafficPolicy,
	p2Ref *ir.AttachedPolicyRef,
	mergeOpts policy.MergeOptions,
	mergeOrigins pluginsdkir.MergeOrigins,
) {
	if p1 == nil || p2 == nil {
		return
	}
	// Call MergeInto method on each policy sub-IR type
	if p2.spec.ai != nil {
		p2.spec.ai.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.extProc != nil {
		p2.spec.extProc.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.transformation != nil {
		p2.spec.transformation.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.rustformation != nil {
		p2.spec.rustformation.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.extAuth != nil {
		p2.spec.extAuth.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.localRateLimit != nil {
		p2.spec.localRateLimit.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.rateLimit != nil {
		p2.spec.rateLimit.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.cors != nil {
		p2.spec.cors.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.csrf != nil {
		p2.spec.csrf.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.buffer != nil {
		p2.spec.buffer.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.autoHostRewrite != nil {
		p2.spec.autoHostRewrite.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
	if p2.spec.hashPolicies != nil {
		p2.spec.hashPolicies.MergeInto(p1, p2, p2Ref, mergeOpts, mergeOrigins)
	}
}
