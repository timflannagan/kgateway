package trafficpolicy

import (
	"maps"

	apiannotations "github.com/kgateway-dev/kgateway/v2/api/annotations"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

// mergePolicies merges the given policy ordered from high to low priority (both hierarchically
// and within the same hierarchy) based on the constraints defined per PolicyAtt.
//
// It iterates policies in reverse order (low to high) to ensure higher priority policies can
// always use an OverridableMerge strategy to override lower priority ones. Iterating policies
// in the given priority order (high to low) requires more complex merging for delegated chains
// because policies anywhere in the chain may enable policy overrides for their children but we
// still need to ensure these children cannot override any policies set by their ancestors that
// are not marked as overridable, i.e., (r1,p1)-delegate->(r2,p2)-delegate->(r3,p3) where
// r=route p=policy needs to ensure p3 does not override p1 (assuming p1 does not enable overrides)
// even if p2 allows overrides. This is easier to guarantee by using an OverridableMerge strategy
// by merging in higher priority policies with different HierarchicalPriority.
func mergePolicies(policies []ir.PolicyAtt) ir.PolicyAtt {
	var out ir.PolicyAtt
	if len(policies) == 0 {
		return out
	}
	_, ok := policies[0].PolicyIr.(*TrafficPolicy)
	// ignore unknown types
	if !ok {
		return out
	}

	// base policy to merge into has an empty PolicyIr so it can always be merged into
	out = ir.PolicyAtt{
		GroupKind:    policies[0].GroupKind,
		PolicyRef:    policies[0].PolicyRef,
		MergeOrigins: map[string]*ir.AttachedPolicyRef{},
		PolicyIr:     &TrafficPolicy{},
	}
	merged := out.PolicyIr.(*TrafficPolicy)

	for i := len(policies) - 1; i >= 0; i-- {
		mergeOpts := policy.MergeOptions{
			Strategy: policy.OverridableMerge,
		}
		// If merging a policy lower in the hierarchy with a policy higher in the hierarchy AND
		// the policy higher in the hierarchy enables policy overrides, use an AugmentedMerge strategy
		// to preserve existing fields set by lower levels.
		// NOTE: the HierarchicalPriority check is necessary to prevent enabling override behavior among
		// policies in the same hierarchy, e.g., ExtensionRef vs TargetRef policy attached to the same route, as
		// DelegationInheritedPolicyPriorityPreferChild strictly applies to parent->child policy inheritance and is not applicable
		// outside delegated policy inheritance.
		if out.HierarchicalPriority < policies[i].HierarchicalPriority && policies[i].DelegationInheritedPolicyPriority == apiannotations.DelegationInheritedPolicyPriorityPreferChild {
			mergeOpts.Strategy = policy.AugmentedMerge
		}

		p2 := policies[i].PolicyIr.(*TrafficPolicy)
		p2Ref := policies[i].PolicyRef

		mergeOrigins := MergeTrafficPolicies(merged, p2, p2Ref, mergeOpts)
		maps.Copy(out.MergeOrigins, mergeOrigins)
		out.HierarchicalPriority = policies[i].HierarchicalPriority
		out.Errors = append(out.Errors, policies[i].Errors...)
	}

	return out
}

// MergeTrafficPolicies merges two TrafficPolicy IRs, returning a map that contains information
// about the origin policy reference for each merged field.
func MergeTrafficPolicies(
	p1, p2 *TrafficPolicy,
	p2Ref *ir.AttachedPolicyRef,
	mergeOpts policy.MergeOptions,
) map[string]*ir.AttachedPolicyRef {
	if p1 == nil || p2 == nil {
		return nil
	}
	mergeOrigins := make(map[string]*ir.AttachedPolicyRef)
	if policy.IsMergeable(p1.spec.AI, p2.spec.AI, mergeOpts) {
		p1.spec.AI = p2.spec.AI
		mergeOrigins["ai"] = p2Ref
	}
	if policy.IsMergeable(p1.spec.ExtProc, p2.spec.ExtProc, mergeOpts) {
		p1.spec.ExtProc = p2.spec.ExtProc
		mergeOrigins["extProc"] = p2Ref
	}
	if policy.IsMergeable(p1.spec.transform, p2.spec.transform, mergeOpts) {
		p1.spec.transform = p2.spec.transform
		mergeOrigins["transformation"] = p2Ref
	}
	if policy.IsMergeable(p1.spec.rustformation, p2.spec.rustformation, mergeOpts) {
		p1.spec.rustformation = p2.spec.rustformation
		mergeOrigins["rustformation"] = p2Ref
	}
	if policy.IsMergeable(p1.spec.extAuth, p2.spec.extAuth, mergeOpts) {
		p1.spec.extAuth = p2.spec.extAuth
		mergeOrigins["extAuth"] = p2Ref
	}
	if policy.IsMergeable(p1.spec.localRateLimit, p2.spec.localRateLimit, mergeOpts) {
		p1.spec.localRateLimit = p2.spec.localRateLimit
		mergeOrigins["rateLimit"] = p2Ref
	}
	// Handle global rate limit merging
	if policy.IsMergeable(p1.spec.rateLimit, p2.spec.rateLimit, mergeOpts) {
		p1.spec.rateLimit = p2.spec.rateLimit
		mergeOrigins["rateLimit"] = p2Ref
	}
	// Handle cors merging
	if policy.IsMergeable(p1.spec.cors, p2.spec.cors, mergeOpts) {
		p1.spec.cors = p2.spec.cors
		mergeOrigins["cors"] = p2Ref
	}

	// Handle CSRF policy merging
	if policy.IsMergeable(p1.spec.csrf, p2.spec.csrf, mergeOpts) {
		p1.spec.csrf = p2.spec.csrf
		mergeOrigins["csrf"] = p2Ref
	}

	// Handle buffer policy merging
	if policy.IsMergeable(p1.spec.buffer, p2.spec.buffer, mergeOpts) {
		p1.spec.buffer = p2.spec.buffer
		mergeOrigins["buffer"] = p2Ref
	}

	return mergeOrigins
}
