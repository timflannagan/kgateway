package trafficpolicy

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	envoycorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyroutev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_ext_authz_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/plugins"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestExtAuthIREquals(t *testing.T) {
	// Helper to create simple extauth configurations for testing
	createSimpleExtAuth := func(disabled bool) *envoy_ext_authz_v3.ExtAuthzPerRoute {
		return &envoy_ext_authz_v3.ExtAuthzPerRoute{
			Override: &envoy_ext_authz_v3.ExtAuthzPerRoute_Disabled{
				Disabled: disabled,
			},
		}
	}

	createProvider := func(name string) *TrafficPolicyGatewayExtensionIR {
		return &TrafficPolicyGatewayExtensionIR{
			Name: name,
			ExtAuth: &envoy_ext_authz_v3.ExtAuthz{
				Services: &envoy_ext_authz_v3.ExtAuthz_GrpcService{
					GrpcService: &envoycorev3.GrpcService{
						TargetSpecifier: &envoycorev3.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoycorev3.GrpcService_EnvoyGrpc{
								ClusterName: name,
							},
						},
					},
				},
			},
		}
	}

	createEnablement := func(disableAll bool) *v1alpha1.ExtAuthEnabled {
		if disableAll {
			return ptr.To(v1alpha1.ExtAuthDisableAll)
		}
		return nil // No specific enablement setting
	}

	// Shared enablement for testing pointer equality
	sharedEnablementTrue := createEnablement(true)

	tests := []struct {
		name     string
		extauth1 *ExtAuthIR
		extauth2 *ExtAuthIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			extauth1: nil,
			extauth2: nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			extauth1: nil,
			extauth2: &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(false)},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			extauth1: &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(false)},
			extauth2: nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			extauth1: &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(false)},
			extauth2: &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(false)},
			expected: true,
		},
		{
			name:     "different disabled settings are not equal",
			extauth1: &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(true)},
			extauth2: &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(false)},
			expected: false,
		},
		{
			name:     "different providers are not equal",
			extauth1: &ExtAuthIR{provider: createProvider("service1")},
			extauth2: &ExtAuthIR{provider: createProvider("service2")},
			expected: false,
		},
		{
			name:     "same providers are equal",
			extauth1: &ExtAuthIR{provider: createProvider("service1")},
			extauth2: &ExtAuthIR{provider: createProvider("service1")},
			expected: true,
		},
		{
			name:     "different enablement settings are not equal",
			extauth1: &ExtAuthIR{enablement: createEnablement(true)},
			extauth2: &ExtAuthIR{enablement: createEnablement(false)},
			expected: false,
		},
		{
			name:     "same enablement settings are equal",
			extauth1: &ExtAuthIR{enablement: sharedEnablementTrue},
			extauth2: &ExtAuthIR{enablement: sharedEnablementTrue},
			expected: true,
		},
		{
			name:     "nil extauth fields are equal",
			extauth1: &ExtAuthIR{extauthPerRoute: nil},
			extauth2: &ExtAuthIR{extauthPerRoute: nil},
			expected: true,
		},
		{
			name:     "nil vs non-nil extauth fields are not equal",
			extauth1: &ExtAuthIR{extauthPerRoute: nil},
			extauth2: &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(false)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.extauth1.Equals(tt.extauth2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.extauth2.Equals(tt.extauth1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		extauth := &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(false)}
		assert.True(t, extauth.Equals(extauth), "extauth should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameExtAuth := func() *ExtAuthIR {
			return &ExtAuthIR{extauthPerRoute: createSimpleExtAuth(true)}
		}

		a := createSameExtAuth()
		b := createSameExtAuth()
		c := createSameExtAuth()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestExtAuthIRValidate(t *testing.T) {
	// Note: ExtAuth validation is permissive at the protobuf level.
	// Provider validation is tested separately.
	tests := []struct {
		name    string
		extauth *ExtAuthIR
	}{
		{
			name:    "nil extauth is valid",
			extauth: nil,
		},
		{
			name: "extauth with nil fields is valid",
			extauth: &ExtAuthIR{
				extauthPerRoute: nil,
				provider:        nil,
			},
		},
		{
			name: "valid extauth config passes validation",
			extauth: &ExtAuthIR{
				extauthPerRoute: &envoy_ext_authz_v3.ExtAuthzPerRoute{
					Override: &envoy_ext_authz_v3.ExtAuthzPerRoute_Disabled{
						Disabled: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.extauth.Validate()
			assert.NoError(t, err, "ExtAuth validation should never fail at protobuf level")
		})
	}
}

func TestExtAuthIRMergeInto(t *testing.T) {
	// Helper to create extauth configurations
	createExtAuthConfig := func(disabled bool) *ExtAuthIR {
		return &ExtAuthIR{
			extauthPerRoute: &envoy_ext_authz_v3.ExtAuthzPerRoute{
				Override: &envoy_ext_authz_v3.ExtAuthzPerRoute_Disabled{
					Disabled: disabled,
				},
			},
		}
	}

	tests := []struct {
		name            string
		p1ExtAuth       *ExtAuthIR
		p2ExtAuth       *ExtAuthIR
		strategy        policy.MergeOptions
		expectOriginSet bool
		expectDisabled  *bool
	}{
		{
			name:            "shallow merge with nil p1",
			p1ExtAuth:       nil,
			p2ExtAuth:       createExtAuthConfig(false),
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet: true,
			expectDisabled:  ptr.To(false),
		},
		{
			name:            "deep merge preserves p1 when set",
			p1ExtAuth:       createExtAuthConfig(true),
			p2ExtAuth:       createExtAuthConfig(false),
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet: false,        // No merge should happen
			expectDisabled:  ptr.To(true), // Should keep p1
		},
		{
			name:            "overridable shallow merge overwrites p1",
			p1ExtAuth:       createExtAuthConfig(true),
			p2ExtAuth:       createExtAuthConfig(false),
			strategy:        policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet: true,
			expectDisabled:  ptr.To(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test policies
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{extAuth: tt.p1ExtAuth}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{extAuth: tt.p2ExtAuth}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			// Perform merge
			if tt.p2ExtAuth != nil {
				tt.p2ExtAuth.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectDisabled != nil {
				assert.NotNil(t, p1.spec.extAuth, "p1.extAuth should not be nil")
				assert.NotNil(t, p1.spec.extAuth.extauthPerRoute, "p1.extAuth.extauthPerRoute should not be nil")

				// Verify disabled setting
				disabled := p1.spec.extAuth.extauthPerRoute.GetDisabled()
				assert.Equal(t, *tt.expectDisabled, disabled)
			}

			// Verify origin tracking
			if tt.expectOriginSet {
				origins := mergeOrigins.Get("extAuth")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("extAuth")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}

func TestExtAuthForSpec(t *testing.T) {
	t.Run("configures request body settings", func(t *testing.T) {
		truthy := true
		// Setup
		spec := &v1alpha1.TrafficPolicy{Spec: v1alpha1.TrafficPolicySpec{
			ExtAuth: &v1alpha1.ExtAuthPolicy{
				ExtensionRef: &corev1.LocalObjectReference{
					Name: "test-extension",
				},
				WithRequestBody: &v1alpha1.BufferSettings{
					MaxRequestBytes:     1024,
					AllowPartialMessage: &truthy,
					PackAsBytes:         &truthy,
				},
			},
		}}

		// Execute
		extauthPerRoute := translatePerFilterConfig(spec.Spec.ExtAuth)

		// Verify
		require.NotNil(t, extauthPerRoute)
		require.NotNil(t, extauthPerRoute.GetCheckSettings().WithRequestBody)
		assert.Equal(t, uint32(1024), extauthPerRoute.GetCheckSettings().WithRequestBody.MaxRequestBytes)
		assert.True(t, extauthPerRoute.GetCheckSettings().WithRequestBody.AllowPartialMessage)
		assert.True(t, extauthPerRoute.GetCheckSettings().WithRequestBody.PackAsBytes)
	})
}

func TestApplyForRoute(t *testing.T) {
	t.Run("applies ext auth configuration to route", func(t *testing.T) {
		// Setup
		plugin := &trafficPolicyPluginGwPass{}
		ctx := context.Background()
		policy := &TrafficPolicy{
			spec: trafficPolicySpecIr{
				extAuth: &ExtAuthIR{
					provider: &TrafficPolicyGatewayExtensionIR{
						Name:    "test-extension",
						ExtType: v1alpha1.GatewayExtensionTypeExtAuth,
						ExtAuth: &envoy_ext_authz_v3.ExtAuthz{
							FailureModeAllow: true,
						},
					},
				},
			},
		}
		pCtx := &ir.RouteContext{
			Policy: policy,
		}
		outputRoute := &envoyroutev3.Route{}

		// Execute
		err := plugin.ApplyForRoute(ctx, pCtx, outputRoute)

		// Verify
		require.NoError(t, err)
		require.NotNil(t, pCtx.TypedFilterConfig)
		extAuthConfig, ok := pCtx.TypedFilterConfig[extAuthFilterName("test-extension")]
		assert.True(t, ok)
		assert.NotNil(t, extAuthConfig)
	})

	t.Run("handles nil ext auth configuration", func(t *testing.T) {
		// Setup
		plugin := &trafficPolicyPluginGwPass{}
		ctx := context.Background()
		policy := &TrafficPolicy{
			spec: trafficPolicySpecIr{
				extAuth: nil,
			},
		}
		pCtx := &ir.RouteContext{
			Policy: policy,
		}
		outputRoute := &envoyroutev3.Route{}

		// Execute
		err := plugin.ApplyForRoute(ctx, pCtx, outputRoute)

		// Verify
		require.NoError(t, err)
		assert.Nil(t, pCtx.TypedFilterConfig)
	})
}

func TestHttpFilters(t *testing.T) {
	t.Run("adds ext auth filter to filter chain", func(t *testing.T) {
		// Setup
		plugin := &trafficPolicyPluginGwPass{
			extAuthPerProvider: ProviderNeededMap{
				Providers: map[string]map[string]*TrafficPolicyGatewayExtensionIR{
					"test-filter-chain": {
						"test-extension": {
							Name:    "test-extension",
							ExtType: v1alpha1.GatewayExtensionTypeExtAuth,
							ExtAuth: &envoy_ext_authz_v3.ExtAuthz{
								FailureModeAllow: true,
							},
						},
					},
				},
			},
		}
		ctx := context.Background()
		fcc := ir.FilterChainCommon{
			FilterChainName: "test-filter-chain",
		}

		// Execute
		filters, err := plugin.HttpFilters(ctx, fcc)

		// Verify
		require.NoError(t, err)
		require.NotNil(t, filters)
		assert.Equal(t, 2, len(filters)) // extauth and metadata filter
		assert.Equal(t, plugins.DuringStage(plugins.AuthZStage), filters[1].Stage)
	})
}

func TestExtAuthPolicyPlugin(t *testing.T) {
	t.Run("applies ext auth configuration to route", func(t *testing.T) {
		// Setup
		plugin := &trafficPolicyPluginGwPass{}
		ctx := context.Background()
		policy := &TrafficPolicy{
			spec: trafficPolicySpecIr{
				extAuth: &ExtAuthIR{
					provider: &TrafficPolicyGatewayExtensionIR{
						Name:    "test-auth-extension",
						ExtType: v1alpha1.GatewayExtensionTypeExtAuth,
						ExtAuth: &envoy_ext_authz_v3.ExtAuthz{
							FailureModeAllow: true,
							WithRequestBody: &envoy_ext_authz_v3.BufferSettings{
								MaxRequestBytes: 1024,
							},
						},
					},
				},
			},
		}
		pCtx := &ir.RouteContext{
			Policy: policy,
		}
		outputRoute := &envoyroutev3.Route{}

		// Execute
		err := plugin.ApplyForRoute(ctx, pCtx, outputRoute)

		// Verify
		require.NoError(t, err)
		require.NotNil(t, pCtx.TypedFilterConfig)
		extAuthConfig, ok := pCtx.TypedFilterConfig[extAuthFilterName("test-auth-extension")]
		assert.True(t, ok)
		assert.NotNil(t, extAuthConfig)
		assert.Empty(t, pCtx.TypedFilterConfig[extAuthGlobalDisableFilterName])
	})

	t.Run("handles disabled ext auth configuration", func(t *testing.T) {
		// Setup
		plugin := &trafficPolicyPluginGwPass{}
		ctx := context.Background()
		policy := &TrafficPolicy{
			spec: trafficPolicySpecIr{
				extAuth: &ExtAuthIR{
					enablement: ptr.To(v1alpha1.ExtAuthDisableAll),
				},
			},
		}
		pCtx := &ir.RouteContext{
			Policy: policy,
		}
		outputRoute := &envoyroutev3.Route{}

		// Execute
		err := plugin.ApplyForRoute(ctx, pCtx, outputRoute)

		// Verify
		require.NoError(t, err)
		// assert.NotNil(t, )
		assert.NotNil(t, pCtx.TypedFilterConfig, pCtx)
		assert.NotEmpty(t, pCtx.TypedFilterConfig[extAuthGlobalDisableFilterName])
	})
}
