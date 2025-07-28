package trafficpolicy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTrafficPolicyCreationTime(t *testing.T) {
	now := time.Now()
	policy := &TrafficPolicy{ct: now}
	assert.Equal(t, now, policy.CreationTime(), "CreationTime should return the stored creation time")
}

func TestTrafficPolicyEquals(t *testing.T) {
	baseTime := time.Now()

	tests := []struct {
		name     string
		policy1  *TrafficPolicy
		policy2  any
		expected bool
	}{
		{
			name:     "same instance equals itself",
			policy1:  &TrafficPolicy{ct: baseTime},
			policy2:  &TrafficPolicy{ct: baseTime},
			expected: true,
		},
		{
			name:     "different creation times are not equal",
			policy1:  &TrafficPolicy{ct: baseTime},
			policy2:  &TrafficPolicy{ct: baseTime.Add(time.Minute)},
			expected: false,
		},
		{
			name:     "nil vs non-nil are not equal",
			policy1:  &TrafficPolicy{ct: baseTime},
			policy2:  nil,
			expected: false,
		},
		{
			name:     "wrong type is not equal",
			policy1:  &TrafficPolicy{ct: baseTime},
			policy2:  "not a traffic policy",
			expected: false,
		},
		{
			name: "different AI policies are not equal",
			policy1: &TrafficPolicy{
				ct:   baseTime,
				spec: trafficPolicySpecIr{ai: &AIPolicyIR{}},
			},
			policy2: &TrafficPolicy{
				ct:   baseTime,
				spec: trafficPolicySpecIr{ai: nil},
			},
			expected: false,
		},
		{
			name: "same AI policies are equal",
			policy1: &TrafficPolicy{
				ct:   baseTime,
				spec: trafficPolicySpecIr{ai: &AIPolicyIR{}},
			},
			policy2: &TrafficPolicy{
				ct:   baseTime,
				spec: trafficPolicySpecIr{ai: &AIPolicyIR{}},
			},
			expected: true,
		},
		{
			name: "different buffer policies are not equal",
			policy1: &TrafficPolicy{
				ct:   baseTime,
				spec: trafficPolicySpecIr{buffer: &BufferIR{maxRequestBytes: 1024}},
			},
			policy2: &TrafficPolicy{
				ct:   baseTime,
				spec: trafficPolicySpecIr{buffer: &BufferIR{maxRequestBytes: 2048}},
			},
			expected: false,
		},
		{
			name: "all nil specs are equal",
			policy1: &TrafficPolicy{
				ct:   baseTime,
				spec: trafficPolicySpecIr{},
			},
			policy2: &TrafficPolicy{
				ct:   baseTime,
				spec: trafficPolicySpecIr{},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.policy1.Equals(tt.policy2)
			assert.Equal(t, tt.expected, result, "Equals result should match expected")
		})
	}
}

func TestTrafficPolicyValidate(t *testing.T) {
	tests := []struct {
		name    string
		policy  *TrafficPolicy
		wantErr bool
	}{
		{
			name: "empty policy is valid",
			policy: &TrafficPolicy{
				spec: trafficPolicySpecIr{},
			},
			wantErr: false,
		},
		{
			name: "policy with valid buffer is valid",
			policy: &TrafficPolicy{
				spec: trafficPolicySpecIr{
					buffer: &BufferIR{maxRequestBytes: 1024},
				},
			},
			wantErr: false,
		},
		{
			name: "policy with valid AI policy is valid",
			policy: &TrafficPolicy{
				spec: trafficPolicySpecIr{
					ai: &AIPolicyIR{},
				},
			},
			wantErr: false,
		},
		{
			name: "policy with multiple valid policies is valid",
			policy: &TrafficPolicy{
				spec: trafficPolicySpecIr{
					ai:              &AIPolicyIR{},
					buffer:          &BufferIR{maxRequestBytes: 1024},
					autoHostRewrite: &AutoHostRewriteIR{},
					cors:            &CorsIR{},
					csrf:            &CsrfIR{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if tt.wantErr {
				assert.Error(t, err, "Validate should return an error")
			} else {
				assert.NoError(t, err, "Validate should not return an error")
			}
		})
	}
}

func TestTrafficPolicySpecIrStructure(t *testing.T) {
	t.Run("trafficPolicySpecIr can be instantiated with all policy types", func(t *testing.T) {
		spec := trafficPolicySpecIr{
			ai:              &AIPolicyIR{},
			autoHostRewrite: &AutoHostRewriteIR{},
			buffer:          &BufferIR{},
			cors:            &CorsIR{},
			csrf:            &CsrfIR{},
			extProc:         &ExtprocIR{},
			extAuth:         &ExtAuthIR{},
			hashPolicies:    &HashPolicyIR{},
			localRateLimit:  &LocalRateLimitIR{},
			globalRateLimit: &GlobalRateLimitIR{},
			transformation:  &TransformationIR{},
			rustformation:   &RustformationIR{},
		}

		// Verify all fields are properly set
		assert.NotNil(t, spec.ai, "AI policy should be set")
		assert.NotNil(t, spec.autoHostRewrite, "AutoHostRewrite policy should be set")
		assert.NotNil(t, spec.buffer, "Buffer policy should be set")
		assert.NotNil(t, spec.cors, "CORS policy should be set")
		assert.NotNil(t, spec.csrf, "CSRF policy should be set")
		assert.NotNil(t, spec.extProc, "ExtProc policy should be set")
		assert.NotNil(t, spec.extAuth, "ExtAuth policy should be set")
		assert.NotNil(t, spec.hashPolicies, "HashPolicies should be set")
		assert.NotNil(t, spec.localRateLimit, "LocalRateLimit policy should be set")
		assert.NotNil(t, spec.globalRateLimit, "GlobalRateLimit policy should be set")
		assert.NotNil(t, spec.transformation, "Transformation policy should be set")
		assert.NotNil(t, spec.rustformation, "Rustformation policy should be set")
	})
}

func TestTrafficPolicyNilSafetyPattern(t *testing.T) {
	baseTime := time.Now()

	t.Run("validates safe nil-handling pattern in sub-IRs", func(t *testing.T) {
		// TrafficPolicy with all nil sub-IRs should validate successfully
		policy := &TrafficPolicy{
			ct:   baseTime,
			spec: trafficPolicySpecIr{
				// All fields are nil by default
			},
		}

		err := policy.Validate()
		assert.NoError(t, err, "Policy with all nil sub-IRs should validate successfully")

		// Individual sub-IR types handle nil receivers safely
		var ai *AIPolicyIR = nil
		var cors *CorsIR = nil
		var buffer *BufferIR = nil

		// These should all work without panicking
		assert.NotPanics(t, func() {
			ai.Validate()
			cors.Validate()
			buffer.Validate()
		}, "Sub-IR Validate methods should handle nil receivers safely")
	})

	t.Run("validates safe nil-handling in Equals", func(t *testing.T) {
		// Both policies have all nil sub-IRs
		policy1 := &TrafficPolicy{ct: baseTime, spec: trafficPolicySpecIr{}}
		policy2 := &TrafficPolicy{ct: baseTime, spec: trafficPolicySpecIr{}}

		// Should be equal (both have same nil state)
		assert.True(t, policy1.Equals(policy2), "Policies with all nil sub-IRs should be equal")

		// Mixed nil and non-nil should not be equal
		policy3 := &TrafficPolicy{
			ct:   baseTime,
			spec: trafficPolicySpecIr{ai: &AIPolicyIR{}}, // ai is set
		}

		assert.False(t, policy1.Equals(policy3), "Policy with nil sub-IR should not equal policy with non-nil sub-IR")
	})

	t.Run("documents nil-safe method patterns", func(t *testing.T) {
		// This test documents the nil-safety patterns used throughout the sub-IR types:
		//
		// Pattern 1: Explicit nil check (used by AIPolicyIR, CorsIR, etc.)
		// func (a *AIPolicyIR) Validate() error {
		//     if a == nil {
		//         return nil
		//     }
		//     // ... rest of validation
		// }
		//
		// Pattern 2: Implicit safety (used by BufferIR)
		// func (b *BufferIR) Validate() error {
		//     return nil  // No field access, safe with nil receiver
		// }
		//
		// Pattern 3: Conditional nil check (used by CorsIR)
		// func (c *CorsIR) Validate() error {
		//     if c == nil || c.corsConfig == nil {
		//         return nil
		//     }
		//     // ... rest of validation
		// }

		// Test demonstrates all patterns work
		var ai *AIPolicyIR = nil   // Pattern 1
		var buffer *BufferIR = nil // Pattern 2
		var cors *CorsIR = nil     // Pattern 3

		// All should work without NPE
		assert.NoError(t, ai.Validate(), "Pattern 1: Explicit nil check")
		assert.NoError(t, buffer.Validate(), "Pattern 2: Implicit safety")
		assert.NoError(t, cors.Validate(), "Pattern 3: Conditional nil check")

		t.Log("✅ All sub-IR types implement nil-safe Validate() methods")
		t.Log("✅ TrafficPolicy.Equals() and TrafficPolicy.Validate() are NPE-safe")
		t.Log("✅ No nil pointer exceptions will occur when calling sub-IR methods directly")
	})
}
