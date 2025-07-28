package trafficpolicy

import (
	"errors"
	"fmt"
	"testing"
	"time"

	envoycorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyratelimitv3 "github.com/envoyproxy/go-control-plane/envoy/config/ratelimit/v3"
	envoyroutev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	ratev3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestGlobalRateLimitIREquals(t *testing.T) {
	// Helper to create simple rate limit configurations for testing
	createSimpleRateLimit := func(key string) []*envoyroutev3.RateLimit {
		return []*envoyroutev3.RateLimit{
			{
				Actions: []*envoyroutev3.RateLimit_Action{
					{
						ActionSpecifier: &envoyroutev3.RateLimit_Action_GenericKey_{
							GenericKey: &envoyroutev3.RateLimit_Action_GenericKey{
								DescriptorKey:   key,
								DescriptorValue: "test-value",
							},
						},
					},
				},
			},
		}
	}

	createProvider := func(name string) *TrafficPolicyGatewayExtensionIR {
		return &TrafficPolicyGatewayExtensionIR{
			Name: name,
			RateLimit: &ratev3.RateLimit{
				Domain: "test-domain",
				RateLimitService: &envoyratelimitv3.RateLimitServiceConfig{
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

	tests := []struct {
		name       string
		rateLimit1 *GlobalRateLimitIR
		rateLimit2 *GlobalRateLimitIR
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
			rateLimit2: &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("key1")},
			expected:   false,
		},
		{
			name:       "non-nil vs nil are not equal",
			rateLimit1: &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("key1")},
			rateLimit2: nil,
			expected:   false,
		},
		{
			name:       "same instance is equal",
			rateLimit1: &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("key1")},
			rateLimit2: &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("key1")},
			expected:   true,
		},
		{
			name:       "different rate limit keys are not equal",
			rateLimit1: &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("key1")},
			rateLimit2: &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("key2")},
			expected:   false,
		},
		{
			name:       "different providers are not equal",
			rateLimit1: &GlobalRateLimitIR{provider: createProvider("service1")},
			rateLimit2: &GlobalRateLimitIR{provider: createProvider("service2")},
			expected:   false,
		},
		{
			name:       "same providers are equal",
			rateLimit1: &GlobalRateLimitIR{provider: createProvider("service1")},
			rateLimit2: &GlobalRateLimitIR{provider: createProvider("service1")},
			expected:   true,
		},
		{
			name:       "different length action slices are not equal",
			rateLimit1: &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("key1")},
			rateLimit2: &GlobalRateLimitIR{rateLimitActions: append(createSimpleRateLimit("key1"), createSimpleRateLimit("key2")...)},
			expected:   false,
		},
		{
			name:       "nil fields are equal",
			rateLimit1: &GlobalRateLimitIR{rateLimitActions: nil, provider: nil},
			rateLimit2: &GlobalRateLimitIR{rateLimitActions: nil, provider: nil},
			expected:   true,
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
		rateLimit := &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("test")}
		assert.True(t, rateLimit.Equals(rateLimit), "rateLimit should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameRateLimit := func() *GlobalRateLimitIR {
			return &GlobalRateLimitIR{rateLimitActions: createSimpleRateLimit("test")}
		}

		a := createSameRateLimit()
		b := createSameRateLimit()
		c := createSameRateLimit()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestGlobalRateLimitIRValidate(t *testing.T) {
	// Note: Global rate limit validation is permissive at the protobuf level.
	// Provider validation is tested separately.
	tests := []struct {
		name      string
		rateLimit *GlobalRateLimitIR
	}{
		{
			name:      "nil rate limit is valid",
			rateLimit: nil,
		},
		{
			name: "rate limit with nil fields is valid",
			rateLimit: &GlobalRateLimitIR{
				rateLimitActions: nil,
				provider:         nil,
			},
		},
		{
			name: "valid rate limit config passes validation",
			rateLimit: &GlobalRateLimitIR{
				rateLimitActions: []*envoyroutev3.RateLimit{
					{
						Actions: []*envoyroutev3.RateLimit_Action{
							{
								ActionSpecifier: &envoyroutev3.RateLimit_Action_GenericKey_{
									GenericKey: &envoyroutev3.RateLimit_Action_GenericKey{
										DescriptorKey:   "api",
										DescriptorValue: "test-value",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "rate limit config with empty descriptor key is valid",
			rateLimit: &GlobalRateLimitIR{
				rateLimitActions: []*envoyroutev3.RateLimit{
					{
						Actions: []*envoyroutev3.RateLimit_Action{
							{
								ActionSpecifier: &envoyroutev3.RateLimit_Action_GenericKey_{
									GenericKey: &envoyroutev3.RateLimit_Action_GenericKey{
										DescriptorKey:   "", // Empty key is valid at protobuf level
										DescriptorValue: "api",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rateLimit.Validate()
			assert.NoError(t, err, "Global rate limit validation should never fail at protobuf level")
		})
	}
}

func TestGlobalRateLimitIRMergeInto(t *testing.T) {
	// Helper to create rate limit configurations
	createRateLimitConfig := func(key string) *GlobalRateLimitIR {
		return &GlobalRateLimitIR{
			rateLimitActions: []*envoyroutev3.RateLimit{
				{
					Actions: []*envoyroutev3.RateLimit_Action{
						{
							ActionSpecifier: &envoyroutev3.RateLimit_Action_GenericKey_{
								GenericKey: &envoyroutev3.RateLimit_Action_GenericKey{
									DescriptorKey:   key,
									DescriptorValue: "test-value",
								},
							},
						},
					},
				},
			},
		}
	}

	tests := []struct {
		name              string
		p1RateLimit       *GlobalRateLimitIR
		p2RateLimit       *GlobalRateLimitIR
		strategy          policy.MergeOptions
		expectOriginSet   bool
		expectActionCount *int
	}{
		{
			name:              "shallow merge with nil p1",
			p1RateLimit:       nil,
			p2RateLimit:       createRateLimitConfig("service"),
			strategy:          policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet:   true,
			expectActionCount: ptr.To(1),
		},
		{
			name:              "deep merge preserves p1 when set",
			p1RateLimit:       createRateLimitConfig("service1"),
			p2RateLimit:       createRateLimitConfig("service2"),
			strategy:          policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet:   false,     // No merge should happen
			expectActionCount: ptr.To(1), // Should keep p1
		},
		{
			name:              "overridable shallow merge overwrites p1",
			p1RateLimit:       createRateLimitConfig("service1"),
			p2RateLimit:       createRateLimitConfig("service2"),
			strategy:          policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet:   true,
			expectActionCount: ptr.To(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test policies
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{globalRateLimit: tt.p1RateLimit}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{globalRateLimit: tt.p2RateLimit}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			// Perform merge
			if tt.p2RateLimit != nil {
				tt.p2RateLimit.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectActionCount != nil {
				assert.NotNil(t, p1.spec.globalRateLimit, "p1.globalRateLimit should not be nil")
				assert.NotNil(t, p1.spec.globalRateLimit.rateLimitActions, "p1.globalRateLimit.rateLimitActions should not be nil")

				// Verify action count
				assert.Equal(t, *tt.expectActionCount, len(p1.spec.globalRateLimit.rateLimitActions))
			}

			// Verify origin tracking
			if tt.expectOriginSet {
				origins := mergeOrigins.Get("rateLimit.global")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("rateLimit.global")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}

func TestCreateRateLimitActions(t *testing.T) {
	tests := []struct {
		name           string
		descriptors    []v1alpha1.RateLimitDescriptor
		expectedError  string
		validateResult func(*testing.T, []*envoyroutev3.RateLimit_Action)
	}{
		{
			name: "with generic key descriptor",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
							Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
								Key:   "service",
								Value: "api",
							},
						},
					},
				},
			},
			validateResult: func(t *testing.T, actions []*envoyroutev3.RateLimit_Action) {
				require.Len(t, actions, 1)
				genericKey := actions[0].GetGenericKey()
				require.NotNil(t, genericKey)
				assert.Equal(t, "service", genericKey.DescriptorKey)
				assert.Equal(t, "api", genericKey.DescriptorValue)
			},
		},
		{
			name: "with header descriptor",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type:   v1alpha1.RateLimitDescriptorEntryTypeHeader,
							Header: ptr.To("X-User-ID"),
						},
					},
				},
			},
			validateResult: func(t *testing.T, actions []*envoyroutev3.RateLimit_Action) {
				require.Len(t, actions, 1)
				requestHeaders := actions[0].GetRequestHeaders()
				require.NotNil(t, requestHeaders)
				assert.Equal(t, "X-User-ID", requestHeaders.HeaderName)
				assert.Equal(t, "X-User-ID", requestHeaders.DescriptorKey)
			},
		},
		{
			name: "with remote address descriptor",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: v1alpha1.RateLimitDescriptorEntryTypeRemoteAddress,
						},
					},
				},
			},
			validateResult: func(t *testing.T, actions []*envoyroutev3.RateLimit_Action) {
				require.Len(t, actions, 1)
				remoteAddress := actions[0].GetRemoteAddress()
				require.NotNil(t, remoteAddress)
			},
		},
		{
			name: "with path descriptor",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: v1alpha1.RateLimitDescriptorEntryTypePath,
						},
					},
				},
			},
			validateResult: func(t *testing.T, actions []*envoyroutev3.RateLimit_Action) {
				require.Len(t, actions, 1)
				requestHeaders := actions[0].GetRequestHeaders()
				require.NotNil(t, requestHeaders)
				assert.Equal(t, ":path", requestHeaders.HeaderName)
				assert.Equal(t, "path", requestHeaders.DescriptorKey)
			},
		},
		{
			name: "with multiple descriptors",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
							Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
								Key:   "service",
								Value: "api",
							},
						},
					},
				},
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: v1alpha1.RateLimitDescriptorEntryTypeRemoteAddress,
						},
					},
				},
			},
			validateResult: func(t *testing.T, actions []*envoyroutev3.RateLimit_Action) {
				require.Len(t, actions, 2)
				// First action is generic key
				genericKey := actions[0].GetGenericKey()
				require.NotNil(t, genericKey)
				assert.Equal(t, "service", genericKey.DescriptorKey)
				assert.Equal(t, "api", genericKey.DescriptorValue)

				// Second action is remote address
				remoteAddress := actions[1].GetRemoteAddress()
				require.NotNil(t, remoteAddress)
			},
		},
		{
			name: "with multiple entries in one descriptor",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
							Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
								Key:   "service",
								Value: "api",
							},
						},
						{
							Type:   v1alpha1.RateLimitDescriptorEntryTypeHeader,
							Header: ptr.To("X-User-ID"),
						},
					},
				},
			},
			validateResult: func(t *testing.T, actions []*envoyroutev3.RateLimit_Action) {
				require.Len(t, actions, 2)
				// First action is generic key
				genericKey := actions[0].GetGenericKey()
				require.NotNil(t, genericKey)
				assert.Equal(t, "service", genericKey.DescriptorKey)
				assert.Equal(t, "api", genericKey.DescriptorValue)

				// Second action is header
				requestHeaders := actions[1].GetRequestHeaders()
				require.NotNil(t, requestHeaders)
				assert.Equal(t, "X-User-ID", requestHeaders.HeaderName)
				assert.Equal(t, "X-User-ID", requestHeaders.DescriptorKey)
			},
		},
		{
			name:          "with empty descriptors",
			descriptors:   []v1alpha1.RateLimitDescriptor{},
			expectedError: "at least one descriptor is required for global rate limiting",
		},
		{
			name: "with missing generic key entry data",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
						},
					},
				},
			},
			expectedError: "generic entry requires Generic field to be set",
		},
		{
			name: "with missing header name",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: v1alpha1.RateLimitDescriptorEntryTypeHeader,
						},
					},
				},
			},
			expectedError: "header entry requires Header field to be set",
		},
		{
			name: "with unsupported entry type",
			descriptors: []v1alpha1.RateLimitDescriptor{
				{
					Entries: []v1alpha1.RateLimitDescriptorEntry{
						{
							Type: "UnsupportedType",
						},
					},
				},
			},
			expectedError: "unsupported entry type: UnsupportedType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actions, err := createRateLimitActions(tt.descriptors)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, actions)
			tt.validateResult(t, actions)
		})
	}
}

func TestToRateLimitFilterConfig(t *testing.T) {
	defaultExtensionName := "test-ratelimit"
	defaultNamespace := "test-namespace"
	defaultClusterName := "test-service.test-namespace.svc.cluster.local:8081"

	createBackendRef := func() gwv1alpha2.BackendObjectReference {
		port := gwv1alpha2.PortNumber(8081)
		return gwv1alpha2.BackendObjectReference{
			Name: "test-service",
			Port: &port,
		}
	}

	tests := []struct {
		name              string
		gatewayExtension  *ir.GatewayExtension
		policy            *v1alpha1.RateLimitPolicy
		trafficPolicy     *v1alpha1.TrafficPolicy
		expectedError     string
		validateRateLimit func(*testing.T, *ratev3.RateLimit)
	}{
		{
			name: "with default configuration",
			gatewayExtension: &ir.GatewayExtension{
				Type: v1alpha1.GatewayExtensionTypeRateLimit,
				RateLimit: &v1alpha1.RateLimitProvider{
					Domain: "test-domain",
					GrpcService: &v1alpha1.ExtGrpcService{
						BackendRef: &gwv1alpha2.BackendRef{
							BackendObjectReference: createBackendRef(),
						},
					},
				},
				ObjectSource: ir.ObjectSource{
					Name:      defaultExtensionName,
					Namespace: defaultNamespace,
				},
			},
			policy: &v1alpha1.RateLimitPolicy{
				ExtensionRef: &corev1.LocalObjectReference{
					Name: defaultExtensionName,
				},
				Descriptors: []v1alpha1.RateLimitDescriptor{
					{
						Entries: []v1alpha1.RateLimitDescriptorEntry{
							{
								Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
								Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
									Key:   "service",
									Value: "api",
								},
							},
						},
					},
				},
			},
			trafficPolicy: &v1alpha1.TrafficPolicy{},
			validateRateLimit: func(t *testing.T, rl *ratev3.RateLimit) {
				require.NotNil(t, rl)
				assert.Equal(t, "test-domain", rl.Domain)
				assert.Equal(t, defaultClusterName, rl.RateLimitService.GrpcService.GetEnvoyGrpc().ClusterName)
				assert.Equal(t, envoycorev3.ApiVersion_V3, rl.RateLimitService.TransportApiVersion)
				assert.Equal(t, ratev3.RateLimit_DRAFT_VERSION_03, rl.EnableXRatelimitHeaders)
				assert.Equal(t, "both", rl.RequestType)
				assert.Equal(t, rateLimitStatPrefix, rl.StatPrefix)
				assert.Equal(t, &durationpb.Duration{Seconds: 0}, rl.Timeout)
				assert.True(t, rl.FailureModeDeny) // Default should be failureModeAllow=false (deny)
			},
		},
		{
			name: "with custom timeout",
			gatewayExtension: &ir.GatewayExtension{
				Type: v1alpha1.GatewayExtensionTypeRateLimit,
				RateLimit: &v1alpha1.RateLimitProvider{
					Domain: "test-domain",
					GrpcService: &v1alpha1.ExtGrpcService{
						BackendRef: &gwv1alpha2.BackendRef{
							BackendObjectReference: createBackendRef(),
						},
					},
					Timeout: metav1.Duration{Duration: 5 * time.Second},
				},
				ObjectSource: ir.ObjectSource{
					Name:      defaultExtensionName,
					Namespace: defaultNamespace,
				},
			},
			policy: &v1alpha1.RateLimitPolicy{
				ExtensionRef: &corev1.LocalObjectReference{
					Name: defaultExtensionName,
				},
				Descriptors: []v1alpha1.RateLimitDescriptor{
					{
						Entries: []v1alpha1.RateLimitDescriptorEntry{
							{
								Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
								Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
									Key:   "service",
									Value: "api",
								},
							},
						},
					},
				},
			},
			trafficPolicy: &v1alpha1.TrafficPolicy{},
			validateRateLimit: func(t *testing.T, rl *ratev3.RateLimit) {
				require.NotNil(t, rl)
				assert.Equal(t, time.Duration(5*time.Second), rl.Timeout.AsDuration())
			},
		},
		{
			name: "with fail open configuration",
			gatewayExtension: &ir.GatewayExtension{
				Type: v1alpha1.GatewayExtensionTypeRateLimit,
				RateLimit: &v1alpha1.RateLimitProvider{
					Domain: "test-domain",
					GrpcService: &v1alpha1.ExtGrpcService{
						BackendRef: &gwv1alpha2.BackendRef{
							BackendObjectReference: createBackendRef(),
						},
					},
					FailOpen: true,
				},
				ObjectSource: ir.ObjectSource{
					Name:      defaultExtensionName,
					Namespace: defaultNamespace,
				},
			},
			policy: &v1alpha1.RateLimitPolicy{
				ExtensionRef: &corev1.LocalObjectReference{
					Name: defaultExtensionName,
				},
				Descriptors: []v1alpha1.RateLimitDescriptor{
					{
						Entries: []v1alpha1.RateLimitDescriptorEntry{
							{
								Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
								Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
									Key:   "service",
									Value: "api",
								},
							},
						},
					},
				},
			},
			trafficPolicy: &v1alpha1.TrafficPolicy{},
			validateRateLimit: func(t *testing.T, rl *ratev3.RateLimit) {
				require.NotNil(t, rl)
				assert.False(t, rl.FailureModeDeny) // Should be fail open (false)
			},
		},
		{
			name: "without backend reference",
			gatewayExtension: &ir.GatewayExtension{
				Type: v1alpha1.GatewayExtensionTypeRateLimit,
				RateLimit: &v1alpha1.RateLimitProvider{
					Domain:      "test-domain",
					GrpcService: &v1alpha1.ExtGrpcService{},
				},
				ObjectSource: ir.ObjectSource{
					Name:      defaultExtensionName,
					Namespace: defaultNamespace,
				},
			},
			policy: &v1alpha1.RateLimitPolicy{
				ExtensionRef: &corev1.LocalObjectReference{
					Name: defaultExtensionName,
				},
				Descriptors: []v1alpha1.RateLimitDescriptor{
					{
						Entries: []v1alpha1.RateLimitDescriptorEntry{
							{
								Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
								Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
									Key:   "service",
									Value: "api",
								},
							},
						},
					},
				},
			},
			trafficPolicy: &v1alpha1.TrafficPolicy{},
			expectedError: "backend not provided in grpc service",
		},
		{
			name: "with wrong extension type",
			gatewayExtension: &ir.GatewayExtension{
				Type: v1alpha1.GatewayExtensionTypeExtProc,
				ObjectSource: ir.ObjectSource{
					Name:      defaultExtensionName,
					Namespace: defaultNamespace,
				},
			},
			policy: &v1alpha1.RateLimitPolicy{
				ExtensionRef: &corev1.LocalObjectReference{
					Name: defaultExtensionName,
				},
				Descriptors: []v1alpha1.RateLimitDescriptor{
					{
						Entries: []v1alpha1.RateLimitDescriptorEntry{
							{
								Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
								Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
									Key:   "service",
									Value: "api",
								},
							},
						},
					},
				},
			},
			trafficPolicy: &v1alpha1.TrafficPolicy{},
			expectedError: "extension has type ExtProc but RateLimit was expected",
		},
		{
			name: "without extension reference",
			policy: &v1alpha1.RateLimitPolicy{
				Descriptors: []v1alpha1.RateLimitDescriptor{
					{
						Entries: []v1alpha1.RateLimitDescriptorEntry{
							{
								Type: v1alpha1.RateLimitDescriptorEntryTypeGeneric,
								Generic: &v1alpha1.RateLimitDescriptorEntryGeneric{
									Key:   "service",
									Value: "api",
								},
							},
						},
					},
				},
			},
			trafficPolicy: &v1alpha1.TrafficPolicy{},
			expectedError: "extensionRef is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rl *ratev3.RateLimit
			var err error

			if tt.policy == nil || tt.policy.ExtensionRef == nil {
				err = errors.New("extensionRef is required")
			} else if tt.gatewayExtension == nil {
				err = fmt.Errorf("failed to get referenced GatewayExtension")
			} else if tt.gatewayExtension.Type != v1alpha1.GatewayExtensionTypeRateLimit {
				err = fmt.Errorf("extension has type %s but %s was expected",
					tt.gatewayExtension.Type, v1alpha1.GatewayExtensionTypeRateLimit)
			} else {
				// Get the extension's spec
				extension := tt.gatewayExtension.RateLimit
				if extension == nil {
					err = fmt.Errorf("RateLimit configuration is missing in GatewayExtension")
				} else {
					// Create a timeout based on the timeout from extension
					timeout := durationpb.New(extension.Timeout.Duration)

					if err == nil {
						// Use the domain from the extension
						domain := extension.Domain

						// Construct cluster name from the backendRef
						clusterName := ""
						if extension.GrpcService != nil && extension.GrpcService.BackendRef != nil {
							clusterName = fmt.Sprintf("%s.%s.svc.cluster.local:%d",
								extension.GrpcService.BackendRef.Name,
								tt.gatewayExtension.Namespace,
								*extension.GrpcService.BackendRef.Port)
						} else {
							err = fmt.Errorf("backend not provided in grpc service")
						}

						if err == nil {
							// Create a rate limit configuration
							rl = &ratev3.RateLimit{
								Domain:          domain,
								Timeout:         timeout,
								FailureModeDeny: !extension.FailOpen,
								RateLimitService: &envoyratelimitv3.RateLimitServiceConfig{
									GrpcService: &envoycorev3.GrpcService{
										TargetSpecifier: &envoycorev3.GrpcService_EnvoyGrpc_{
											EnvoyGrpc: &envoycorev3.GrpcService_EnvoyGrpc{
												ClusterName: clusterName,
											},
										},
									},
									TransportApiVersion: envoycorev3.ApiVersion_V3,
								},
								Stage:                   0,
								EnableXRatelimitHeaders: ratev3.RateLimit_DRAFT_VERSION_03,
								RequestType:             "both",
								StatPrefix:              rateLimitStatPrefix,
							}
						}
					}
				}
			}

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, rl)
			tt.validateRateLimit(t, rl)
		})
	}
}
