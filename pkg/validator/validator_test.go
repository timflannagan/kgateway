package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBinaryValidator(t *testing.T) {
	validator := NewBinaryValidator("/usr/local/bin/envoy")
	assert.NotNil(t, validator)
}

func TestNewCachedValidator(t *testing.T) {
	baseValidator := NewBinaryValidator("/usr/local/bin/envoy")
	cachedValidator := NewCachedValidator(baseValidator)
	assert.NotNil(t, cachedValidator)
}

func TestValidationCache(t *testing.T) {
	cache := &validationCache{
		cache: make(map[string]error),
	}

	// Mock validator that counts calls
	callCount := 0
	mockValidator := func(ctx context.Context, yaml string) error {
		callCount++
		return nil
	}

	yaml := "test: config"

	// First call should invoke the validator
	err := cache.validate(context.Background(), yaml, mockValidator)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Second call with same YAML should use cache
	err = cache.validate(context.Background(), yaml, mockValidator)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount) // Should not increment

	// Different YAML should invoke validator again
	err = cache.validate(context.Background(), "different: config", mockValidator)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)
}
