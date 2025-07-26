package trafficpolicy

import (
	"testing"

	envoycorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_ext_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestExtprocIREquals(t *testing.T) {
	createSimpleExtproc := func(requestMode envoy_ext_proc_v3.ProcessingMode_HeaderSendMode) *envoy_ext_proc_v3.ExtProcPerRoute {
		return &envoy_ext_proc_v3.ExtProcPerRoute{
			Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{
				Overrides: &envoy_ext_proc_v3.ExtProcOverrides{
					ProcessingMode: &envoy_ext_proc_v3.ProcessingMode{
						RequestHeaderMode: requestMode,
					},
				},
			},
		}
	}
	createProvider := func(name string) *TrafficPolicyGatewayExtensionIR {
		return &TrafficPolicyGatewayExtensionIR{
			Name: name,
			ExtProc: &envoy_ext_proc_v3.ExternalProcessor{
				GrpcService: &envoycorev3.GrpcService{
					TargetSpecifier: &envoycorev3.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &envoycorev3.GrpcService_EnvoyGrpc{
							ClusterName: name,
						},
					},
				},
			},
		}
	}

	tests := []struct {
		name     string
		extproc1 *ExtprocIR
		extproc2 *ExtprocIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			extproc1: nil,
			extproc2: nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			extproc1: nil,
			extproc2: &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SEND)},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			extproc1: &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SEND)},
			extproc2: nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			extproc1: &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SEND)},
			extproc2: &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SEND)},
			expected: true,
		},
		{
			name:     "different processing modes are not equal",
			extproc1: &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SEND)},
			extproc2: &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SKIP)},
			expected: false,
		},
		{
			name:     "different providers are not equal",
			extproc1: &ExtprocIR{provider: createProvider("service1")},
			extproc2: &ExtprocIR{provider: createProvider("service2")},
			expected: false,
		},
		{
			name:     "same providers are equal",
			extproc1: &ExtprocIR{provider: createProvider("service1")},
			extproc2: &ExtprocIR{provider: createProvider("service1")},
			expected: true,
		},
		{
			name:     "nil perRoute fields are equal",
			extproc1: &ExtprocIR{perRoute: nil},
			extproc2: &ExtprocIR{perRoute: nil},
			expected: true,
		},
		{
			name:     "nil vs non-nil perRoute fields are not equal",
			extproc1: &ExtprocIR{perRoute: nil},
			extproc2: &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SEND)},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.extproc1.Equals(tt.extproc2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.extproc2.Equals(tt.extproc1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		extproc := &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SEND)}
		assert.True(t, extproc.Equals(extproc), "extproc should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameExtproc := func() *ExtprocIR {
			return &ExtprocIR{perRoute: createSimpleExtproc(envoy_ext_proc_v3.ProcessingMode_SKIP)}
		}

		a := createSameExtproc()
		b := createSameExtproc()
		c := createSameExtproc()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestExtprocIRValidate(t *testing.T) {
	// Note: ExtProc validation is permissive at the protobuf level.
	// Provider validation is tested separately.
	tests := []struct {
		name    string
		extproc *ExtprocIR
	}{
		{
			name:    "nil extproc is valid",
			extproc: nil,
		},
		{
			name: "extproc with nil fields is valid",
			extproc: &ExtprocIR{
				perRoute: nil,
				provider: nil,
			},
		},
		{
			name: "valid extproc config passes validation",
			extproc: &ExtprocIR{
				perRoute: &envoy_ext_proc_v3.ExtProcPerRoute{
					Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{
						Overrides: &envoy_ext_proc_v3.ExtProcOverrides{
							ProcessingMode: &envoy_ext_proc_v3.ProcessingMode{
								RequestHeaderMode: envoy_ext_proc_v3.ProcessingMode_SEND,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.extproc.Validate()
			assert.NoError(t, err, "ExtProc validation should never fail at protobuf level")
		})
	}
}

func TestExtprocIRMergeInto(t *testing.T) {
	// Helper to create extproc configurations
	createExtprocConfig := func(mode envoy_ext_proc_v3.ProcessingMode_HeaderSendMode) *ExtprocIR {
		return &ExtprocIR{
			perRoute: &envoy_ext_proc_v3.ExtProcPerRoute{
				Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{
					Overrides: &envoy_ext_proc_v3.ExtProcOverrides{
						ProcessingMode: &envoy_ext_proc_v3.ProcessingMode{
							RequestHeaderMode: mode,
						},
					},
				},
			},
		}
	}

	tests := []struct {
		name            string
		p1Extproc       *ExtprocIR
		p2Extproc       *ExtprocIR
		strategy        policy.MergeOptions
		expectOriginSet bool
		expectMode      *envoy_ext_proc_v3.ProcessingMode_HeaderSendMode
	}{
		{
			name:            "shallow merge with nil p1",
			p1Extproc:       nil,
			p2Extproc:       createExtprocConfig(envoy_ext_proc_v3.ProcessingMode_SEND),
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet: true,
			expectMode:      &[]envoy_ext_proc_v3.ProcessingMode_HeaderSendMode{envoy_ext_proc_v3.ProcessingMode_SEND}[0],
		},
		{
			name:            "deep merge preserves p1 when set",
			p1Extproc:       createExtprocConfig(envoy_ext_proc_v3.ProcessingMode_SKIP),
			p2Extproc:       createExtprocConfig(envoy_ext_proc_v3.ProcessingMode_SEND),
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet: false,                                                                                        // No merge should happen
			expectMode:      &[]envoy_ext_proc_v3.ProcessingMode_HeaderSendMode{envoy_ext_proc_v3.ProcessingMode_SKIP}[0], // Should keep p1
		},
		{
			name:            "overridable shallow merge overwrites p1",
			p1Extproc:       createExtprocConfig(envoy_ext_proc_v3.ProcessingMode_SKIP),
			p2Extproc:       createExtprocConfig(envoy_ext_proc_v3.ProcessingMode_SEND),
			strategy:        policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet: true,
			expectMode:      &[]envoy_ext_proc_v3.ProcessingMode_HeaderSendMode{envoy_ext_proc_v3.ProcessingMode_SEND}[0],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test policies
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{extProc: tt.p1Extproc}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{extProc: tt.p2Extproc}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			// Perform merge
			if tt.p2Extproc != nil {
				tt.p2Extproc.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectMode != nil {
				assert.NotNil(t, p1.spec.extProc, "p1.extProc should not be nil")
				assert.NotNil(t, p1.spec.extProc.perRoute, "p1.extProc.perRoute should not be nil")

				// Verify processing mode
				overrides := p1.spec.extProc.perRoute.GetOverrides()
				if overrides != nil && overrides.ProcessingMode != nil {
					assert.Equal(t, *tt.expectMode, overrides.ProcessingMode.RequestHeaderMode)
				}
			}

			// Verify origin tracking
			if tt.expectOriginSet {
				origins := mergeOrigins.Get("extProc")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("extProc")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}

func TestBuildEnvoyExtProc(t *testing.T) {
	tests := []struct {
		name           string
		gatewayExt     *ir.GatewayExtension
		extprocConfig  *v1alpha1.ExtProcPolicy
		expectedError  string
		validateResult func(*testing.T, *envoy_ext_proc_v3.ExtProcPerRoute)
	}{
		{
			name: "with all processing modes",
			gatewayExt: &ir.GatewayExtension{
				ExtProc: &v1alpha1.ExtProcProvider{
					GrpcService: &v1alpha1.ExtGrpcService{},
				},
			},
			extprocConfig: &v1alpha1.ExtProcPolicy{
				ProcessingMode: &v1alpha1.ProcessingMode{
					RequestHeaderMode:   ptr.To("SEND"),
					ResponseHeaderMode:  ptr.To("SKIP"),
					RequestBodyMode:     ptr.To("STREAMED"),
					ResponseBodyMode:    ptr.To("BUFFERED"),
					RequestTrailerMode:  ptr.To("SEND"),
					ResponseTrailerMode: ptr.To("SKIP"),
				},
			},
			validateResult: func(t *testing.T, result *envoy_ext_proc_v3.ExtProcPerRoute) {
				processingMode := result.GetOverrides().GetProcessingMode()
				assert.NotNil(t, processingMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_SEND, processingMode.RequestHeaderMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_SKIP, processingMode.ResponseHeaderMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_STREAMED, processingMode.RequestBodyMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_BUFFERED, processingMode.ResponseBodyMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_SEND, processingMode.RequestTrailerMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_SKIP, processingMode.ResponseTrailerMode)
			},
		},
		{
			name: "with default processing modes",
			gatewayExt: &ir.GatewayExtension{
				ExtProc: &v1alpha1.ExtProcProvider{
					GrpcService: &v1alpha1.ExtGrpcService{},
				},
			},
			extprocConfig: &v1alpha1.ExtProcPolicy{
				ProcessingMode: &v1alpha1.ProcessingMode{},
			},
			validateResult: func(t *testing.T, result *envoy_ext_proc_v3.ExtProcPerRoute) {
				processingMode := result.GetOverrides().GetProcessingMode()
				assert.NotNil(t, processingMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_DEFAULT, processingMode.RequestHeaderMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_DEFAULT, processingMode.ResponseHeaderMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_NONE, processingMode.RequestBodyMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_NONE, processingMode.ResponseBodyMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_DEFAULT, processingMode.RequestTrailerMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_DEFAULT, processingMode.ResponseTrailerMode)
			},
		},
		{
			name: "with invalid processing modes",
			gatewayExt: &ir.GatewayExtension{
				ExtProc: &v1alpha1.ExtProcProvider{
					GrpcService: &v1alpha1.ExtGrpcService{},
				},
			},
			extprocConfig: &v1alpha1.ExtProcPolicy{
				ProcessingMode: &v1alpha1.ProcessingMode{
					RequestHeaderMode:   ptr.To("INVALID"),
					ResponseHeaderMode:  ptr.To("INVALID"),
					RequestBodyMode:     ptr.To("INVALID"),
					ResponseBodyMode:    ptr.To("INVALID"),
					RequestTrailerMode:  ptr.To("INVALID"),
					ResponseTrailerMode: ptr.To("INVALID"),
				},
			},
			validateResult: func(t *testing.T, result *envoy_ext_proc_v3.ExtProcPerRoute) {
				processingMode := result.GetOverrides().GetProcessingMode()
				assert.NotNil(t, processingMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_DEFAULT, processingMode.RequestHeaderMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_DEFAULT, processingMode.ResponseHeaderMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_NONE, processingMode.RequestBodyMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_NONE, processingMode.ResponseBodyMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_DEFAULT, processingMode.RequestTrailerMode)
				assert.Equal(t, envoy_ext_proc_v3.ProcessingMode_DEFAULT, processingMode.ResponseTrailerMode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateExtProcPerFilterConfig(tt.extprocConfig)

			//require.NoError(t, err)
			require.NotNil(t, result)
			tt.validateResult(t, result)
		})
	}
}
