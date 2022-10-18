package glooinstance_handler

import (
	"context"
	"os"

	"github.com/rotisserie/eris"
	apps_v1 "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1"
	core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	sk_sets "github.com/solo-io/skv2/contrib/pkg/sets/v2"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/pkg/ezkube"
	enterprise_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	gateway_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	types "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	ratelimit_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	enterprise_gloo_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/enterprise.gloo.solo.io/v1/handler"
	gateway_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/gateway.solo.io/v1/handler"
	gloo_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/gloo.solo.io/v1/handler"
	ratelimit_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/ratelimit.api.solo.io/v1alpha1/handler"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/checker"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/images"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/locality"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/proxy"
	"go.uber.org/zap"
	k8s_apps_types "k8s.io/api/apps/v1"
	k8s_core_types "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Single-cluster Gloo Edge doesn't have a notion of clusters like Gloo Fed (which requires cluster registration
	// with user-provided cluster names), so we hardcode an arbitrary cluster name here for display purposes
	ClusterName = "cluster"

	// Note: currently we only show the GlooInstance in the current install namespace. In the future, if we need to
	// support listing all GlooInstances on a cluster, make this a configurable helm flag
	SingleInstanceOnly = true
)

//go:generate mockgen -source ./single_cluster_lister.go -destination mocks/single_cluster_lister.go

func NewSingleClusterGlooInstanceLister(
	coreClientset core_v1.Clientset,
	appsClientset apps_v1.Clientset,
	gatewayClientset gateway_v1.Clientset,
	glooClientset gloo_v1.Clientset,
	enterpriseGlooClientset enterprise_gloo_v1.Clientset,
	ratelimitClientset ratelimit_v1alpha1.Clientset,
) SingleClusterGlooInstanceLister {
	return &singleClusterGlooInstanceLister{
		coreClientset:           coreClientset,
		appsClientset:           appsClientset,
		gatewayClientset:        gatewayClientset,
		glooClientset:           glooClientset,
		enterpriseGlooClientset: enterpriseGlooClientset,
		ratelimitClientset:      ratelimitClientset,
	}
}

type SingleClusterGlooInstanceLister interface {
	ListGlooInstances(ctx context.Context) ([]*rpc_edge_v1.GlooInstance, error)
	GetGlooInstance(ctx context.Context, glooInstanceRef *skv2_v1.ObjectRef) (*rpc_edge_v1.GlooInstance, error)
}

type singleClusterGlooInstanceLister struct {
	coreClientset           core_v1.Clientset
	appsClientset           apps_v1.Clientset
	gatewayClientset        gateway_v1.Clientset
	glooClientset           gloo_v1.Clientset
	enterpriseGlooClientset enterprise_gloo_v1.Clientset
	ratelimitClientset      ratelimit_v1alpha1.Clientset
}

// Gets a list of Gloo instances on this cluster.
// If SingleInstanceOnly is true, only returns the Gloo instance in the same install namespace as
// the currently running apiserver instance (rather than all Gloo instances on the cluster).
func (l *singleClusterGlooInstanceLister) ListGlooInstances(ctx context.Context) ([]*rpc_edge_v1.GlooInstance, error) {
	var instances []*rpc_edge_v1.GlooInstance
	installNamespace := apiserverutils.GetInstallNamespace()

	// if we only want the Gloo instance in the current install namespace, we can get it directly
	// in the event that a use has requested limited RBAC restrictions, we have only a single gloo instance
	if SingleInstanceOnly || os.Getenv("NAMESPACE_RESTRICTED_MODE") == "true" {
		instance, err := l.GetGlooInstance(ctx, &skv2_v1.ObjectRef{
			Name:      images.ControlPlaneDeploymentName,
			Namespace: installNamespace,
		})
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
		return instances, nil
	}

	deployments, err := l.appsClientset.Deployments().ListDeployment(ctx)
	if err != nil {
		return nil, err
	}
	daemonSets, err := l.appsClientset.DaemonSets().ListDaemonSet(ctx)
	if err != nil {
		return nil, err
	}
	services, err := l.coreClientset.Services().ListService(ctx)
	if err != nil {
		return nil, err
	}
	pods, err := l.coreClientset.Pods().ListPod(ctx)
	if err != nil {
		return nil, err
	}

	// go through all the deployments to find the Gloo control planes
	for _, deployment := range deployments.Items {
		deployment := deployment

		// get info about the control plane
		tag, isEnterprise, isControlPlane := images.GetControlPlaneImage(&deployment)
		if !isControlPlane {
			continue
		}

		instance, err := l.buildGlooInstanceFromDeployment(ctx, &deployment, tag, isEnterprise, deployments, daemonSets, services, pods)
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	sortGlooInstances(instances)
	return instances, nil
}

func (l *singleClusterGlooInstanceLister) GetGlooInstance(ctx context.Context, glooInstanceRef *skv2_v1.ObjectRef) (*rpc_edge_v1.GlooInstance, error) {
	if os.Getenv("NAMESPACE_RESTRICTED_MODE") == "true" {
		glooInstance, err := l.buildGlooInstanceFromInstanceRef(ctx, glooInstanceRef)
		if err != nil {
			return nil, err
		}
		return glooInstance, nil
	}

	// get the deployment with the same name/namespace as the glooInstanceRef
	deployment, err := l.appsClientset.Deployments().GetDeployment(ctx, ezkube.MakeClientObjectKey(glooInstanceRef))
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorw("Failed to get deployment for GlooInstance", zap.Error(err), zap.Any("glooInstanceRef", glooInstanceRef))
		return nil, err
	}
	tag, isEnterprise, isControlPlane := images.GetControlPlaneImage(deployment)
	if !isControlPlane {
		return nil, eris.Errorf("Ref does not refer to a GlooInstance: %v", glooInstanceRef)
	}

	deployments, err := l.appsClientset.Deployments().ListDeployment(ctx)
	if err != nil {
		return nil, err
	}
	daemonSets, err := l.appsClientset.DaemonSets().ListDaemonSet(ctx)
	if err != nil {
		return nil, err
	}
	services, err := l.coreClientset.Services().ListService(ctx)
	if err != nil {
		return nil, err
	}
	pods, err := l.coreClientset.Pods().ListPod(ctx)
	if err != nil {
		return nil, err
	}

	glooInstance, err := l.buildGlooInstanceFromDeployment(ctx, deployment, tag, isEnterprise, deployments, daemonSets, services, pods)
	if err != nil {
		return nil, err
	}
	return glooInstance, nil
}

func (l *singleClusterGlooInstanceLister) buildGlooInstanceFromInstanceRef(ctx context.Context, glooInstanceRef *skv2_v1.ObjectRef) (*rpc_edge_v1.GlooInstance, error) {
	settingsClient, watchedNamespaces, err := l.getSettingsClient(ctx, glooInstanceRef.Namespace)
	if err != nil {
		return nil, err
	}

	instance := &rpc_edge_v1.GlooInstance{
		Metadata: &rpc_edge_v1.ObjectMeta{
			Name:      glooInstanceRef.Name,
			Namespace: glooInstanceRef.Namespace,
		},
		Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
			Cluster:      ClusterName,
			IsEnterprise: true,
			ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
				Version:           "indeterminable version",
				Namespace:         glooInstanceRef.Namespace,
				WatchedNamespaces: watchedNamespaces,
			},
			Region: "indeterminable region",
			Check:  &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{},
		},
	}

	l.addCrdSummaries(ctx, instance, watchedNamespaces, settingsClient)
	return instance, nil
}

func (l *singleClusterGlooInstanceLister) buildGlooInstanceFromDeployment(
	ctx context.Context,
	deployment *k8s_apps_types.Deployment,
	tag string,
	isEnterprise bool,
	deployments *k8s_apps_types.DeploymentList,
	daemonSets *k8s_apps_types.DaemonSetList,
	services *k8s_core_types.ServiceList,
	pods *k8s_core_types.PodList) (*rpc_edge_v1.GlooInstance, error) {

	// get locality and ingress info
	podClient := l.coreClientset.Pods()
	nodeClient := l.coreClientset.Nodes()
	ipFinder := locality.NewExternalIpFinder(ClusterName, podClient, nodeClient)
	localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
	region, err := localityFinder.GetRegion(ctx)
	if err != nil {
		// log the error but keep going; we just won't have region info
		contextutils.LoggerFrom(ctx).Warnw("Failed to get region of control plane", zap.Error(err), zap.Any("deployment", deployment))
	}

	// get the watched namespaces from this Gloo instance's settings
	settingsClient, watchedNamespaces, err := l.getSettingsClient(ctx, deployment.GetNamespace())
	if err != nil {
		return nil, err
	}

	instance := &rpc_edge_v1.GlooInstance{
		Metadata: &rpc_edge_v1.ObjectMeta{
			// use the gloo deployment name and namespace
			Name:      deployment.GetName(),
			Namespace: deployment.GetNamespace(),
		},
		Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
			Cluster:      ClusterName,
			IsEnterprise: isEnterprise,
			ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
				Version:           tag,
				Namespace:         deployment.GetNamespace(),
				WatchedNamespaces: watchedNamespaces,
			},
			Region: region,
			Check:  &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{},
		},
	}

	// proxy info
	instance.Spec.Proxies = convertProxies(proxy.GetGlooInstanceProxies(ctx, "", deployment.GetNamespace(),
		newDeploymentSetFromList(deployments), newDaemonSetSetFromList(daemonSets), newServiceSetFromList(services), localityFinder, ipFinder))

	// pod and deployment summaries
	podSet := newPodSetFromList(pods)
	deploymentSet := newDeploymentSetFromList(deployments)
	instance.Spec.Check.Pods = convertSummary(checker.GetPodsSummary(ctx, podSet, instance.Spec.ControlPlane.GetNamespace(), ""))
	instance.Spec.Check.Deployments = convertSummary(checker.GetDeploymentsSummary(ctx, deploymentSet, instance.Spec.ControlPlane.GetNamespace(), ""))

	l.addCrdSummaries(ctx, instance, watchedNamespaces, settingsClient)

	return instance, nil
}

func (l *singleClusterGlooInstanceLister) getSettingsClient(ctx context.Context, namespace string) (types.SettingsClient, []string, error) {
	settingsClient := l.glooClientset.Settings()
	settings, err := settingsClient.GetSettings(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      defaults.SettingsName,
	})
	if err != nil {
		return nil, nil, err
	}
	// get the watched namespaces from this Gloo instance's settings
	watchedNamespaces := settings.Spec.GetWatchNamespaces()

	return settingsClient, watchedNamespaces, nil
}

func (l *singleClusterGlooInstanceLister) addCrdSummaries(ctx context.Context, instance *rpc_edge_v1.GlooInstance, watchedNamespaces []string, settingsClient types.SettingsClient) {
	// gloo resource summaries
	instance.Spec.Check.Settings = gloo_resource_handler.GetSettingsSummary(ctx, settingsClient, watchedNamespaces)
	instance.Spec.Check.Upstreams = gloo_resource_handler.GetUpstreamSummary(ctx, l.glooClientset.Upstreams(), watchedNamespaces)
	instance.Spec.Check.UpstreamGroups = gloo_resource_handler.GetUpstreamGroupSummary(ctx, l.glooClientset.UpstreamGroups(), watchedNamespaces)
	instance.Spec.Check.Proxies = gloo_resource_handler.GetProxySummary(ctx, l.glooClientset.Proxies(), watchedNamespaces)
	instance.Spec.Check.AuthConfigs = enterprise_gloo_resource_handler.GetAuthConfigSummary(ctx, l.enterpriseGlooClientset.AuthConfigs(), watchedNamespaces)
	instance.Spec.Check.RateLimitConfigs = ratelimit_resource_handler.GetRateLimitConfigSummary(ctx, l.ratelimitClientset.RateLimitConfigs(), watchedNamespaces)

	// gateway resource summaries
	instance.Spec.Check.Gateways = gateway_resource_handler.GetGatewaySummary(ctx, l.gatewayClientset.Gateways(), watchedNamespaces)
	instance.Spec.Check.VirtualServices = gateway_resource_handler.GetVirtualServiceSummary(ctx, l.gatewayClientset.VirtualServices(), watchedNamespaces)
	instance.Spec.Check.RouteTables = gateway_resource_handler.GetRouteTableSummary(ctx, l.gatewayClientset.RouteTables(), watchedNamespaces)
	instance.Spec.Check.MatchableHttpGateways = gateway_resource_handler.GetMatchableHttpGatewaySummary(ctx, l.gatewayClientset.MatchableHttpGateways(), watchedNamespaces)
}

// Convert from the generic structs to the GlooEE apiserver types

func convertProxies(proxies []*proxy.Proxy) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy {
	convertedProxies := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy, 0, len(proxies))
	for _, p := range proxies {
		convertedProxies = append(convertedProxies, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{
			Replicas:                      p.Replicas,
			AvailableReplicas:             p.AvailableReplicas,
			ReadyReplicas:                 p.ReadyReplicas,
			WasmEnabled:                   p.WasmEnabled,
			ReadConfigMulticlusterEnabled: p.ReadConfigMulticlusterEnabled,
			Version:                       p.Version,
			Name:                          p.Name,
			Namespace:                     p.Namespace,
			WorkloadControllerType:        convertWorkloadControllerType(p.WorkloadControllerType),
			Zones:                         p.Zones,
			IngressEndpoints:              convertIngressEndpoints(p.IngressEndpoints),
		})
	}
	return convertedProxies
}

func convertWorkloadControllerType(workloadControllerType proxy.WorkloadController) rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_WorkloadController {
	switch workloadControllerType {
	case proxy.DEPLOYMENT:
		return rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_DEPLOYMENT
	case proxy.DAEMON_SET:
		return rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_DAEMON_SET
	default:
		return rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_UNDEFINED
	}
}

func convertIngressEndpoints(ingressEndpoints []*locality.IngressEndpoint) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint {
	convertedEndpoints := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint, 0, len(ingressEndpoints))
	for _, endpoint := range ingressEndpoints {
		convertedEndpoints = append(convertedEndpoints, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint{
			Address:     endpoint.Address,
			Ports:       convertPorts(endpoint.Ports),
			ServiceName: endpoint.ServiceName,
		})
	}
	return convertedEndpoints
}

func convertPorts(ports []*locality.Port) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port {
	convertedPorts := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port, 0, len(ports))
	for _, port := range ports {
		convertedPorts = append(convertedPorts, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port{
			Port: port.Port,
			Name: port.Name,
		})
	}
	return convertedPorts
}

func convertSummary(summary *checker.Summary) *rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary {
	return &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
		Total:    summary.Total,
		Errors:   convertResourceReports(summary.Errors),
		Warnings: convertResourceReports(summary.Warnings),
	}
}

func convertResourceReports(resourceReports []*checker.ResourceReport) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport {
	if len(resourceReports) == 0 {
		return nil
	}
	convertedResourceReports := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport, 0, len(resourceReports))
	for _, resourceReport := range resourceReports {
		convertedResourceReports = append(convertedResourceReports, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
			Ref:     resourceReport.Ref,
			Message: resourceReport.Message,
		})
	}
	return convertedResourceReports
}

// Helpers to replace New<T>SetFromList functions from external-apis to use generic sets
// Take a list of <T> as returned from the k8s client and return a set of pointers
func newDeploymentSetFromList(deploymentList *k8s_apps_types.DeploymentList) sk_sets.ResourceSet[*k8s_apps_types.Deployment] {
	set := sk_sets.NewResourceSet[*k8s_apps_types.Deployment]()
	for _, item := range deploymentList.Items {
		set.Insert(&item)
	}
	return set
}
func newPodSetFromList(serviceList *k8s_core_types.PodList) sk_sets.ResourceSet[*k8s_core_types.Pod] {
	set := sk_sets.NewResourceSet[*k8s_core_types.Pod]()
	for _, item := range serviceList.Items {
		set.Insert(&item)
	}
	return set
}
func newServiceSetFromList(serviceList *k8s_core_types.ServiceList) sk_sets.ResourceSet[*k8s_core_types.Service] {
	set := sk_sets.NewResourceSet[*k8s_core_types.Service]()
	for _, item := range serviceList.Items {
		set.Insert(&item)
	}
	return set
}
func newDaemonSetSetFromList(deploymentList *k8s_apps_types.DaemonSetList) sk_sets.ResourceSet[*k8s_apps_types.DaemonSet] {
	set := sk_sets.NewResourceSet[*k8s_apps_types.DaemonSet]()
	for _, item := range deploymentList.Items {
		set.Insert(&item)
	}
	return set
}
