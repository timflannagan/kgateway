package validator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
)

// Validator defines the interface for validating Envoy xDS configurations
type Validator interface {
	// Validate validates the given xDS configuration from the provided reader
	// Returns an error if validation fails or if there's an error during validation
	Validate(ctx context.Context, config io.Reader) error
}

// envoyValidator implements the Validator interface using the Envoy binary
type envoyValidator struct {
	envoyPath string
}

var _ Validator = &envoyValidator{}

// New creates a new Validator instance
//
// TODO: avoid hardcoding the envoy path.
func New() Validator {
	return &envoyValidator{
		envoyPath: "/usr/local/bin/envoy",
	}
}

// Validate validates the given xDS configuration from the provided reader to
// determine whether Envoy will ACK or NACK the config. This is provides an additional
// guardrail to prevent invalid xDS configuration from being pushed to managed Envoy
// instances.
func (v *envoyValidator) Validate(ctx context.Context, config io.Reader) error {
	// TODO: measure time taken to validate for larger configs.
	configBytes, err := io.ReadAll(config)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	cmd := exec.CommandContext(ctx, v.envoyPath, "--mode", "validate", "--config-yaml", "/dev/fd/0")
	cmd.Stdin = bytes.NewReader(configBytes)
	return cmd.Run()
}
