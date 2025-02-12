package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=upstreams,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.kgateway.dev,resources=upstreams/status,verbs=get;update;patch

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:metadata:labels={app=gateway,app.kubernetes.io/name=gateway}
// +kubebuilder:resource:categories=kgateway,shortName=up
// +kubebuilder:subresource:status
type Upstream struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpstreamSpec   `json:"spec,omitempty"`
	Status UpstreamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type UpstreamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Upstream `json:"items"`
}

const (
	UpstreamTypeAws    = "aws"
	UpstreamTypeStatic = "static"
)

// +kubebuilder:validation:XValidation:message="Type field must match the upstream configuration: aws requires Aws field, static requires Static field",rule="(self.type == 'aws' && has(self.aws) && !has(self.static)) || (self.type == 'static' && has(self.static) && !has(self.aws))"
type UpstreamSpec struct {
	// Type controls the type of upstream.
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=aws;static
	Type string `json:"type"`
	// AWS upstream configuration. Allows for referencing AWS Lambda functions.
	// +optional
	Aws *AwsUpstream `json:"aws,omitempty"`
	// Static upstream configuration. Allows for referencing a list of hosts.
	// +optional
	Static *StaticUpstream `json:"static,omitempty"`
}

type AwsUpstream struct {
	Region    string                      `json:"region,omitempty"`
	AccountId string                      `json:"accountId,omitempty"`
	SecretRef corev1.LocalObjectReference `json:"secretRef,omitempty"`
	Lambda    *AwsLambdaUpstream          `json:"lambda,omitempty"`
}

type AwsLambdaUpstream struct {
	// The name of the Lambda function
	// +kubebuilder:validation:Required
	FunctionName string `json:"functionName"`
	// Optional qualifier (version or alias)
	// +optional
	Qualifier string `json:"qualifier,omitempty"`
}

type StaticUpstream struct {
	Hosts []Host `json:"hosts,omitempty"`
}

type Host struct {
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Host string          `json:"host"`
	Port gwv1.PortNumber `json:"port"`
}

type UpstreamStatus struct {
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
