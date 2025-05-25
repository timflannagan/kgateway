package policy

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorType indicates the severity of the error
type ErrorType int

const (
	// ErrorTypeWarning indicates a non-fatal error that should be reported but doesn't prevent policy application
	ErrorTypeWarning ErrorType = iota
	// ErrorTypeTerminal indicates a fatal error that prevents policy application
	ErrorTypeTerminal
)

// PolicyError represents an error that occurred during policy processing
type PolicyError struct {
	// Type indicates whether this is a warning or terminal error
	Type ErrorType
	// Message contains the error message
	Message string
	// Err contains the original error if any
	Err error
}

// Error implements the error interface
func (e *PolicyError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *PolicyError) Unwrap() error {
	return e.Err
}

// IsWarning returns true if this is a warning error
func (e *PolicyError) IsWarning() bool {
	return e.Type == ErrorTypeWarning
}

// IsTerminal returns true if this is a terminal error
func (e *PolicyError) IsTerminal() bool {
	return e.Type == ErrorTypeTerminal
}

// NewWarningError creates a new warning error
func NewWarningError(msg string, err error) *PolicyError {
	return &PolicyError{
		Type:    ErrorTypeWarning,
		Message: msg,
		Err:     err,
	}
}

// NewTerminalError creates a new terminal error
func NewTerminalError(msg string, err error) *PolicyError {
	return &PolicyError{
		Type:    ErrorTypeTerminal,
		Message: msg,
		Err:     err,
	}
}

// NewEnvoyValidationError creates a new terminal error specifically for Envoy validation failures
func NewEnvoyValidationError(msg string, err error) *PolicyError {
	return &PolicyError{
		Type:    ErrorTypeTerminal,
		Message: fmt.Sprintf("EnvoyValidationFailed: %s", msg),
		Err:     err,
	}
}

// IsEnvoyValidationError checks if the error is an Envoy validation error
func IsEnvoyValidationError(err error) bool {
	var policyErr *PolicyError
	if !errors.As(err, &policyErr) {
		return false
	}
	return strings.HasPrefix(policyErr.Message, "EnvoyValidationFailed:")
}
