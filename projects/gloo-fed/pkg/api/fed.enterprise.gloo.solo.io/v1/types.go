// Code generated by skv2. DO NOT EDIT.

// Definitions for the Kubernetes types
package v1

import (
	i8e2432a4f6c240b117e19c109dc2dcd2 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status

// GroupVersionKind for FederatedAuthConfig
var FederatedAuthConfigGVK = schema.GroupVersionKind{
	Group:   "fed.enterprise.gloo.solo.io",
	Version: "v1",
	Kind:    "FederatedAuthConfig",
}

// FederatedAuthConfig is the Schema for the federatedAuthConfig API
type FederatedAuthConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   i8e2432a4f6c240b117e19c109dc2dcd2.FederatedAuthConfigSpec   `json:"spec,omitempty"`
	Status i8e2432a4f6c240b117e19c109dc2dcd2.FederatedAuthConfigStatus `json:"status,omitempty"`
}

// GVK returns the GroupVersionKind associated with the resource type.
func (FederatedAuthConfig) GVK() schema.GroupVersionKind {
	return FederatedAuthConfigGVK
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FederatedAuthConfigList contains a list of FederatedAuthConfig
type FederatedAuthConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FederatedAuthConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FederatedAuthConfig{}, &FederatedAuthConfigList{})
}
