package trafficpolicy

import (
	"testing"

	exteniondynamicmodulev3 "github.com/envoyproxy/go-control-plane/envoy/extensions/dynamic_modules/v3"
	dynamicmodulesv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/dynamic_modules/v3"
	transformationpb "github.com/solo-io/envoy-gloo/go/config/filter/http/transformation/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/utils/ptr"

	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestTransformationIREquals(t *testing.T) {
	// Helper to create simple transformations for testing
	createSimpleTransformation := func() *transformationpb.RouteTransformations {
		// Create a basic transformation using the actual structure
		return &transformationpb.RouteTransformations{
			Transformations: []*transformationpb.RouteTransformations_RouteTransformation{
				{
					Match: &transformationpb.RouteTransformations_RouteTransformation_RequestMatch_{
						RequestMatch: &transformationpb.RouteTransformations_RouteTransformation_RequestMatch{
							RequestTransformation: &transformationpb.Transformation{
								TransformationType: &transformationpb.Transformation_TransformationTemplate{
									TransformationTemplate: &transformationpb.TransformationTemplate{
										Headers: map[string]*transformationpb.InjaTemplate{
											"x-test": {Text: "test-value"},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	tests := []struct {
		name     string
		trans1   *TransformationIR
		trans2   *TransformationIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			trans1:   nil,
			trans2:   nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			trans1:   nil,
			trans2:   &TransformationIR{transformation: createSimpleTransformation()},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			trans1:   &TransformationIR{transformation: createSimpleTransformation()},
			trans2:   nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			trans1:   &TransformationIR{transformation: createSimpleTransformation()},
			trans2:   &TransformationIR{transformation: createSimpleTransformation()},
			expected: true,
		},
		{
			name:     "nil transformation fields are equal",
			trans1:   &TransformationIR{transformation: nil},
			trans2:   &TransformationIR{transformation: nil},
			expected: true,
		},
		{
			name:     "nil vs non-nil transformation fields are not equal",
			trans1:   &TransformationIR{transformation: nil},
			trans2:   &TransformationIR{transformation: createSimpleTransformation()},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.trans1.Equals(tt.trans2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.trans2.Equals(tt.trans1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		transformation := &TransformationIR{
			transformation: &transformationpb.RouteTransformations{
				Transformations: []*transformationpb.RouteTransformations_RouteTransformation{},
			},
		}
		assert.True(t, transformation.Equals(transformation), "transformation should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameTransformation := func() *TransformationIR {
			return &TransformationIR{
				transformation: createSimpleTransformation(),
			}
		}

		a := createSameTransformation()
		b := createSameTransformation()
		c := createSameTransformation()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestTransformationIRValidate(t *testing.T) {
	// Note: Transformation validation is permissive at the protobuf level.
	// Complex validation logic is handled during translation.
	tests := []struct {
		name  string
		trans *TransformationIR
	}{
		{
			name:  "nil transformation is valid",
			trans: nil,
		},
		{
			name: "transformation with nil field is valid",
			trans: &TransformationIR{
				transformation: nil,
			},
		},
		{
			name: "empty transformation is valid",
			trans: &TransformationIR{
				transformation: &transformationpb.RouteTransformations{
					Transformations: []*transformationpb.RouteTransformations_RouteTransformation{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.trans.Validate()
			assert.NoError(t, err, "Transformation validation should never fail at protobuf level")
		})
	}
}

func TestTransformationIRMergeInto(t *testing.T) {
	// Simple helper to create a basic transformation
	createTransformation := func(headerValue string) *TransformationIR {
		return &TransformationIR{
			transformation: &transformationpb.RouteTransformations{
				Transformations: []*transformationpb.RouteTransformations_RouteTransformation{
					{
						Match: &transformationpb.RouteTransformations_RouteTransformation_RequestMatch_{
							RequestMatch: &transformationpb.RouteTransformations_RouteTransformation_RequestMatch{
								RequestTransformation: &transformationpb.Transformation{
									TransformationType: &transformationpb.Transformation_TransformationTemplate{
										TransformationTemplate: &transformationpb.TransformationTemplate{
											Headers: map[string]*transformationpb.InjaTemplate{
												"x-test": {Text: headerValue},
											},
										},
									},
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
		p1Trans           *TransformationIR
		p2Trans           *TransformationIR
		strategy          policy.MergeOptions
		expectOriginSet   bool
		expectHeaderValue *string
	}{
		{
			name:              "shallow merge with nil p1",
			p1Trans:           nil,
			p2Trans:           createTransformation("p2-value"),
			strategy:          policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet:   true,
			expectHeaderValue: ptr.To("p2-value"),
		},
		{
			name:              "deep merge preserves p1 when set",
			p1Trans:           createTransformation("p1-value"),
			p2Trans:           createTransformation("p2-value"),
			strategy:          policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet:   true,               // In transformation, deep merge actually merges
			expectHeaderValue: ptr.To("p1-value"), // Should preserve p1 content after merge
		},
		{
			name:              "overridable shallow merge overwrites p1",
			p1Trans:           createTransformation("p1-value"),
			p2Trans:           createTransformation("p2-value"),
			strategy:          policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet:   true,
			expectHeaderValue: ptr.To("p2-value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test policies
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{transformation: tt.p1Trans}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{transformation: tt.p2Trans}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			// Perform merge
			if tt.p2Trans != nil {
				tt.p2Trans.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectHeaderValue != nil {
				assert.NotNil(t, p1.spec.transformation, "p1.transformation should not be nil")
				assert.NotNil(t, p1.spec.transformation.transformation, "p1.transformation.transformation should not be nil")

				// Verify the actual header value in the transformation
				if len(p1.spec.transformation.transformation.Transformations) > 0 {
					trans := p1.spec.transformation.transformation.Transformations[0]
					if reqMatch := trans.GetRequestMatch(); reqMatch != nil {
						if headers := reqMatch.RequestTransformation.GetTransformationTemplate().Headers; headers != nil {
							if testHeader := headers["x-test"]; testHeader != nil {
								assert.Equal(t, *tt.expectHeaderValue, testHeader.Text)
							}
						}
					}
				}
			}

			// Verify origin tracking
			if tt.expectOriginSet {
				origins := mergeOrigins.Get("transformation")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("transformation")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}

func TestRustformationIREquals(t *testing.T) {
	// Helper to create simple rustformation configurations for testing
	createSimpleRustformation := func() *dynamicmodulesv3.DynamicModuleFilter {
		// Create a basic dynamic module filter for rustformation
		filterConfig, _ := anypb.New(&exteniondynamicmodulev3.DynamicModuleConfig{
			Name: "rust_module",
		})

		return &dynamicmodulesv3.DynamicModuleFilter{
			DynamicModuleConfig: &exteniondynamicmodulev3.DynamicModuleConfig{
				Name: "rust_module",
			},
			FilterName:   "http_simple_mutations",
			FilterConfig: filterConfig,
		}
	}

	tests := []struct {
		name     string
		rust1    *RustformationIR
		rust2    *RustformationIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			rust1:    nil,
			rust2:    nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			rust1:    nil,
			rust2:    &RustformationIR{rustformation: createSimpleRustformation()},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			rust1:    &RustformationIR{rustformation: createSimpleRustformation()},
			rust2:    nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			rust1:    &RustformationIR{rustformation: createSimpleRustformation()},
			rust2:    &RustformationIR{rustformation: createSimpleRustformation()},
			expected: true,
		},
		{
			name:     "nil rustformation fields are equal",
			rust1:    &RustformationIR{rustformation: nil},
			rust2:    &RustformationIR{rustformation: nil},
			expected: true,
		},
		{
			name:     "nil vs non-nil rustformation fields are not equal",
			rust1:    &RustformationIR{rustformation: nil},
			rust2:    &RustformationIR{rustformation: createSimpleRustformation()},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rust1.Equals(tt.rust2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.rust2.Equals(tt.rust1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		rustformation := &RustformationIR{
			rustformation: &dynamicmodulesv3.DynamicModuleFilter{
				FilterName: "http_simple_mutations",
			},
		}
		assert.True(t, rustformation.Equals(rustformation), "rustformation should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		createSameRustformation := func() *RustformationIR {
			return &RustformationIR{
				rustformation: createSimpleRustformation(),
			}
		}

		a := createSameRustformation()
		b := createSameRustformation()
		c := createSameRustformation()

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestRustformationIRValidate(t *testing.T) {
	// Note: Rustformation validation is permissive at the protobuf level.
	tests := []struct {
		name string
		rust *RustformationIR
	}{
		{
			name: "nil rustformation is valid",
			rust: nil,
		},
		{
			name: "rustformation with nil field is valid",
			rust: &RustformationIR{
				rustformation: nil,
			},
		},
		{
			name: "empty rustformation is valid",
			rust: &RustformationIR{
				rustformation: &dynamicmodulesv3.DynamicModuleFilter{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rust.Validate()
			assert.NoError(t, err, "Rustformation validation should never fail at protobuf level")
		})
	}
}

func TestRustformationIRMergeInto(t *testing.T) {
	// Helper to create rustformation configurations with different module names
	createRustformationConfig := func(moduleName string) *RustformationIR {
		return &RustformationIR{
			rustformation: &dynamicmodulesv3.DynamicModuleFilter{
				DynamicModuleConfig: &exteniondynamicmodulev3.DynamicModuleConfig{
					Name: moduleName,
				},
				FilterName: "http_simple_mutations",
			},
		}
	}

	tests := []struct {
		name             string
		p1Rust           *RustformationIR
		p2Rust           *RustformationIR
		strategy         policy.MergeOptions
		expectOriginSet  bool
		expectModuleName *string
	}{
		{
			name:             "shallow merge with nil p1",
			p1Rust:           nil,
			p2Rust:           createRustformationConfig("rust_module_p2"),
			strategy:         policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectOriginSet:  true,
			expectModuleName: ptr.To("rust_module_p2"),
		},
		{
			name:             "deep merge preserves p1 when set",
			p1Rust:           createRustformationConfig("rust_module_p1"),
			p2Rust:           createRustformationConfig("rust_module_p2"),
			strategy:         policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectOriginSet:  false,                    // No merge should happen
			expectModuleName: ptr.To("rust_module_p1"), // Should preserve p1
		},
		{
			name:             "overridable shallow merge overwrites p1",
			p1Rust:           createRustformationConfig("rust_module_p1"),
			p2Rust:           createRustformationConfig("rust_module_p2"),
			strategy:         policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectOriginSet:  true,
			expectModuleName: ptr.To("rust_module_p2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test policies
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{rustformation: tt.p1Rust}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{rustformation: tt.p2Rust}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			// Perform merge
			if tt.p2Rust != nil {
				tt.p2Rust.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Validate actual content instead of just nil checking
			if tt.expectModuleName != nil {
				assert.NotNil(t, p1.spec.rustformation, "p1.rustformation should not be nil")
				assert.NotNil(t, p1.spec.rustformation.rustformation, "p1.rustformation.rustformation should not be nil")

				// Verify the actual module name in the rustformation
				if p1.spec.rustformation.rustformation.DynamicModuleConfig != nil {
					assert.Equal(t, *tt.expectModuleName, p1.spec.rustformation.rustformation.DynamicModuleConfig.Name)
				}
			}

			// Verify origin tracking
			if tt.expectOriginSet {
				origins := mergeOrigins.Get("rustformation")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("rustformation")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}
