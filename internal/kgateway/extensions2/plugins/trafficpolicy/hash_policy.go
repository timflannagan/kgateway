package trafficpolicy

import (
	"sort"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	envoyroutev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

type HashPolicyIR struct {
	hashPolicies []*envoyroutev3.RouteAction_HashPolicy
}

func (h *HashPolicyIR) Equals(other *HashPolicyIR) bool {
	if h == nil && other == nil {
		return true
	}
	if h == nil || other == nil {
		return false
	}
	if len(h.hashPolicies) != len(other.hashPolicies) {
		return false
	}
	for i, policy := range h.hashPolicies {
		if !proto.Equal(policy, other.hashPolicies[i]) {
			return false
		}
	}
	return true
}

// Validate performs validation on the hash policy components
func (h *HashPolicyIR) Validate() error {
	if h == nil || h.hashPolicies == nil {
		return nil
	}
	for _, policy := range h.hashPolicies {
		if err := policy.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// MergeInto handles merging hash policies from p2 into p1
func (h *HashPolicyIR) MergeInto(
	p1, p2 *TrafficPolicy,
	p2Ref *pluginsdkir.AttachedPolicyRef,
	opts policy.MergeOptions,
	mergeOrigins pluginsdkir.MergeOrigins,
) {
	if !policy.IsMergeable(p1.spec.hashPolicies, p2.spec.hashPolicies, opts) {
		return
	}

	switch opts.Strategy {
	case policy.AugmentedDeepMerge, policy.OverridableDeepMerge:
		if p1.spec.hashPolicies != nil {
			return
		}
		fallthrough // can override p1 if it is unset

	case policy.AugmentedShallowMerge, policy.OverridableShallowMerge:
		p1.spec.hashPolicies = p2.spec.hashPolicies
		mergeOrigins.SetOne("hashPolicies", p2Ref)

	default:
		logger.Warn("unsupported merge strategy for hashPolicies policy", "strategy", opts.Strategy, "policy", p2Ref)
	}
}

func (h *HashPolicyIR) HashPolicies() []*envoyroutev3.RouteAction_HashPolicy {
	if h == nil {
		return nil
	}
	return h.hashPolicies
}

func hashPolicyForSpec(spec v1alpha1.TrafficPolicySpec, outSpec *trafficPolicySpecIr) {
	if len(spec.HashPolicies) == 0 {
		return
	}
	policies := make([]*envoyroutev3.RouteAction_HashPolicy, 0, len(spec.HashPolicies))
	for _, hashPolicy := range spec.HashPolicies {
		policy := &envoyroutev3.RouteAction_HashPolicy{}
		if hashPolicy.Terminal != nil {
			policy.Terminal = *hashPolicy.Terminal
		}
		switch {
		case hashPolicy.Header != nil:
			policy.PolicySpecifier = &envoyroutev3.RouteAction_HashPolicy_Header_{
				Header: &envoyroutev3.RouteAction_HashPolicy_Header{
					HeaderName: hashPolicy.Header.Name,
				},
			}
		case hashPolicy.Cookie != nil:
			policy.PolicySpecifier = &envoyroutev3.RouteAction_HashPolicy_Cookie_{
				Cookie: &envoyroutev3.RouteAction_HashPolicy_Cookie{
					Name: hashPolicy.Cookie.Name,
				},
			}
			if hashPolicy.Cookie.TTL != nil {
				policy.GetCookie().Ttl = durationpb.New(hashPolicy.Cookie.TTL.Duration)
			}
			if hashPolicy.Cookie.Path != nil {
				policy.GetCookie().Path = *hashPolicy.Cookie.Path
			}
			if hashPolicy.Cookie.Attributes != nil {
				// Get all attribute names and sort them for consistent ordering
				names := make([]string, 0, len(hashPolicy.Cookie.Attributes))
				for name := range hashPolicy.Cookie.Attributes {
					names = append(names, name)
				}
				sort.Strings(names)

				attributes := make([]*envoyroutev3.RouteAction_HashPolicy_CookieAttribute, 0, len(hashPolicy.Cookie.Attributes))
				for _, name := range names {
					attributes = append(attributes, &envoyroutev3.RouteAction_HashPolicy_CookieAttribute{
						Name:  name,
						Value: hashPolicy.Cookie.Attributes[name],
					})
				}
				policy.GetCookie().Attributes = attributes
			}
		case hashPolicy.SourceIP != nil:
			policy.PolicySpecifier = &envoyroutev3.RouteAction_HashPolicy_ConnectionProperties_{
				ConnectionProperties: &envoyroutev3.RouteAction_HashPolicy_ConnectionProperties{
					SourceIp: true,
				},
			}
		}
		policies = append(policies, policy)
	}
	outSpec.hashPolicies = &HashPolicyIR{
		hashPolicies: policies,
	}
}
