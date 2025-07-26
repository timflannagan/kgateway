package trafficpolicy

import (
	"context"
	"testing"

	envoyroutev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestAutoHostRewriteIREquals(t *testing.T) {
	tests := []struct {
		name     string
		rewrite1 *AutoHostRewriteIR
		rewrite2 *AutoHostRewriteIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			rewrite1: nil,
			rewrite2: nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			rewrite1: nil,
			rewrite2: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			rewrite1: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			rewrite2: nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			rewrite1: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			rewrite2: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			expected: true,
		},
		{
			name:     "same false values are equal",
			rewrite1: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)},
			rewrite2: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)},
			expected: true,
		},
		{
			name:     "different values are not equal",
			rewrite1: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			rewrite2: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)},
			expected: false,
		},
		{
			name:     "nil proto fields are equal",
			rewrite1: &AutoHostRewriteIR{autoHostRewrite: nil},
			rewrite2: &AutoHostRewriteIR{autoHostRewrite: nil},
			expected: true,
		},
		{
			name:     "nil vs non-nil proto fields are not equal",
			rewrite1: &AutoHostRewriteIR{autoHostRewrite: nil},
			rewrite2: &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rewrite1.Equals(tt.rewrite2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.rewrite2.Equals(tt.rewrite1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		rewrite := &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)}
		assert.True(t, rewrite.Equals(rewrite), "rewrite should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		a := &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)}
		b := &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)}
		c := &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)}

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestAutoHostRewriteIRValidate(t *testing.T) {
	// Note: AutoHostRewrite validation is minimal as it only contains a boolean wrapper.
	tests := []struct {
		name    string
		rewrite *AutoHostRewriteIR
	}{
		{
			name:    "nil rewrite is valid",
			rewrite: nil,
		},
		{
			name: "rewrite with nil field is valid",
			rewrite: &AutoHostRewriteIR{
				autoHostRewrite: nil,
			},
		},
		{
			name: "valid true rewrite passes validation",
			rewrite: &AutoHostRewriteIR{
				autoHostRewrite: wrapperspb.Bool(true),
			},
		},
		{
			name: "valid false rewrite passes validation",
			rewrite: &AutoHostRewriteIR{
				autoHostRewrite: wrapperspb.Bool(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rewrite.Validate()
			assert.NoError(t, err, "AutoHostRewrite validation should never fail at protobuf level")
		})
	}
}

func TestAutoHostRewriteIRMergeInto(t *testing.T) {
	tests := []struct {
		name            string
		p1Rewrite       *AutoHostRewriteIR
		p2Rewrite       *AutoHostRewriteIR
		strategy        policy.MergeOptions
		expectedP1      *AutoHostRewriteIR
		expectOriginSet bool
	}{
		{
			name:            "shallow merge with nil p1",
			p1Rewrite:       nil,
			p2Rewrite:       &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectedP1:      &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			expectOriginSet: true,
		},
		{
			name:            "deep merge preserves p1 when set",
			p1Rewrite:       &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)},
			p2Rewrite:       &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectedP1:      &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)},
			expectOriginSet: false,
		},
		{
			name:            "overridable shallow merge overwrites p1",
			p1Rewrite:       &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(false)},
			p2Rewrite:       &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			strategy:        policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectedP1:      &AutoHostRewriteIR{autoHostRewrite: wrapperspb.Bool(true)},
			expectOriginSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{autoHostRewrite: tt.p1Rewrite}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{autoHostRewrite: tt.p2Rewrite}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			if tt.p2Rewrite != nil {
				tt.p2Rewrite.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			if tt.expectedP1 == nil {
				assert.Nil(t, p1.spec.autoHostRewrite)
			} else {
				assert.NotNil(t, p1.spec.autoHostRewrite)
				if tt.expectedP1.autoHostRewrite != nil && p1.spec.autoHostRewrite.autoHostRewrite != nil {
					assert.Equal(t, tt.expectedP1.autoHostRewrite.Value, p1.spec.autoHostRewrite.autoHostRewrite.Value)
				}
			}

			if tt.expectOriginSet {
				origins := mergeOrigins.Get("autoHostRewrite")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("autoHostRewrite")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}

func TestApplyForRoute_SetsRouteActionFlag(t *testing.T) {
	ctx := context.Background()
	plugin := &trafficPolicyPluginGwPass{}

	t.Run("autoHostRewrite true - RouteAction flag set", func(t *testing.T) {
		policy := &TrafficPolicy{
			spec: trafficPolicySpecIr{
				autoHostRewrite: &AutoHostRewriteIR{
					autoHostRewrite: wrapperspb.Bool(true),
				},
			},
		}

		pCtx := &ir.RouteContext{Policy: policy}
		out := &envoyroutev3.Route{
			Action: &envoyroutev3.Route_Route{
				Route: &envoyroutev3.RouteAction{},
			},
		}

		require.NoError(t, plugin.ApplyForRoute(ctx, pCtx, out))

		ra := out.GetRoute()
		require.NotNil(t, ra)
		assert.NotNil(t, ra.GetAutoHostRewrite())
		assert.True(t, ra.GetAutoHostRewrite().GetValue())
	})

	t.Run("autoHostRewrite nil - RouteAction untouched", func(t *testing.T) {
		policy := &TrafficPolicy{
			spec: trafficPolicySpecIr{autoHostRewrite: nil},
		}
		pCtx := &ir.RouteContext{Policy: policy}
		out := &envoyroutev3.Route{
			Action: &envoyroutev3.Route_Route{Route: &envoyroutev3.RouteAction{}},
		}

		require.NoError(t, plugin.ApplyForRoute(ctx, pCtx, out))

		ra := out.GetRoute()
		require.NotNil(t, ra)
		assert.Nil(t, ra.HostRewriteSpecifier) // nothing written
	})
}
