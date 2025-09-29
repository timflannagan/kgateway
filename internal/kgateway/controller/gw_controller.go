package controller

import (
	"context"
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	api "sigs.k8s.io/gateway-api/apis/v1"

	intdeployer "github.com/kgateway-dev/kgateway/v2/internal/kgateway/deployer"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/deployer"
)

type gatewayReconciler struct {
	cli client.Client

	controllerName    string
	agwControllerName string

	scheme   *runtime.Scheme
	deployer *deployer.Deployer
}

func NewGatewayReconciler(ctx context.Context, cfg GatewayConfig, deployer *deployer.Deployer) *gatewayReconciler {
	return &gatewayReconciler{
		cli:               cfg.Mgr.GetClient(),
		scheme:            cfg.Mgr.GetScheme(),
		controllerName:    cfg.ControllerName,
		agwControllerName: cfg.AgwControllerName,
		deployer:          deployer,
	}
}

// Reconcile handles Gateway infrastructure deployment and manages the Programmed condition.
// It works alongside the Status Syncer, which handles translation and the Accepted condition.
//
// For managed gateways, this controller deploys Kubernetes resources (Deployments, Services)
// and reports infrastructure status. For self-managed gateways, it validates configuration
// and sets appropriate status without deploying anything.
//
// The controller owns two parts of Gateway status:
//   - Programmed condition: true when infrastructure is ready, false for config errors
//   - Addresses field: network endpoints from deployed services (managed gateways only)
//
// The Status Syncer owns the Accepted condition and Listeners field, avoiding the race
// conditions that occurred when both components tried to manage the same status fields.
func (r *gatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, rErr error) {
	log := log.FromContext(ctx).WithValues("gw", req.NamespacedName)
	log.V(1).Info("reconciling request", "req", req)

	finishMetrics := collectReconciliationMetrics("gateway", req)
	defer func() {
		finishMetrics(rErr)
	}()

	var gw api.Gateway
	if err := r.cli.Get(ctx, req.NamespacedName, &gw); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if gw.GetDeletionTimestamp() != nil {
		// no need to do anything as we have owner refs, so children will be deleted
		log.Info("gateway deleted, no need for reconciling")
		return ctrl.Result{}, nil
	}

	// make sure we're the right controller for this
	var gwc api.GatewayClass
	if err := r.cli.Get(
		ctx,
		client.ObjectKey{Name: string(gw.Spec.GatewayClassName)},
		&gwc,
	); err != nil {
		log.Error(err, "failed to check controller for GatewayClass")
		return ctrl.Result{}, err
	}
	if gwc.Spec.ControllerName != api.GatewayController(r.controllerName) && gwc.Spec.ControllerName != api.GatewayController(r.agwControllerName) {
		// ignore, not our GatewayClass
		return ctrl.Result{}, nil
	}

	// at this point, we now need to prepare to deploy the gateway. shell out to the deployer
	// component to get the objects to deploy. this will additionally consider any parameterRefs
	// on the gateway to further customize the objects.
	objs, err := r.deployer.GetObjsToDeploy(ctx, &gw)
	if err != nil {
		if errors.Is(err, intdeployer.ErrNoValidPorts) {
			// status is reported from translator, so return normally
			return ctrl.Result{}, err
		}
		// if we fail to either reference a valid GatewayParameters or
		// the GatewayParameters configuration leads to issues building the
		// objects, we want to set the Programmed status to false with InvalidParameters.
		condition := metav1.Condition{
			Type:               string(api.GatewayConditionProgrammed),
			Status:             metav1.ConditionFalse,
			ObservedGeneration: gw.Generation,
			Reason:             string(api.GatewayReasonInvalidParameters),
			Message:            err.Error(),
		}
		patcher := utils.NewGatewayStatusPatcher(r.cli)
		if statusErr := patcher.PatchCondition(ctx, client.ObjectKeyFromObject(&gw), condition); statusErr != nil {
			log.Error(statusErr, "failed to update Gateway status")
			return ctrl.Result{}, statusErr
		}
		return ctrl.Result{}, err
	}

	objs = r.deployer.SetNamespaceAndOwner(&gw, objs)
	err = r.deployer.DeployObjsWithSource(ctx, objs, &gw)
	if err != nil {
		return ctrl.Result{}, err
	}

	// If we reach here, deployment was successful - set Programmed=True
	patcher := utils.NewGatewayStatusPatcher(r.cli)
	if statusErr := patcher.PatchCondition(ctx, client.ObjectKeyFromObject(&gw), metav1.Condition{
		Type:               string(api.GatewayConditionProgrammed),
		Status:             metav1.ConditionTrue,
		ObservedGeneration: gw.Generation,
		Reason:             string(api.GatewayReasonProgrammed),
		Message:            "Gateway is programmed",
	}); statusErr != nil {
		log.Error(statusErr, "failed to update Gateway status")
		return ctrl.Result{}, statusErr
	}

	// find the name/ns of the service we own so we can grab addresses
	// from it for status
	var generatedSvc *metav1.ObjectMeta
	for _, obj := range objs {
		if svc, ok := obj.(*corev1.Service); ok {
			generatedSvc = &svc.ObjectMeta
			break
		}
	}

	// update status (whether we generated a service or not, for unmanaged)
	result := ctrl.Result{}
	err = updateStatus(ctx, r.cli, &gw, generatedSvc)
	if err != nil {
		log.Error(err, "failed to update status")
		result.Requeue = true
	}

	return result, nil
}

func updateStatus(ctx context.Context, cli client.Client, gw *api.Gateway, svcmd *metav1.ObjectMeta) error {
	var svc *corev1.Service
	if svcmd != nil {
		svcnns := client.ObjectKey{
			Namespace: svcmd.Namespace,
			Name:      svcmd.Name,
		}

		svc = &corev1.Service{}
		if err := cli.Get(ctx, svcnns, svc); err != nil {
			return client.IgnoreNotFound(err)
		}

		// make sure we own this service
		controller := metav1.GetControllerOf(svc)
		if controller == nil {
			return nil
		}

		if gw.UID != controller.UID {
			return nil
		}
	}

	// update gateway addresses in the status
	desiredAddresses := getDesiredAddresses(gw, svc)

	return updateGatewayAddresses(ctx, cli, client.ObjectKeyFromObject(gw), desiredAddresses)
}

func getDesiredAddresses(gw *api.Gateway, svc *corev1.Service) []api.GatewayStatusAddress {
	var ret []api.GatewayStatusAddress
	seen := sets.New[api.GatewayStatusAddress]()

	if svc != nil && svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(svc.Status.LoadBalancer.Ingress) == 0 {
			return nil
		}
		for _, ing := range svc.Status.LoadBalancer.Ingress {
			if addr, ok := convertIngressAddr(ing); ok {
				seen.Insert(addr)
				ret = append(ret, addr)
			}
		}

		return ret
	} else if svc != nil {
		t := api.IPAddressType
		if len(svc.Spec.ClusterIPs) != 0 {
			for _, ip := range svc.Spec.ClusterIPs {
				ret = append(ret, api.GatewayStatusAddress{
					Type:  &t,
					Value: ip,
				})
			}
		} else if svc.Spec.ClusterIP != "" {
			ret = append(ret, api.GatewayStatusAddress{
				Type:  &t,
				Value: svc.Spec.ClusterIP,
			})
		}
	}

	for _, specAddr := range gw.Spec.Addresses {
		addr := api.GatewayStatusAddress{
			Type:  specAddr.Type,
			Value: specAddr.Value,
		}
		if !seen.Has(addr) {
			ret = append(ret, addr)
		}
	}

	return ret
}

// updateGatewayAddresses updates the addresses of a Gateway resource using strategic patching
// to avoid conflicts with other components updating different status fields.
func updateGatewayAddresses(
	ctx context.Context,
	cli client.Client,
	gwNN types.NamespacedName,
	desired []api.GatewayStatusAddress,
) error {
	// Use strategic patching to only update the addresses field
	patcher := utils.NewGatewayStatusPatcher(cli)
	return patcher.PatchAddresses(ctx, gwNN, desired)
}

func convertIngressAddr(ing corev1.LoadBalancerIngress) (api.GatewayStatusAddress, bool) {
	if ing.Hostname != "" {
		t := api.HostnameAddressType
		return api.GatewayStatusAddress{
			Type:  &t,
			Value: ing.Hostname,
		}, true
	}
	if ing.IP != "" {
		t := api.IPAddressType
		return api.GatewayStatusAddress{
			Type:  &t,
			Value: ing.IP,
		}, true
	}
	return api.GatewayStatusAddress{}, false
}
