package trafficpolicy

import (
	"testing"
	"time"

	envoyroutev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestHashPolicyIREquals(t *testing.T) {
	// Helper to create simple hash policy configurations for testing
	createSimpleHashPolicies := func(headerName string) []*envoyroutev3.RouteAction_HashPolicy {
		return []*envoyroutev3.RouteAction_HashPolicy{
			{
				PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Header_{
					Header: &envoyroutev3.RouteAction_HashPolicy_Header{
						HeaderName: headerName,
					},
				},
				Terminal: false,
			},
		}
	}

	createHashPoliciesWithTerminal := func(terminal bool) []*envoyroutev3.RouteAction_HashPolicy {
		return []*envoyroutev3.RouteAction_HashPolicy{
			{
				PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Header_{
					Header: &envoyroutev3.RouteAction_HashPolicy_Header{
						HeaderName: "x-user-id",
					},
				},
				Terminal: terminal,
			},
		}
	}

	tests := []struct {
		name     string
		hash1    *HashPolicyIR
		hash2    *HashPolicyIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			hash1:    nil,
			hash2:    nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			hash1:    nil,
			hash2:    &HashPolicyIR{hashPolicies: createSimpleHashPolicies("x-user-id")},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			hash1:    &HashPolicyIR{hashPolicies: createSimpleHashPolicies("x-user-id")},
			hash2:    nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			hash1:    &HashPolicyIR{hashPolicies: createSimpleHashPolicies("x-user-id")},
			hash2:    &HashPolicyIR{hashPolicies: createSimpleHashPolicies("x-user-id")},
			expected: true,
		},
		{
			name:     "different header names are not equal",
			hash1:    &HashPolicyIR{hashPolicies: createSimpleHashPolicies("x-user-id")},
			hash2:    &HashPolicyIR{hashPolicies: createSimpleHashPolicies("x-session-id")},
			expected: false,
		},
		{
			name:     "different terminal settings are not equal",
			hash1:    &HashPolicyIR{hashPolicies: createHashPoliciesWithTerminal(true)},
			hash2:    &HashPolicyIR{hashPolicies: createHashPoliciesWithTerminal(false)},
			expected: false,
		},
		{
			name:     "same terminal settings are equal",
			hash1:    &HashPolicyIR{hashPolicies: createHashPoliciesWithTerminal(true)},
			hash2:    &HashPolicyIR{hashPolicies: createHashPoliciesWithTerminal(true)},
			expected: true,
		},
		{
			name:     "nil hash policy fields are equal",
			hash1:    &HashPolicyIR{hashPolicies: nil},
			hash2:    &HashPolicyIR{hashPolicies: nil},
			expected: true,
		},
		{
			name:     "nil vs non-nil hash policy fields are not equal",
			hash1:    &HashPolicyIR{hashPolicies: nil},
			hash2:    &HashPolicyIR{hashPolicies: createSimpleHashPolicies("x-test")},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hash1.Equals(tt.hash2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.hash2.Equals(tt.hash1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		hash := &HashPolicyIR{hashPolicies: createSimpleHashPolicies("x-test")}
		assert.True(t, hash.Equals(hash), "hash should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameHash := func() *HashPolicyIR {
			return &HashPolicyIR{hashPolicies: createHashPoliciesWithTerminal(false)}
		}

		a := createSameHash()
		b := createSameHash()
		c := createSameHash()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestHashPolicyIRValidate(t *testing.T) {
	// Note: Hash policy validation is permissive at the protobuf level.
	tests := []struct {
		name string
		hash *HashPolicyIR
	}{
		{
			name: "nil hash policy is valid",
			hash: nil,
		},
		{
			name: "hash policy with nil policies is valid",
			hash: &HashPolicyIR{
				hashPolicies: nil,
			},
		},
		{
			name: "valid hash policy passes validation",
			hash: &HashPolicyIR{
				hashPolicies: []*envoyroutev3.RouteAction_HashPolicy{
					{
						PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Header_{
							Header: &envoyroutev3.RouteAction_HashPolicy_Header{
								HeaderName: "x-user-id",
							},
						},
						Terminal: false,
					},
				},
			},
		},
		{
			name: "empty hash policies list is valid",
			hash: &HashPolicyIR{
				hashPolicies: []*envoyroutev3.RouteAction_HashPolicy{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.hash.Validate()
			assert.NoError(t, err, "Hash policy validation should never fail at protobuf level")
		})
	}
}

func TestHashPolicyIRMergeInto(t *testing.T) {
	// Helper to create hash policy configurations
	createHashPolicyConfig := func(headerName string, terminal bool) *HashPolicyIR {
		return &HashPolicyIR{
			hashPolicies: []*envoyroutev3.RouteAction_HashPolicy{
				{
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Header_{
						Header: &envoyroutev3.RouteAction_HashPolicy_Header{
							HeaderName: headerName,
						},
					},
					Terminal: terminal,
				},
			},
		}
	}

	tests := []struct {
		name              string
		p1Hash            *HashPolicyIR
		p2Hash            *HashPolicyIR
		strategy          policy.MergeOptions
		expectOriginSet   bool
		expectPolicyCount *int
	}{
		{
			name:              "shallow merge with nil p1",
			p1Hash:            nil,
			p2Hash:            createHashPolicyConfig("x-user-id", true),
			strategy:          policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet:   true,
			expectPolicyCount: ptr.To(1),
		},
		{
			name:              "deep merge preserves p1 when set",
			p1Hash:            createHashPolicyConfig("x-session-id", false),
			p2Hash:            createHashPolicyConfig("x-user-id", true),
			strategy:          policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet:   false,
			expectPolicyCount: ptr.To(1),
		},
		{
			name:              "overridable shallow merge overwrites p1",
			p1Hash:            createHashPolicyConfig("x-session-id", false),
			p2Hash:            createHashPolicyConfig("x-user-id", true),
			strategy:          policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet:   true,
			expectPolicyCount: ptr.To(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{hashPolicies: tt.p1Hash}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{hashPolicies: tt.p2Hash}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			if tt.p2Hash != nil {
				tt.p2Hash.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectPolicyCount != nil {
				assert.NotNil(t, p1.spec.hashPolicies, "p1.hashPolicies should not be nil")
				assert.NotNil(t, p1.spec.hashPolicies.hashPolicies, "p1.hashPolicies.hashPolicies should not be nil")
				assert.Equal(t, *tt.expectPolicyCount, len(p1.spec.hashPolicies.hashPolicies))
			}

			if tt.expectOriginSet {
				origins := mergeOrigins.Get("hashPolicies")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("hashPolicies")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}

func TestHashPolicyForSpec(t *testing.T) {
	tests := []struct {
		name     string
		spec     v1alpha1.TrafficPolicySpec
		expected []*envoyroutev3.RouteAction_HashPolicy
	}{
		{
			name: "nil hash policies",
			spec: v1alpha1.TrafficPolicySpec{
				HashPolicies: nil,
			},
			expected: nil,
		},
		{
			name: "empty hash policies",
			spec: v1alpha1.TrafficPolicySpec{
				HashPolicies: []*v1alpha1.HashPolicy{},
			},
			expected: nil,
		},
		{
			name: "header hash policy",
			spec: v1alpha1.TrafficPolicySpec{
				HashPolicies: []*v1alpha1.HashPolicy{
					{
						Header: &v1alpha1.Header{
							Name: "x-user-id",
						},
						Terminal: ptr.To(true),
					},
				},
			},
			expected: []*envoyroutev3.RouteAction_HashPolicy{
				{
					Terminal: true,
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Header_{
						Header: &envoyroutev3.RouteAction_HashPolicy_Header{
							HeaderName: "x-user-id",
						},
					},
				},
			},
		},
		{
			name: "cookie hash policy without TTL and path",
			spec: v1alpha1.TrafficPolicySpec{
				HashPolicies: []*v1alpha1.HashPolicy{
					{
						Cookie: &v1alpha1.Cookie{
							Name: "session-id",
						},
					},
				},
			},
			expected: []*envoyroutev3.RouteAction_HashPolicy{
				{
					Terminal: false,
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Cookie_{
						Cookie: &envoyroutev3.RouteAction_HashPolicy_Cookie{
							Name: "session-id",
						},
					},
				},
			},
		},
		{
			name: "cookie hash policy with TTL and path",
			spec: v1alpha1.TrafficPolicySpec{
				HashPolicies: []*v1alpha1.HashPolicy{
					{
						Cookie: &v1alpha1.Cookie{
							Name: "session-id",
							TTL: &metav1.Duration{
								Duration: 30 * time.Minute,
							},
							Path: ptr.To("/api"),
							Attributes: map[string]string{
								"domain": "example.com",
								"secure": "true",
							},
						},
						Terminal: ptr.To(true),
					},
				},
			},
			expected: []*envoyroutev3.RouteAction_HashPolicy{
				{
					Terminal: true,
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Cookie_{
						Cookie: &envoyroutev3.RouteAction_HashPolicy_Cookie{
							Name: "session-id",
							Ttl:  durationpb.New(30 * time.Minute),
							Path: "/api",
							Attributes: []*envoyroutev3.RouteAction_HashPolicy_CookieAttribute{
								{
									Name:  "domain",
									Value: "example.com",
								},
								{
									Name:  "secure",
									Value: "true",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "source IP hash policy",
			spec: v1alpha1.TrafficPolicySpec{
				HashPolicies: []*v1alpha1.HashPolicy{
					{
						SourceIP: &v1alpha1.SourceIP{},
						Terminal: ptr.To(false),
					},
				},
			},
			expected: []*envoyroutev3.RouteAction_HashPolicy{
				{
					Terminal: false,
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_ConnectionProperties_{
						ConnectionProperties: &envoyroutev3.RouteAction_HashPolicy_ConnectionProperties{
							SourceIp: true,
						},
					},
				},
			},
		},
		{
			name: "multiple hash policies",
			spec: v1alpha1.TrafficPolicySpec{
				HashPolicies: []*v1alpha1.HashPolicy{
					{
						Header: &v1alpha1.Header{
							Name: "x-user-id",
						},
						Terminal: ptr.To(true),
					},
					{
						Cookie: &v1alpha1.Cookie{
							Name: "session-id",
							TTL: &metav1.Duration{
								Duration: 1 * time.Hour,
							},
						},
					},
					{
						SourceIP: &v1alpha1.SourceIP{},
					},
				},
			},
			expected: []*envoyroutev3.RouteAction_HashPolicy{
				{
					Terminal: true,
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Header_{
						Header: &envoyroutev3.RouteAction_HashPolicy_Header{
							HeaderName: "x-user-id",
						},
					},
				},
				{
					Terminal: false,
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Cookie_{
						Cookie: &envoyroutev3.RouteAction_HashPolicy_Cookie{
							Name: "session-id",
							Ttl:  durationpb.New(1 * time.Hour),
						},
					},
				},
				{
					Terminal: false,
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_ConnectionProperties_{
						ConnectionProperties: &envoyroutev3.RouteAction_HashPolicy_ConnectionProperties{
							SourceIp: true,
						},
					},
				},
			},
		},
		{
			name: "cookie hash policy with nil TTL",
			spec: v1alpha1.TrafficPolicySpec{
				HashPolicies: []*v1alpha1.HashPolicy{
					{
						Cookie: &v1alpha1.Cookie{
							Name: "session-id",
							TTL:  nil,
							Path: ptr.To("/api"),
						},
						Terminal: ptr.To(false),
					},
				},
			},
			expected: []*envoyroutev3.RouteAction_HashPolicy{
				{
					Terminal: false,
					PolicySpecifier: &envoyroutev3.RouteAction_HashPolicy_Cookie_{
						Cookie: &envoyroutev3.RouteAction_HashPolicy_Cookie{
							Name: "session-id",
							Path: "/api",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outSpec := &trafficPolicySpecIr{}

			hashPolicyForSpec(tt.spec, outSpec)

			var actual []*envoyroutev3.RouteAction_HashPolicy
			if outSpec.hashPolicies != nil {
				actual = outSpec.hashPolicies.HashPolicies()
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}
