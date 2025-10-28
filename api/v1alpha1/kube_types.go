package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

// A container image. See https://kubernetes.io/docs/concepts/containers/images
// for details.
type Image struct {
	// The image registry.
	//
	// +optional
	Registry *string `json:"registry,omitempty"`

	// The image repository (name).
	//
	// +optional
	Repository *string `json:"repository,omitempty"`

	// The image tag.
	//
	// +optional
	Tag *string `json:"tag,omitempty"`

	// The hash digest of the image, e.g. `sha256:12345...`
	//
	// +optional
	Digest *string `json:"digest,omitempty"`

	// The image pull policy for the container. See
	// https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy
	// for details.
	//
	// +optional
	PullPolicy *corev1.PullPolicy `json:"pullPolicy,omitempty"`
}


// Configuration for a Kubernetes Service.
type Service struct {
	// The Kubernetes Service type.
	//
	// +optional
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer;ExternalName
	Type *corev1.ServiceType `json:"type,omitempty"`

	// The manually specified IP address of the service, if a randomly assigned
	// IP is not desired. See
	// https://kubernetes.io/docs/concepts/services-networking/service/#choosing-your-own-ip-address
	// and
	// https://kubernetes.io/docs/concepts/services-networking/service/#headless-services
	// on the implications of setting `clusterIP`.
	//
	// +optional
	ClusterIP *string `json:"clusterIP,omitempty"`

	// Additional labels to add to the Service object metadata.
	// If the same label is present on `Gateway.spec.infrastructure.labels`, the `Gateway` takes precedence.
	//
	// +optional
	ExtraLabels map[string]string `json:"extraLabels,omitempty"`

	// Additional annotations to add to the Service object metadata.
	// If the same annotation is present on `Gateway.spec.infrastructure.annotations`, the `Gateway` takes precedence.
	//
	// +optional
	ExtraAnnotations map[string]string `json:"extraAnnotations,omitempty"`

	// Additional configuration for the service ports.
	// The actual port numbers are specified in the Gateway resource.
	//
	// +optional
	Ports []Port `json:"ports,omitempty"`

	// ExternalTrafficPolicy defines the external traffic policy for the service.
	// Valid values are Cluster and Local. Default value is Cluster.
	//
	// +optional
	ExternalTrafficPolicy *string `json:"externalTrafficPolicy,omitempty"`
}


type Port struct {
	// The port number to match on the Gateway
	//
	// +required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// The NodePort to be used for the service. If not specified, a random port
	// will be assigned by the Kubernetes API server.
	//
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	NodePort *int32 `json:"nodePort,omitempty"`
}



type ServiceAccount struct {
	// Additional labels to add to the ServiceAccount object metadata.
	//
	// +optional
	ExtraLabels map[string]string `json:"extraLabels,omitempty"`

	// Additional annotations to add to the ServiceAccount object metadata.
	// If the same annotation is present on `Gateway.spec.infrastructure.annotations`, the `Gateway` takes precedence.
	//
	// +optional
	ExtraAnnotations map[string]string `json:"extraAnnotations,omitempty"`
}


// Configuration for a Kubernetes Pod template.
type Pod struct {
	// Additional labels to add to the Pod object metadata.
	// If the same label is present on `Gateway.spec.infrastructure.labels`, the `Gateway` takes precedence.
	//
	// +optional
	ExtraLabels map[string]string `json:"extraLabels,omitempty"`

	// Additional annotations to add to the Pod object metadata.
	// If the same annotation is present on `Gateway.spec.infrastructure.annotations`, the `Gateway` takes precedence.
	//
	// +optional
	ExtraAnnotations map[string]string `json:"extraAnnotations,omitempty"`

	// The pod security context. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core
	// for details.
	//
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`

	// An optional list of references to secrets in the same namespace to use for
	// pulling any of the images used by this Pod spec. See
	// https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod
	// for details.
	//
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// A selector which must be true for the pod to fit on a node. See
	// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/ for
	// details.
	//
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// If specified, the pod's scheduling constraints. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core
	// for details.
	//
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// do not use slice of pointers: https://github.com/kubernetes/code-generator/issues/166
	// If specified, the pod's tolerations. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core
	// for details.
	//
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// If specified, the pod's graceful shutdown spec.
	//
	// +optional
	GracefulShutdown *GracefulShutdownSpec `json:"gracefulShutdown,omitempty"`

	// If specified, the pod's termination grace period in seconds. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pod-v1-core
	// for details
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=31536000
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`

	// If specified, the pod's startup probe. A probe of container startup readiness.
	// Container will be only be added to service endpoints if the probe succeeds. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#probe-v1-core
	// for details.
	//
	// +optional
	StartupProbe *corev1.Probe `json:"startupProbe,omitempty"`

	// If specified, the pod's readiness probe. Periodic probe of container service readiness.
	// Container will be removed from service endpoints if the probe fails. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#probe-v1-core
	// for details.
	//
	// +optional
	ReadinessProbe *corev1.Probe `json:"readinessProbe,omitempty"`

	// If specified, the pod's liveness probe. Periodic probe of container service readiness.
	// Container will be restarted if the probe fails. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#probe-v1-core
	// for details.
	//
	// +optional
	LivenessProbe *corev1.Probe `json:"livenessProbe,omitempty"`

	// If specified, the pod's topology spread constraints. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#topologyspreadconstraint-v1-core
	// for details.
	//
	// +optional
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`

	// Additional volumes to add to the pod. See
	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#volume-v1-core
	// for details.
	//
	// +optional
	ExtraVolumes []corev1.Volume `json:"extraVolumes,omitempty"`
}


type GracefulShutdownSpec struct {
	// Enable grace period before shutdown to finish current requests while Envoy health checks fail to e.g. notify external load balancers. *NOTE:* This will not have any effect if you have not defined health checks via the health check filter
	//
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Time (in seconds) for the preStop hook to wait before allowing Envoy to terminate
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=31536000
	SleepTimeSeconds *int64 `json:"sleepTimeSeconds,omitempty"`
}

