package policy

import "fmt"

// ErrorType indicates the severity of the error
type ErrorType int

const (
	// ErrorTypeTerminal indicates a fatal error that prevents policy application
	ErrorTypeTerminal ErrorType = iota
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

// IsTerminal returns true if this is a terminal error
func (e *PolicyError) IsTerminal() bool {
	return e.Type == ErrorTypeTerminal
}

// NewTerminalError creates a new terminal error
func NewTerminalError(msg string, err error) *PolicyError {
	return &PolicyError{
		Type:    ErrorTypeTerminal,
		Message: msg,
		Err:     err,
	}
}
