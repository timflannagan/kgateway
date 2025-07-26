package trafficpolicy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
)

func TestBufferForSpec(t *testing.T) {
	tests := []struct {
		name     string
		spec     v1alpha1.TrafficPolicySpec
		expected *BufferIR
	}{
		{
			name:     "nil buffer spec",
			spec:     v1alpha1.TrafficPolicySpec{},
			expected: nil,
		},
		{
			name: "valid buffer spec",
			spec: v1alpha1.TrafficPolicySpec{
				Buffer: &v1alpha1.Buffer{
					MaxRequestSize: ptr.To(resource.MustParse("1Ki")),
				},
			},
			expected: &BufferIR{
				maxRequestBytes: 1024,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &trafficPolicySpecIr{}
			bufferForSpec(tt.spec, out)

			if tt.expected == nil {
				assert.Nil(t, out.buffer)
			} else {
				assert.NotNil(t, out.buffer)
				assert.Equal(t, tt.expected.maxRequestBytes, out.buffer.maxRequestBytes)
			}
		})
	}
}

func TestBufferIREquals(t *testing.T) {
	tests := []struct {
		name     string
		buffer1  *BufferIR
		buffer2  *BufferIR
		expected bool
	}{
		{
			name:     "both nil are equal",
			buffer1:  nil,
			buffer2:  nil,
			expected: true,
		},
		{
			name:     "nil vs non-nil are not equal",
			buffer1:  nil,
			buffer2:  &BufferIR{maxRequestBytes: 1024},
			expected: false,
		},
		{
			name:     "non-nil vs nil are not equal",
			buffer1:  &BufferIR{maxRequestBytes: 1024},
			buffer2:  nil,
			expected: false,
		},
		{
			name:     "same instance is equal",
			buffer1:  &BufferIR{maxRequestBytes: 1024},
			buffer2:  &BufferIR{maxRequestBytes: 1024},
			expected: true,
		},
		{
			name:     "same values are equal",
			buffer1:  &BufferIR{maxRequestBytes: 2048},
			buffer2:  &BufferIR{maxRequestBytes: 2048},
			expected: true,
		},
		{
			name:     "different values are not equal",
			buffer1:  &BufferIR{maxRequestBytes: 1024},
			buffer2:  &BufferIR{maxRequestBytes: 2048},
			expected: false,
		},
		{
			name:     "zero values are equal",
			buffer1:  &BufferIR{maxRequestBytes: 0},
			buffer2:  &BufferIR{maxRequestBytes: 0},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.buffer1.Equals(tt.buffer2)
			assert.Equal(t, tt.expected, result)

			// Test symmetry: a.Equals(b) should equal b.Equals(a)
			reverseResult := tt.buffer2.Equals(tt.buffer1)
			assert.Equal(t, result, reverseResult, "Equals should be symmetric")
		})
	}

	// Test reflexivity: x.Equals(x) should always be true for non-nil values
	t.Run("reflexivity", func(t *testing.T) {
		buffer := &BufferIR{maxRequestBytes: 1024}
		assert.True(t, buffer.Equals(buffer), "buffer should equal itself")
	})

	// Test transitivity: if a.Equals(b) && b.Equals(c), then a.Equals(c)
	t.Run("transitivity", func(t *testing.T) {
		a := &BufferIR{maxRequestBytes: 1024}
		b := &BufferIR{maxRequestBytes: 1024}
		c := &BufferIR{maxRequestBytes: 1024}

		assert.True(t, a.Equals(b), "a should equal b")
		assert.True(t, b.Equals(c), "b should equal c")
		assert.True(t, a.Equals(c), "a should equal c (transitivity)")
	})
}

func TestBufferIRValidate(t *testing.T) {
	// Note: Buffer validation is minimal as it only contains a single uint32 field.
	tests := []struct {
		name   string
		buffer *BufferIR
	}{
		{
			name:   "nil buffer is valid",
			buffer: nil,
		},
		{
			name: "valid buffer passes validation",
			buffer: &BufferIR{
				maxRequestBytes: 1024,
			},
		},
		{
			name: "zero value buffer is valid",
			buffer: &BufferIR{
				maxRequestBytes: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.buffer.Validate()
			assert.NoError(t, err, "Buffer validation should never fail at protobuf level")
		})
	}
}

func TestBufferIRMergeInto(t *testing.T) {
	tests := []struct {
		name            string
		p1Buffer        *BufferIR
		p2Buffer        *BufferIR
		strategy        policy.MergeOptions
		expectedP1      *BufferIR
		expectOriginSet bool
	}{
		{
			name:            "shallow merge with nil p1",
			p1Buffer:        nil,
			p2Buffer:        &BufferIR{maxRequestBytes: 2048},
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedShallowMerge},
			expectedP1:      &BufferIR{maxRequestBytes: 2048},
			expectOriginSet: true,
		},
		{
			name:            "deep merge preserves p1 when set",
			p1Buffer:        &BufferIR{maxRequestBytes: 1024},
			p2Buffer:        &BufferIR{maxRequestBytes: 2048},
			strategy:        policy.MergeOptions{Strategy: policy.AugmentedDeepMerge},
			expectedP1:      &BufferIR{maxRequestBytes: 1024}, // Should keep p1
			expectOriginSet: false,                            // No merge should happen
		},
		{
			name:            "overridable shallow merge overwrites p1",
			p1Buffer:        &BufferIR{maxRequestBytes: 1024},
			p2Buffer:        &BufferIR{maxRequestBytes: 2048},
			strategy:        policy.MergeOptions{Strategy: policy.OverridableShallowMerge},
			expectedP1:      &BufferIR{maxRequestBytes: 2048},
			expectOriginSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test policies
			p1 := &TrafficPolicy{spec: trafficPolicySpecIr{buffer: tt.p1Buffer}}
			p2 := &TrafficPolicy{spec: trafficPolicySpecIr{buffer: tt.p2Buffer}}
			p2Ref := &pluginsdkir.AttachedPolicyRef{Name: "test-policy"}
			mergeOrigins := make(pluginsdkir.MergeOrigins)

			// Perform merge
			if tt.p2Buffer != nil {
				tt.p2Buffer.MergeInto(p1, p2, p2Ref, tt.strategy, mergeOrigins)
			}

			// Verify result
			if tt.expectedP1 == nil {
				assert.Nil(t, p1.spec.buffer)
			} else {
				assert.NotNil(t, p1.spec.buffer)
				assert.Equal(t, tt.expectedP1.maxRequestBytes, p1.spec.buffer.maxRequestBytes)
			}

			// Verify origin tracking
			if tt.expectOriginSet {
				origins := mergeOrigins.Get("buffer")
				assert.NotEmpty(t, origins, "merge origin should be set")
			} else {
				origins := mergeOrigins.Get("buffer")
				assert.Empty(t, origins, "merge origin should not be set")
			}
		})
	}
}
