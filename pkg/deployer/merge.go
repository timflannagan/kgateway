package deployer

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

func DeepMergeGatewayParameters(dst, src *v1alpha1.GatewayParameters) {
	if src != nil && src.Spec.SelfManaged != nil {
		// The src override specifies a self-managed gateway, set this on the dst
		// and skip merging of kube fields that are irrelevant because of using
		// a self-managed gateway
		dst.Spec.SelfManaged = src.Spec.SelfManaged
		dst.Spec.Kube = nil
		return
	}

	// nil src override means just use dst
	if src == nil || src.Spec.Kube == nil {
		return
	}

	if dst == nil || dst.Spec.Kube == nil {
		return
	}

	dstKube := dst.Spec.Kube
	srcKube := src.Spec.Kube.DeepCopy()

	dstKube.Deployment = deepMergeDeployment(dstKube.Deployment, srcKube.Deployment)
	dstKube.EnvoyContainer = deepMergeEnvoyContainer(dstKube.EnvoyContainer, srcKube.EnvoyContainer)
	dstKube.SdsContainer = deepMergeSdsContainer(dstKube.SdsContainer, srcKube.SdsContainer)
	dstKube.PodTemplate = deepMergePodTemplate(dstKube.PodTemplate, srcKube.PodTemplate)
	dstKube.Service = deepMergeService(dstKube.Service, srcKube.Service)
	dstKube.ServiceAccount = deepMergeServiceAccount(dstKube.ServiceAccount, srcKube.ServiceAccount)
	dstKube.Istio = deepMergeIstioIntegration(dstKube.Istio, srcKube.Istio)
	dstKube.Stats = deepMergeStatsConfig(dstKube.Stats, srcKube.Stats)
	dstKube.AiExtension = deepMergeAIExtension(dstKube.AiExtension, srcKube.AiExtension)   //nolint:staticcheck
	dstKube.FloatingUserId = MergePointers(dstKube.FloatingUserId, srcKube.FloatingUserId) //nolint:staticcheck
	dstKube.OmitDefaultSecurityContext = MergePointers(dstKube.OmitDefaultSecurityContext, srcKube.OmitDefaultSecurityContext)
	dstKube.Agentgateway = deepMergeAgentgateway(dstKube.Agentgateway, srcKube.Agentgateway)
}

// MergePointers will decide whether to use dst or src without dereferencing or recursing
func MergePointers[T any](dst, src *T) *T {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	// given non-nil src override, use that instead
	return src
}

// DeepMergeMaps will use dst if src is nil, src if dest is nil, or add all entries from src into dst
// if neither are nil
func DeepMergeMaps[keyT comparable, valT any](dst, src map[keyT]valT) map[keyT]valT {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil || len(src) == 0 {
		return src
	}

	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func DeepMergeSlices[T any](dst, src []T) []T {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil || len(src) == 0 {
		return src
	}

	dst = append(dst, src...)

	return dst
}

func OverrideSlices[T any](dst, src []T) []T {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	return src
}

// Check against base value
func MergeComparable[T comparable](dst, src T) T {
	var t T
	if src == t {
		return dst
	}

	return src
}

func deepMergeStatsConfig(dst *v1alpha1.StatsConfig, src *v1alpha1.StatsConfig) *v1alpha1.StatsConfig {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Enabled = MergePointers(dst.Enabled, src.Enabled)
	dst.RoutePrefixRewrite = MergeComparable(dst.RoutePrefixRewrite, src.RoutePrefixRewrite)
	dst.EnableStatsRoute = MergeComparable(dst.EnableStatsRoute, src.EnableStatsRoute)
	dst.StatsRoutePrefixRewrite = MergeComparable(dst.StatsRoutePrefixRewrite, src.StatsRoutePrefixRewrite)

	return dst
}

func deepMergePodTemplate(dst, src *v1alpha1.Pod) *v1alpha1.Pod {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.ExtraLabels = DeepMergeMaps(dst.ExtraLabels, src.ExtraLabels)
	dst.ExtraAnnotations = DeepMergeMaps(dst.ExtraAnnotations, src.ExtraAnnotations)
	dst.SecurityContext = deepMergePodSecurityContext(dst.SecurityContext, src.SecurityContext)
	dst.ImagePullSecrets = DeepMergeSlices(dst.ImagePullSecrets, src.ImagePullSecrets)
	dst.NodeSelector = DeepMergeMaps(dst.NodeSelector, src.NodeSelector)
	dst.Affinity = DeepMergeAffinity(dst.Affinity, src.Affinity)
	dst.Tolerations = DeepMergeSlices(dst.Tolerations, src.Tolerations)
	dst.GracefulShutdown = deepMergeGracefulShutdown(dst.GracefulShutdown, src.GracefulShutdown)
	dst.TerminationGracePeriodSeconds = MergePointers(dst.TerminationGracePeriodSeconds, src.TerminationGracePeriodSeconds)
	dst.StartupProbe = deepMergeProbe(dst.StartupProbe, src.StartupProbe)
	dst.ReadinessProbe = deepMergeProbe(dst.ReadinessProbe, src.ReadinessProbe)
	dst.LivenessProbe = deepMergeProbe(dst.LivenessProbe, src.LivenessProbe)
	dst.TopologySpreadConstraints = DeepMergeSlices(dst.TopologySpreadConstraints, src.TopologySpreadConstraints)
	dst.ExtraVolumes = DeepMergeSlices(dst.ExtraVolumes, src.ExtraVolumes)

	return dst
}

func deepMergePodSecurityContext(dst, src *corev1.PodSecurityContext) *corev1.PodSecurityContext {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.SELinuxOptions = deepMergeSELinuxOptions(dst.SELinuxOptions, src.SELinuxOptions)
	dst.WindowsOptions = deepMergeWindowsSecurityContextOptions(dst.WindowsOptions, src.WindowsOptions)
	dst.RunAsUser = MergePointers(dst.RunAsUser, src.RunAsUser)
	dst.RunAsGroup = MergePointers(dst.RunAsGroup, src.RunAsGroup)
	dst.RunAsNonRoot = MergePointers(dst.RunAsNonRoot, src.RunAsNonRoot)
	dst.SupplementalGroups = DeepMergeSlices(dst.SupplementalGroups, src.SupplementalGroups)
	dst.FSGroup = MergePointers(dst.FSGroup, src.FSGroup)
	dst.Sysctls = DeepMergeSlices(dst.Sysctls, src.Sysctls)
	dst.FSGroupChangePolicy = MergePointers(dst.FSGroupChangePolicy, src.FSGroupChangePolicy)
	dst.SeccompProfile = deepMergeSeccompProfile(dst.SeccompProfile, src.SeccompProfile)

	return dst
}

func deepMergeSELinuxOptions(dst, src *corev1.SELinuxOptions) *corev1.SELinuxOptions {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.User = MergeComparable(dst.User, src.User)
	dst.Role = MergeComparable(dst.Role, src.Role)
	dst.Type = MergeComparable(dst.Type, src.Type)
	dst.Level = MergeComparable(dst.Level, src.Level)

	return dst
}

func deepMergeWindowsSecurityContextOptions(dst, src *corev1.WindowsSecurityContextOptions) *corev1.WindowsSecurityContextOptions {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.GMSACredentialSpecName = MergePointers(dst.GMSACredentialSpec, src.GMSACredentialSpec)
	dst.GMSACredentialSpec = MergePointers(dst.GMSACredentialSpec, src.GMSACredentialSpec)
	dst.RunAsUserName = MergePointers(dst.RunAsUserName, src.RunAsUserName)
	dst.HostProcess = MergePointers(dst.HostProcess, src.HostProcess)

	return dst
}

func deepMergeSeccompProfile(dst, src *corev1.SeccompProfile) *corev1.SeccompProfile {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Type = MergeComparable(dst.Type, src.Type)
	dst.LocalhostProfile = MergePointers(dst.LocalhostProfile, src.LocalhostProfile)

	return dst
}

func DeepMergeAffinity(dst, src *corev1.Affinity) *corev1.Affinity {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.NodeAffinity = deepMergeNodeAffinity(dst.NodeAffinity, src.NodeAffinity)
	dst.PodAffinity = deepMergePodAffinity(dst.PodAffinity, src.PodAffinity)
	dst.PodAntiAffinity = deepMergePodAntiAffinity(dst.PodAntiAffinity, src.PodAntiAffinity)

	return dst
}

func deepMergeNodeAffinity(dst, src *corev1.NodeAffinity) *corev1.NodeAffinity {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.RequiredDuringSchedulingIgnoredDuringExecution = deepMergeNodeSelector(dst.RequiredDuringSchedulingIgnoredDuringExecution, src.RequiredDuringSchedulingIgnoredDuringExecution)
	dst.PreferredDuringSchedulingIgnoredDuringExecution = DeepMergeSlices(dst.PreferredDuringSchedulingIgnoredDuringExecution, src.PreferredDuringSchedulingIgnoredDuringExecution)

	return dst
}

func deepMergeNodeSelector(dst, src *corev1.NodeSelector) *corev1.NodeSelector {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.NodeSelectorTerms = DeepMergeSlices(dst.NodeSelectorTerms, src.NodeSelectorTerms)

	return dst
}

func deepMergePodAffinity(dst, src *corev1.PodAffinity) *corev1.PodAffinity {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.RequiredDuringSchedulingIgnoredDuringExecution = DeepMergeSlices(dst.RequiredDuringSchedulingIgnoredDuringExecution, src.RequiredDuringSchedulingIgnoredDuringExecution)
	dst.PreferredDuringSchedulingIgnoredDuringExecution = DeepMergeSlices(dst.PreferredDuringSchedulingIgnoredDuringExecution, src.PreferredDuringSchedulingIgnoredDuringExecution)

	return dst
}

func deepMergePodAntiAffinity(dst, src *corev1.PodAntiAffinity) *corev1.PodAntiAffinity {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.RequiredDuringSchedulingIgnoredDuringExecution = DeepMergeSlices(dst.RequiredDuringSchedulingIgnoredDuringExecution, src.RequiredDuringSchedulingIgnoredDuringExecution)
	dst.PreferredDuringSchedulingIgnoredDuringExecution = DeepMergeSlices(dst.PreferredDuringSchedulingIgnoredDuringExecution, src.PreferredDuringSchedulingIgnoredDuringExecution)

	return dst
}

func deepMergeProbe(dst, src *corev1.Probe) *corev1.Probe {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.ProbeHandler = deepMergeProbeHandler(dst.ProbeHandler, src.ProbeHandler)
	dst.InitialDelaySeconds = MergeComparable(dst.InitialDelaySeconds, src.InitialDelaySeconds)
	dst.TimeoutSeconds = MergeComparable(dst.TimeoutSeconds, src.TimeoutSeconds)
	dst.PeriodSeconds = MergeComparable(dst.PeriodSeconds, src.PeriodSeconds)
	dst.SuccessThreshold = MergeComparable(dst.SuccessThreshold, src.SuccessThreshold)
	dst.FailureThreshold = MergeComparable(dst.FailureThreshold, src.FailureThreshold)
	dst.TerminationGracePeriodSeconds = MergePointers(dst.TerminationGracePeriodSeconds, src.TerminationGracePeriodSeconds)

	return dst
}

func deepMergeProbeHandler(dst, src corev1.ProbeHandler) corev1.ProbeHandler {
	srcHasExecAction := src.Exec != nil
	srcHasHTTPGetAction := src.HTTPGet != nil
	srcHasTCPSocketAction := src.TCPSocket != nil
	srcHasGRPCAction := src.GRPC != nil
	srcHasAction := srcHasExecAction || srcHasHTTPGetAction || srcHasTCPSocketAction || srcHasGRPCAction
	if srcHasAction {
		// Reset the dest so it does not conflict with the src Action as there should only be one Action defined per probe
		dst.Exec = nil
		dst.HTTPGet = nil
		dst.TCPSocket = nil
		dst.GRPC = nil
	}
	// kube-builder validation ensures that the src has only one action
	dst.Exec = deepMergeExecAction(dst.Exec, src.Exec)
	dst.HTTPGet = deepMergeHTTPGetAction(dst.HTTPGet, src.HTTPGet)
	dst.TCPSocket = deepMergeTCPSocketAction(dst.TCPSocket, src.TCPSocket)
	dst.GRPC = deepMergeGRPCAction(dst.GRPC, src.GRPC)

	return dst
}

func deepMergeExecAction(dst, src *corev1.ExecAction) *corev1.ExecAction {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	// Don't merge the command string as that can break the entire probe
	dst.Command = OverrideSlices(dst.Command, src.Command)

	return dst
}

func deepMergeHTTPGetAction(dst, src *corev1.HTTPGetAction) *corev1.HTTPGetAction {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Path = MergeComparable(dst.Path, src.Path)
	dst.Port = mergeIntOrString(dst.Port, src.Port)
	dst.Host = MergeComparable(dst.Host, src.Host)
	dst.Scheme = MergeComparable(dst.Scheme, src.Scheme)
	dst.HTTPHeaders = DeepMergeSlices(dst.HTTPHeaders, src.HTTPHeaders)

	return dst
}

func mergeIntOrString(dst, src intstr.IntOrString) intstr.IntOrString {
	// Do not deep merge as this can cause a conflict between the name and number of the port to access on the container
	return MergeComparable(dst, src)
}

func deepMergeTCPSocketAction(dst, src *corev1.TCPSocketAction) *corev1.TCPSocketAction {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Port = mergeIntOrString(dst.Port, src.Port)
	dst.Host = MergeComparable(dst.Host, src.Host)

	return dst
}

func deepMergeGRPCAction(dst, src *corev1.GRPCAction) *corev1.GRPCAction {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Port = MergeComparable(dst.Port, src.Port)
	dst.Service = MergePointers(dst.Service, src.Service)

	return dst
}

func deepMergeGracefulShutdown(dst, src *v1alpha1.GracefulShutdownSpec) *v1alpha1.GracefulShutdownSpec {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Enabled = MergePointers(dst.Enabled, src.Enabled)
	dst.SleepTimeSeconds = MergePointers(dst.SleepTimeSeconds, src.SleepTimeSeconds)

	return dst
}

func deepMergeService(dst, src *v1alpha1.Service) *v1alpha1.Service {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	if src.Type != nil {
		dst.Type = src.Type
	}

	if src.ClusterIP != nil {
		dst.ClusterIP = src.ClusterIP
	}

	dst.ExtraLabels = DeepMergeMaps(dst.ExtraLabels, src.ExtraLabels)
	dst.ExtraAnnotations = DeepMergeMaps(dst.ExtraAnnotations, src.ExtraAnnotations)
	dst.Ports = DeepMergeSlices(dst.Ports, src.Ports)
	if src.ExternalTrafficPolicy != nil {
		dst.ExternalTrafficPolicy = src.ExternalTrafficPolicy
	}

	return dst
}

func deepMergeServiceAccount(dst, src *v1alpha1.ServiceAccount) *v1alpha1.ServiceAccount {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.ExtraLabels = DeepMergeMaps(dst.ExtraLabels, src.ExtraLabels)
	dst.ExtraAnnotations = DeepMergeMaps(dst.ExtraAnnotations, src.ExtraAnnotations)

	return dst
}

func deepMergeSdsContainer(dst, src *v1alpha1.SdsContainer) *v1alpha1.SdsContainer {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Image = DeepMergeImage(dst.Image, src.Image)
	dst.SecurityContext = DeepMergeSecurityContext(dst.SecurityContext, src.SecurityContext)
	dst.Resources = DeepMergeResourceRequirements(dst.Resources, src.Resources)
	dst.Bootstrap = deepMergeSdsBootstrap(dst.Bootstrap, src.Bootstrap)

	return dst
}

func deepMergeSdsBootstrap(dst, src *v1alpha1.SdsBootstrap) *v1alpha1.SdsBootstrap {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	if src.LogLevel != nil {
		dst.LogLevel = src.LogLevel
	}

	return dst
}

func deepMergeIstioIntegration(dst, src *v1alpha1.IstioIntegration) *v1alpha1.IstioIntegration {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.IstioProxyContainer = deepMergeIstioContainer(dst.IstioProxyContainer, src.IstioProxyContainer)
	dst.CustomSidecars = mergeCustomSidecars(dst.CustomSidecars, src.CustomSidecars)

	return dst
}

// mergeCustomSidecars will decide whether to use dst or src custom sidecar containers
func mergeCustomSidecars(dst, src []corev1.Container) []corev1.Container {
	// nil src override means just use dst
	if len(src) == 0 {
		return dst
	}

	// given non-nil src override, use that instead
	return src
}

func deepMergeIstioContainer(dst, src *v1alpha1.IstioContainer) *v1alpha1.IstioContainer {
	// nil src override means just use dst
	if src == nil {
		return dst
	}
	if dst == nil {
		return src
	}

	dst.Image = DeepMergeImage(dst.Image, src.Image)
	dst.SecurityContext = DeepMergeSecurityContext(dst.SecurityContext, src.SecurityContext)
	dst.Resources = DeepMergeResourceRequirements(dst.Resources, src.Resources)

	if src.LogLevel != nil {
		dst.LogLevel = src.LogLevel
	}

	if src.IstioDiscoveryAddress != nil {
		dst.IstioDiscoveryAddress = src.IstioDiscoveryAddress
	}

	if src.IstioMetaMeshId != nil {
		dst.IstioMetaMeshId = src.IstioMetaMeshId
	}

	if src.IstioMetaClusterId != nil {
		dst.IstioMetaClusterId = src.IstioMetaClusterId
	}

	return dst
}

func deepMergeEnvoyContainer(dst, src *v1alpha1.EnvoyContainer) *v1alpha1.EnvoyContainer {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Bootstrap = deepMergeEnvoyBootstrap(dst.Bootstrap, src.Bootstrap)
	dst.Image = DeepMergeImage(dst.Image, src.Image)
	dst.SecurityContext = DeepMergeSecurityContext(dst.SecurityContext, src.SecurityContext)
	dst.Resources = DeepMergeResourceRequirements(dst.Resources, src.Resources)
	dst.Env = DeepMergeSlices(dst.Env, src.Env)
	dst.ExtraVolumeMounts = DeepMergeSlices(dst.ExtraVolumeMounts, src.ExtraVolumeMounts)

	return dst
}

func DeepMergeImage(dst, src *v1alpha1.Image) *v1alpha1.Image {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	if src.Registry != nil {
		dst.Registry = src.Registry
	}

	if src.Repository != nil {
		dst.Repository = src.Repository
	}

	if src.Tag != nil {
		dst.Tag = src.Tag
	}

	if src.Digest != nil {
		dst.Digest = src.Digest
	}

	if src.PullPolicy != nil {
		dst.PullPolicy = src.PullPolicy
	}

	return dst
}

func deepMergeEnvoyBootstrap(dst, src *v1alpha1.EnvoyBootstrap) *v1alpha1.EnvoyBootstrap {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}
	if src.LogLevel != nil {
		dst.LogLevel = src.LogLevel
	}

	dst.ComponentLogLevels = DeepMergeMaps(dst.ComponentLogLevels, src.ComponentLogLevels)

	return dst
}

func DeepMergeResourceRequirements(dst, src *corev1.ResourceRequirements) *corev1.ResourceRequirements {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Limits = DeepMergeMaps(dst.Limits, src.Limits)
	dst.Requests = DeepMergeMaps(dst.Requests, src.Requests)

	return dst
}

func DeepMergeSecurityContext(dst, src *corev1.SecurityContext) *corev1.SecurityContext {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Capabilities = deepMergeCapabilities(dst.Capabilities, src.Capabilities)
	dst.SELinuxOptions = deepMergeSELinuxOptions(dst.SELinuxOptions, src.SELinuxOptions)
	dst.WindowsOptions = deepMergeWindowsSecurityContextOptions(dst.WindowsOptions, src.WindowsOptions)
	dst.RunAsUser = MergePointers(dst.RunAsUser, src.RunAsUser)
	dst.RunAsGroup = MergePointers(dst.RunAsGroup, src.RunAsGroup)
	dst.RunAsNonRoot = MergePointers(dst.RunAsNonRoot, src.RunAsNonRoot)
	dst.Privileged = MergePointers(dst.Privileged, src.Privileged)
	dst.ReadOnlyRootFilesystem = MergePointers(dst.ReadOnlyRootFilesystem, src.ReadOnlyRootFilesystem)
	dst.AllowPrivilegeEscalation = MergePointers(dst.AllowPrivilegeEscalation, src.AllowPrivilegeEscalation)
	dst.ProcMount = MergePointers(dst.ProcMount, src.ProcMount)
	dst.SeccompProfile = deepMergeSeccompProfile(dst.SeccompProfile, src.SeccompProfile)

	return dst
}

func deepMergeCapabilities(dst, src *corev1.Capabilities) *corev1.Capabilities {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Add = DeepMergeSlices(dst.Add, src.Add)
	dst.Drop = DeepMergeSlices(dst.Drop, src.Drop)

	return dst
}

func deepMergeDeployment(dst, src *v1alpha1.ProxyDeployment) *v1alpha1.ProxyDeployment {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Replicas = MergePointers(dst.Replicas, src.Replicas)
	dst.Strategy = MergePointers(dst.Strategy, src.Strategy)

	return dst
}

func deepMergeAIExtension(dst, src *v1alpha1.AiExtension) *v1alpha1.AiExtension {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Enabled = MergePointers(dst.Enabled, src.Enabled)
	dst.Image = DeepMergeImage(dst.Image, src.Image)
	dst.SecurityContext = DeepMergeSecurityContext(dst.SecurityContext, src.SecurityContext)
	dst.Resources = DeepMergeResourceRequirements(dst.Resources, src.Resources)
	dst.Env = DeepMergeSlices(dst.Env, src.Env)
	dst.Ports = DeepMergeSlices(dst.Ports, src.Ports)
	dst.Stats = deepMergeAIExtensionStats(dst.Stats, src.Stats)
	dst.Tracing = deepMergeAIExtensionTracing(dst.Tracing, src.Tracing)
	return dst
}

func deepMergeAIExtensionTracing(dst, src *v1alpha1.AiExtensionTrace) *v1alpha1.AiExtensionTrace {
	// nil src override means just use dst
	if src == nil {
		return dst
	}
	if dst == nil {
		return src
	}
	dst.EndPoint = MergeComparable(dst.EndPoint, src.EndPoint)
	dst.Sampler = MergePointers(dst.Sampler, src.Sampler)
	dst.Timeout = MergePointers(dst.Timeout, src.Timeout)
	dst.Protocol = MergePointers(dst.Protocol, src.Protocol)
	return dst
}

func deepMergeAIExtensionStats(dst, src *v1alpha1.AiExtensionStats) *v1alpha1.AiExtensionStats {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.CustomLabels = DeepMergeSlices(dst.CustomLabels, src.CustomLabels)

	return dst
}

func deepMergeAgentgateway(dst, src *v1alpha1.Agentgateway) *v1alpha1.Agentgateway {
	// nil src override means just use dst
	if src == nil {
		return dst
	}

	if dst == nil {
		return src
	}

	dst.Enabled = MergePointers(dst.Enabled, src.Enabled)
	dst.LogLevel = MergePointers(dst.LogLevel, src.LogLevel)
	dst.Image = DeepMergeImage(dst.Image, src.Image)
	dst.SecurityContext = DeepMergeSecurityContext(dst.SecurityContext, src.SecurityContext)
	dst.Resources = DeepMergeResourceRequirements(dst.Resources, src.Resources)
	dst.Env = DeepMergeSlices(dst.Env, src.Env)
	dst.CustomConfigMapName = MergePointers(dst.CustomConfigMapName, src.CustomConfigMapName)
	dst.ExtraVolumeMounts = DeepMergeSlices(dst.ExtraVolumeMounts, src.ExtraVolumeMounts)

	return dst
}
