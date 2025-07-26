package trafficpolicy

import (
	"testing"
	"time"

	localratelimitv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/local_ratelimit/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/utils/ptr"

	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestLocalRateLimitIREquals(t *testing.T) {
	// Helper to create simple rate limit configurations for testing
	createSimpleRateLimit := func(tokensPerSecond uint32) *localratelimitv3.LocalRateLimit {
		return &localratelimitv3.LocalRateLimit{
			TokenBucket: &typev3.TokenBucket{
				MaxTokens:     tokensPerSecond * 10,
				TokensPerFill: wrapperspb.UInt32(tokensPerSecond),
				FillInterval:  durationpb.New(time.Second),
			},
		}
	}

	createRateLimitWithPrefix := func(prefix string) *localratelimitv3.LocalRateLimit {
		return &localratelimitv3.LocalRateLimit{
			StatPrefix: prefix,
		}
	}

	tests := []struct {
		name       string
		rateLimit1 *LocalRateLimitIR
		rateLimit2 *LocalRateLimitIR
		expected   bool
	}{
		{
			name:       "both nil are equal",
			rateLimit1: nil,
			rateLimit2: nil,
			expected:   true,
		},
		{
			name:       "nil vs non-nil are not equal",
			rateLimit1: nil,
			rateLimit2: &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(100)},
			expected:   false,
		},
		{
			name:       "non-nil vs nil are not equal",
			rateLimit1: &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(100)},
			rateLimit2: nil,
			expected:   false,
		},
		{
			name:       "same instance is equal",
			rateLimit1: &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(100)},
			rateLimit2: &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(100)},
			expected:   true,
		},
		{
			name:       "different token rates are not equal",
			rateLimit1: &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(100)},
			rateLimit2: &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(200)},
			expected:   false,
		},
		{
			name:       "different stat prefixes are not equal",
			rateLimit1: &LocalRateLimitIR{localRateLimit: createRateLimitWithPrefix("prefix1")},
			rateLimit2: &LocalRateLimitIR{localRateLimit: createRateLimitWithPrefix("prefix2")},
			expected:   false,
		},
		{
			name:       "same stat prefixes are equal",
			rateLimit1: &LocalRateLimitIR{localRateLimit: createRateLimitWithPrefix("prefix1")},
			rateLimit2: &LocalRateLimitIR{localRateLimit: createRateLimitWithPrefix("prefix1")},
			expected:   true,
		},
		{
			name:       "nil rate limit config fields are equal",
			rateLimit1: &LocalRateLimitIR{localRateLimit: nil},
			rateLimit2: &LocalRateLimitIR{localRateLimit: nil},
			expected:   true,
		},
		{
			name:       "nil vs non-nil rate limit config fields are not equal",
			rateLimit1: &LocalRateLimitIR{localRateLimit: nil},
			rateLimit2: &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(100)},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rateLimit1.Equals(tt.rateLimit2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.rateLimit2.Equals(tt.rateLimit1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		rateLimit := &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(50)}
		assert.True(t, rateLimit.Equals(rateLimit), "rateLimit should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameRateLimit := func() *LocalRateLimitIR {
			return &LocalRateLimitIR{localRateLimit: createSimpleRateLimit(75)}
		}

		a := createSameRateLimit()
		b := createSameRateLimit()
		c := createSameRateLimit()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestLocalRateLimitIRValidate(t *testing.T) {
	tests := []struct {
		name      string
		rateLimit *LocalRateLimitIR
		wantErr   bool
	}{
		{
			name:      "nil rate limit is valid",
			rateLimit: nil,
			wantErr:   false,
		},
		{
			name:      "rate limit with nil config is valid",
			rateLimit: &LocalRateLimitIR{localRateLimit: nil},
			wantErr:   false,
		},
		{
			name: "valid rate limit config passes validation",
			rateLimit: &LocalRateLimitIR{
				localRateLimit: &localratelimitv3.LocalRateLimit{
					StatPrefix: "test_prefix",
					TokenBucket: &typev3.TokenBucket{
						MaxTokens:     1000,
						TokensPerFill: wrapperspb.UInt32(100),
						FillInterval:  durationpb.New(time.Second),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty rate limit config fails validation",
			rateLimit: &LocalRateLimitIR{
				localRateLimit: &localratelimitv3.LocalRateLimit{},
			},
			wantErr: true,
		},
		{
			name: "rate limit config with invalid fill interval fails validation",
			rateLimit: &LocalRateLimitIR{
				localRateLimit: &localratelimitv3.LocalRateLimit{
					StatPrefix: "test_prefix",
					TokenBucket: &typev3.TokenBucket{
						MaxTokens:     100,
						TokensPerFill: wrapperspb.UInt32(10),
						FillInterval:  &durationpb.Duration{Seconds: -1},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rateLimit.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLocalRateLimitIRMergeInto(t *testing.T) {
	// Helper to create rate limit configurations
	createRateLimitConfig := func(tokensPerSecond uint32) *LocalRateLimitIR {
		return &LocalRateLimitIR{
			localRateLimit: &localratelimitv3.LocalRateLimit{
				TokenBucket: &typev3.TokenBucket{
					MaxTokens:     tokensPerSecond * 10,
					TokensPerFill: wrapperspb.UInt32(tokensPerSecond),
					FillInterval:  durationpb.New(time.Second),
				},
			},
		}
	}

	tests := []struct {
		name            string
		p1RateLimit     *LocalRateLimitIR
		p2RateLimit     *LocalRateLimitIR
		strategy        policy.MergeOptions
		expectOriginSet bool
		expectTokens    *uint32
	}{
		{
			name:            "shallow merge with nil p1",
			p1RateLimit:     nil,
			p2RateLimit:     createRateLimitConfig(100),
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet: true,
			expectTokens:    ptr.To(uint32(100)),
		},
		{
			name:            "deep merge preserves p1 when set",
			p1RateLimit:     createRateLimitConfig(50),
			p2RateLimit:     createRateLimitConfig(100),
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet: false,
			expectTokens:    ptr.To(uint32(50)),
		},
		{
			name:            "overridable shallow merge overwrites p1",
			p1RateLimit:     createRateLimitConfig(50),
			p2RateLimit:     createRateLimitConfig(100),
			strategy:        policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet: true,
			expectTokens:    ptr.To(uint32(100)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{localRateLimit: tt.p1RateLimit}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{localRateLimit: tt.p2RateLimit}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			if tt.p2RateLimit != nil {
				tt.p2RateLimit.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectTokens != nil {
				assert.NotNil(t, p1.spec.localRateLimit, "p1.localRateLimit should not be nil")
				assert.NotNil(t, p1.spec.localRateLimit.localRateLimit, "p1.localRateLimit.localRateLimit should not be nil")
				if p1.spec.localRateLimit.localRateLimit.TokenBucket != nil && p1.spec.localRateLimit.localRateLimit.TokenBucket.TokensPerFill != nil {
					assert.Equal(t, *tt.expectTokens, p1.spec.localRateLimit.localRateLimit.TokenBucket.TokensPerFill.Value)
				}
			}

			if tt.expectOriginSet {
				origins := mergeOrigins.Get("rateLimit.local")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("rateLimit.local")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}
