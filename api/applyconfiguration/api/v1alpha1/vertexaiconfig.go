// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	apiv1alpha1 "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

// VertexAIConfigApplyConfiguration represents a declarative configuration of the VertexAIConfig type for use
// with apply.
type VertexAIConfigApplyConfiguration struct {
	AuthToken  *SingleAuthTokenApplyConfiguration `json:"authToken,omitempty"`
	Model      *string                            `json:"model,omitempty"`
	ApiVersion *string                            `json:"apiVersion,omitempty"`
	ProjectId  *string                            `json:"projectId,omitempty"`
	Location   *string                            `json:"location,omitempty"`
	ModelPath  *string                            `json:"modelPath,omitempty"`
	Publisher  *apiv1alpha1.Publisher             `json:"publisher,omitempty"`
}

// VertexAIConfigApplyConfiguration constructs a declarative configuration of the VertexAIConfig type for use with
// apply.
func VertexAIConfig() *VertexAIConfigApplyConfiguration {
	return &VertexAIConfigApplyConfiguration{}
}

// WithAuthToken sets the AuthToken field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AuthToken field is set to the value of the last call.
func (b *VertexAIConfigApplyConfiguration) WithAuthToken(value *SingleAuthTokenApplyConfiguration) *VertexAIConfigApplyConfiguration {
	b.AuthToken = value
	return b
}

// WithModel sets the Model field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Model field is set to the value of the last call.
func (b *VertexAIConfigApplyConfiguration) WithModel(value string) *VertexAIConfigApplyConfiguration {
	b.Model = &value
	return b
}

// WithApiVersion sets the ApiVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ApiVersion field is set to the value of the last call.
func (b *VertexAIConfigApplyConfiguration) WithApiVersion(value string) *VertexAIConfigApplyConfiguration {
	b.ApiVersion = &value
	return b
}

// WithProjectId sets the ProjectId field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ProjectId field is set to the value of the last call.
func (b *VertexAIConfigApplyConfiguration) WithProjectId(value string) *VertexAIConfigApplyConfiguration {
	b.ProjectId = &value
	return b
}

// WithLocation sets the Location field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Location field is set to the value of the last call.
func (b *VertexAIConfigApplyConfiguration) WithLocation(value string) *VertexAIConfigApplyConfiguration {
	b.Location = &value
	return b
}

// WithModelPath sets the ModelPath field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ModelPath field is set to the value of the last call.
func (b *VertexAIConfigApplyConfiguration) WithModelPath(value string) *VertexAIConfigApplyConfiguration {
	b.ModelPath = &value
	return b
}

// WithPublisher sets the Publisher field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Publisher field is set to the value of the last call.
func (b *VertexAIConfigApplyConfiguration) WithPublisher(value apiv1alpha1.Publisher) *VertexAIConfigApplyConfiguration {
	b.Publisher = &value
	return b
}
