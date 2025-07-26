package trafficpolicy

import (
	"testing"

	corsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	envoy_type_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/utils/ptr"

	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestCorsIREquals(t *testing.T) {
	createSimpleCors := func(allowOrigin string) *corsv3.CorsPolicy {
		return &corsv3.CorsPolicy{
			AllowOriginStringMatch: []*envoy_type_matcher_v3.StringMatcher{
				{
					MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
						Exact: allowOrigin,
					},
				},
			},
			AllowCredentials: wrapperspb.Bool(true),
		}
	}
	createCorsWith := func(allowCredentials bool) *corsv3.CorsPolicy {
		return &corsv3.CorsPolicy{
			AllowCredentials: wrapperspb.Bool(allowCredentials),
		}
	}

	tests := []struct {
		name     string
		cors1    *CorsIR
		cors2    *CorsIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			cors1:    nil,
			cors2:    nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			cors1:    nil,
			cors2:    &CorsIR{corsConfig: createSimpleCors("*")},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			cors1:    &CorsIR{corsConfig: createSimpleCors("*")},
			cors2:    nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			cors1:    &CorsIR{corsConfig: createSimpleCors("https://example.com")},
			cors2:    &CorsIR{corsConfig: createSimpleCors("https://example.com")},
			expected: true,
		},
		{
			name:     "different origins are not equal",
			cors1:    &CorsIR{corsConfig: createSimpleCors("https://example.com")},
			cors2:    &CorsIR{corsConfig: createSimpleCors("https://other.com")},
			expected: false,
		},
		{
			name:     "different credentials settings are not equal",
			cors1:    &CorsIR{corsConfig: createCorsWith(true)},
			cors2:    &CorsIR{corsConfig: createCorsWith(false)},
			expected: false,
		},
		{
			name:     "same credentials settings are equal",
			cors1:    &CorsIR{corsConfig: createCorsWith(true)},
			cors2:    &CorsIR{corsConfig: createCorsWith(true)},
			expected: true,
		},
		{
			name:     "nil cors config fields are equal",
			cors1:    &CorsIR{corsConfig: nil},
			cors2:    &CorsIR{corsConfig: nil},
			expected: true,
		},
		{
			name:     "nil vs non-nil cors config fields are not equal",
			cors1:    &CorsIR{corsConfig: nil},
			cors2:    &CorsIR{corsConfig: createSimpleCors("*")},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cors1.Equals(tt.cors2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.cors2.Equals(tt.cors1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		cors := &CorsIR{corsConfig: createSimpleCors("https://test.com")}
		assert.True(t, cors.Equals(cors), "cors should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameCors := func() *CorsIR {
			return &CorsIR{corsConfig: createCorsWith(false)}
		}

		a := createSameCors()
		b := createSameCors()
		c := createSameCors()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestCorsIRValidate(t *testing.T) {
	// Note: CORS validation is very permissive at the protobuf level.
	// Most validation constraints are applied at the Kubernetes API level.
	tests := []struct {
		name string
		cors *CorsIR
	}{
		{
			name: "nil cors is valid",
			cors: nil,
		},
		{
			name: "cors with nil config is valid",
			cors: &CorsIR{
				corsConfig: nil,
			},
		},
		{
			name: "valid cors config passes validation",
			cors: &CorsIR{
				corsConfig: &corsv3.CorsPolicy{
					AllowCredentials: wrapperspb.Bool(true),
				},
			},
		},
		{
			name: "empty cors config is valid",
			cors: &CorsIR{
				corsConfig: &corsv3.CorsPolicy{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cors.Validate()
			assert.NoError(t, err, "CORS validation should never fail at protobuf level")
		})
	}
}

func TestCorsIRMergeInto(t *testing.T) {
	createCorsConfig := func(allowCredentials bool) *CorsIR {
		return &CorsIR{
			corsConfig: &corsv3.CorsPolicy{
				AllowCredentials: wrapperspb.Bool(allowCredentials),
			},
		}
	}

	tests := []struct {
		name              string
		p1Cors            *CorsIR
		p2Cors            *CorsIR
		strategy          policy.MergeOptions
		expectOriginSet   bool
		expectCredentials *bool
	}{
		{
			name:              "shallow merge with nil p1",
			p1Cors:            nil,
			p2Cors:            createCorsConfig(true),
			strategy:          policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet:   true,
			expectCredentials: ptr.To(true),
		},
		{
			name:              "deep merge preserves p1 when set",
			p1Cors:            createCorsConfig(false),
			p2Cors:            createCorsConfig(true),
			strategy:          policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet:   false,
			expectCredentials: ptr.To(false),
		},
		{
			name:              "overridable shallow merge overwrites p1",
			p1Cors:            createCorsConfig(false),
			p2Cors:            createCorsConfig(true),
			strategy:          policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet:   true,
			expectCredentials: ptr.To(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{cors: tt.p1Cors}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{cors: tt.p2Cors}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			if tt.p2Cors != nil {
				tt.p2Cors.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectCredentials != nil {
				assert.NotNil(t, p1.spec.cors, "p1.cors should not be nil")
				assert.NotNil(t, p1.spec.cors.corsConfig, "p1.cors.corsConfig should not be nil")
				if p1.spec.cors.corsConfig.AllowCredentials != nil {
					assert.Equal(t, *tt.expectCredentials, p1.spec.cors.corsConfig.AllowCredentials.Value)
				}
			}

			if tt.expectOriginSet {
				origins := mergeOrigins.Get("cors")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("cors")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}
