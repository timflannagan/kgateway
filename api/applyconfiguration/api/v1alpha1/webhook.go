// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

// WebhookApplyConfiguration represents a declarative configuration of the Webhook type for use
// with apply.
type WebhookApplyConfiguration struct {
	Host           *HostApplyConfiguration `json:"host,omitempty"`
	ForwardHeaders []v1.HTTPHeaderMatch    `json:"forwardHeaders,omitempty"`
}

// WebhookApplyConfiguration constructs a declarative configuration of the Webhook type for use with
// apply.
func Webhook() *WebhookApplyConfiguration {
	return &WebhookApplyConfiguration{}
}

// WithHost sets the Host field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Host field is set to the value of the last call.
func (b *WebhookApplyConfiguration) WithHost(value *HostApplyConfiguration) *WebhookApplyConfiguration {
	b.Host = value
	return b
}

// WithForwardHeaders adds the given value to the ForwardHeaders field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the ForwardHeaders field.
func (b *WebhookApplyConfiguration) WithForwardHeaders(values ...v1.HTTPHeaderMatch) *WebhookApplyConfiguration {
	for i := range values {
		b.ForwardHeaders = append(b.ForwardHeaders, values[i])
	}
	return b
}
