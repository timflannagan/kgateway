package trafficpolicy

import (
	"testing"

	envoycorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	csrfv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/csrf/v3"
	envoy_type_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"

	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestCsrfIREquals(t *testing.T) {
	createSimpleCsrf := func(additionalOrigin string) *csrfv3.CsrfPolicy {
		return &csrfv3.CsrfPolicy{
			AdditionalOrigins: []*envoy_type_matcher_v3.StringMatcher{
				{
					MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
						Exact: additionalOrigin,
					},
				},
			},
			FilterEnabled: &envoycorev3.RuntimeFractionalPercent{
				DefaultValue: &typev3.FractionalPercent{
					Numerator:   100,
					Denominator: typev3.FractionalPercent_HUNDRED,
				},
			},
		}
	}
	createCsrfWithShadowEnabled := func(shadowEnabled bool) *csrfv3.CsrfPolicy {
		return &csrfv3.CsrfPolicy{
			ShadowEnabled: &envoycorev3.RuntimeFractionalPercent{
				DefaultValue: &typev3.FractionalPercent{
					Numerator: func() uint32 {
						if shadowEnabled {
							return 100
						}
						return 0
					}(),
					Denominator: typev3.FractionalPercent_HUNDRED,
				},
			},
		}
	}

	tests := []struct {
		name     string
		csrf1    *CsrfIR
		csrf2    *CsrfIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			csrf1:    nil,
			csrf2:    nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			csrf1:    nil,
			csrf2:    &CsrfIR{csrfPolicy: createSimpleCsrf("https://example.com")},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			csrf1:    &CsrfIR{csrfPolicy: createSimpleCsrf("https://example.com")},
			csrf2:    nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			csrf1:    &CsrfIR{csrfPolicy: createSimpleCsrf("https://example.com")},
			csrf2:    &CsrfIR{csrfPolicy: createSimpleCsrf("https://example.com")},
			expected: true,
		},
		{
			name:     "different origins are not equal",
			csrf1:    &CsrfIR{csrfPolicy: createSimpleCsrf("https://example.com")},
			csrf2:    &CsrfIR{csrfPolicy: createSimpleCsrf("https://other.com")},
			expected: false,
		},
		{
			name:     "different shadow enabled settings are not equal",
			csrf1:    &CsrfIR{csrfPolicy: createCsrfWithShadowEnabled(true)},
			csrf2:    &CsrfIR{csrfPolicy: createCsrfWithShadowEnabled(false)},
			expected: false,
		},
		{
			name:     "same shadow enabled settings are equal",
			csrf1:    &CsrfIR{csrfPolicy: createCsrfWithShadowEnabled(true)},
			csrf2:    &CsrfIR{csrfPolicy: createCsrfWithShadowEnabled(true)},
			expected: true,
		},
		{
			name:     "nil csrf config fields are equal",
			csrf1:    &CsrfIR{csrfPolicy: nil},
			csrf2:    &CsrfIR{csrfPolicy: nil},
			expected: true,
		},
		{
			name:     "nil vs non-nil csrf config fields are not equal",
			csrf1:    &CsrfIR{csrfPolicy: nil},
			csrf2:    &CsrfIR{csrfPolicy: createSimpleCsrf("https://test.com")},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.csrf1.Equals(tt.csrf2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.csrf2.Equals(tt.csrf1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		csrf := &CsrfIR{csrfPolicy: createSimpleCsrf("https://test.com")}
		assert.True(t, csrf.Equals(csrf), "csrf should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameCsrf := func() *CsrfIR {
			return &CsrfIR{csrfPolicy: createCsrfWithShadowEnabled(false)}
		}

		a := createSameCsrf()
		b := createSameCsrf()
		c := createSameCsrf()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestCsrfIRValidate(t *testing.T) {
	tests := []struct {
		name    string
		csrf    *CsrfIR
		wantErr bool
	}{
		{
			name:    "nil csrf is valid",
			csrf:    nil,
			wantErr: false,
		},
		{
			name:    "csrf with nil config is valid",
			csrf:    &CsrfIR{csrfPolicy: nil},
			wantErr: false,
		},
		{
			name: "valid csrf config passes validation",
			csrf: &CsrfIR{
				csrfPolicy: &csrfv3.CsrfPolicy{
					FilterEnabled: &envoycorev3.RuntimeFractionalPercent{
						DefaultValue: &typev3.FractionalPercent{
							Numerator:   100,
							Denominator: typev3.FractionalPercent_HUNDRED,
						},
					},
					ShadowEnabled: &envoycorev3.RuntimeFractionalPercent{
						DefaultValue: &typev3.FractionalPercent{
							Numerator:   100,
							Denominator: typev3.FractionalPercent_HUNDRED,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty csrf config fails validation",
			csrf: &CsrfIR{
				csrfPolicy: &csrfv3.CsrfPolicy{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.csrf.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCsrfIRMergeInto(t *testing.T) {
	// Helper to create CSRF configurations
	createCsrfConfig := func(shadowEnabled bool) *CsrfIR {
		return &CsrfIR{
			csrfPolicy: &csrfv3.CsrfPolicy{
				ShadowEnabled: &envoycorev3.RuntimeFractionalPercent{
					DefaultValue: &typev3.FractionalPercent{
						Numerator: func() uint32 {
							if shadowEnabled {
								return 100
							}
							return 0
						}(),
						Denominator: typev3.FractionalPercent_HUNDRED,
					},
				},
			},
		}
	}

	tests := []struct {
		name                string
		p1Csrf              *CsrfIR
		p2Csrf              *CsrfIR
		strategy            policy.MergeOptions
		expectOriginSet     bool
		expectShadowEnabled *bool
	}{
		{
			name:                "shallow merge with nil p1",
			p1Csrf:              nil,
			p2Csrf:              createCsrfConfig(true),
			strategy:            policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet:     true,
			expectShadowEnabled: ptr.To(true),
		},
		{
			name:                "deep merge preserves p1 when set",
			p1Csrf:              createCsrfConfig(false),
			p2Csrf:              createCsrfConfig(true),
			strategy:            policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet:     false,         // No merge should happen
			expectShadowEnabled: ptr.To(false), // Should keep p1
		},
		{
			name:                "overridable shallow merge overwrites p1",
			p1Csrf:              createCsrfConfig(false),
			p2Csrf:              createCsrfConfig(true),
			strategy:            policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet:     true,
			expectShadowEnabled: ptr.To(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test policies
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{csrf: tt.p1Csrf}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{csrf: tt.p2Csrf}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			// Perform merge
			if tt.p2Csrf != nil {
				tt.p2Csrf.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectShadowEnabled != nil {
				assert.NotNil(t, p1.spec.csrf, "p1.csrf should not be nil")
				assert.NotNil(t, p1.spec.csrf.csrfPolicy, "p1.csrf.csrfPolicy should not be nil")

				// Verify shadow enabled setting
				if p1.spec.csrf.csrfPolicy.ShadowEnabled != nil && p1.spec.csrf.csrfPolicy.ShadowEnabled.DefaultValue != nil {
					expectedNumerator := uint32(0)
					if *tt.expectShadowEnabled {
						expectedNumerator = 100
					}
					assert.Equal(t, expectedNumerator, p1.spec.csrf.csrfPolicy.ShadowEnabled.DefaultValue.Numerator)
				}
			}

			// Verify origin tracking
			if tt.expectOriginSet {
				origins := mergeOrigins.Get("csrf")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("csrf")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}
