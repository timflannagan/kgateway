package controller

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	apiv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
)

// gatewayClassProvisioner reconciles GatewayClass objects to ensure they exist
type gatewayClassProvisioner struct {
	client         client.Client
	controllerName string
}

// NewGatewayClassProvisioner creates a new GatewayClassProvisioner
func NewGatewayClassProvisioner(mgr ctrl.Manager, controllerName string) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.GatewayClass{}).
		Complete(&gatewayClassProvisioner{
			client:         mgr.GetClient(),
			controllerName: controllerName,
		})
}

func (r *gatewayClassProvisioner) Reconcile(ctx context.Context, _ ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	type config struct {
		labels      map[string]string
		annotations map[string]string
	}

	// TODO: we just need labels/annotations here.
	classes := map[string]config{
		wellknown.GatewayClassName: {},
		wellknown.WaypointClassName: {
			annotations: map[string]string{
				"ambient.istio.io/waypoint-inbound-binding": "PROXY/15088",
			},
		},
	}

	// Need to ensure the following requirements:
	// 1. GatewayClass exists
	// 2. GatewayClass has the correct controller
	// 3. GatewayClass has the correct annotations
	// 4. GatewayClass has the correct labels
	//
	// Any parametersRef will be handled by users. If they modify the default GC
	// that's managed by us, then we will overwrite their changes. Is that the
	// right behavior? OpenShift had a create-only toggle for the openshift to
	// handle this scenario, but imo, users can disable it in the helm values,
	// upgrade with the new values, and workaround it.
	for name := range classes {
		gc := &apiv1.GatewayClass{}
		err := r.client.Get(ctx, client.ObjectKey{Name: name}, gc)
		if apierrors.IsNotFound(err) {
			// GatewayClass does not exist, create it
			if err := r.client.Create(ctx, &apiv1.GatewayClass{
				ObjectMeta: ctrl.ObjectMeta{
					Name:        name,
					Annotations: classes[name].annotations,
					Labels:      classes[name].labels,
				},
				Spec: apiv1.GatewayClassSpec{
					ControllerName: apiv1.GatewayController(r.controllerName),
				},
			}); err != nil {
				return ctrl.Result{}, err
			}
			log.Info("created GatewayClass", "name", name)
			continue
		}
		if err != nil {
			return ctrl.Result{}, err
		}

		// TODO: Verify the desired and current annotations/labels are the same
		// If not, update the GatewayClass.

		// GatewayClass exists, ensure it has the correct controller
		if gc.Spec.ControllerName != apiv1.GatewayController(r.controllerName) {
			log.Info("skipping existing GatewayClass with different controller",
				"name", name,
				"controller", gc.Spec.ControllerName,
			)
			continue
		}
	}

	return ctrl.Result{}, nil
}

// hasRequiredAnnotations checks if obj has all required annotations
func hasRequiredAnnotations(existing, required map[string]string) bool {
	for k, v := range required {
		if existing[k] != v {
			return false
		}
	}
	return true
}

// mergeAnnotations merges required annotations into existing ones
func mergeAnnotations(existing, required map[string]string) map[string]string {
	if existing == nil {
		existing = make(map[string]string)
	}
	for k, v := range required {
		existing[k] = v
	}
	return existing
}
