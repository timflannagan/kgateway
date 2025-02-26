// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	apiv1alpha1 "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

// RegexApplyConfiguration represents a declarative configuration of the Regex type for use
// with apply.
type RegexApplyConfiguration struct {
	Matches  []RegexMatchApplyConfiguration `json:"matches,omitempty"`
	Builtins []apiv1alpha1.BuiltIn          `json:"builtins,omitempty"`
	Action   *apiv1alpha1.Action            `json:"action,omitempty"`
}

// RegexApplyConfiguration constructs a declarative configuration of the Regex type for use with
// apply.
func Regex() *RegexApplyConfiguration {
	return &RegexApplyConfiguration{}
}

// WithMatches adds the given value to the Matches field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Matches field.
func (b *RegexApplyConfiguration) WithMatches(values ...*RegexMatchApplyConfiguration) *RegexApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithMatches")
		}
		b.Matches = append(b.Matches, *values[i])
	}
	return b
}

// WithBuiltins adds the given value to the Builtins field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Builtins field.
func (b *RegexApplyConfiguration) WithBuiltins(values ...apiv1alpha1.BuiltIn) *RegexApplyConfiguration {
	for i := range values {
		b.Builtins = append(b.Builtins, values[i])
	}
	return b
}

// WithAction sets the Action field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Action field is set to the value of the last call.
func (b *RegexApplyConfiguration) WithAction(value apiv1alpha1.Action) *RegexApplyConfiguration {
	b.Action = &value
	return b
}
