package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	apiv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
)

func TestGatewayClassProvisioner(t *testing.T) {
	tests := []struct {
		name           string
		existingGCs    []client.Object
		expectedGCs    map[string]*apiv1.GatewayClass
		controllerName string
	}{
		{
			name:           "creates missing gateway classes",
			existingGCs:    nil,
			controllerName: "kgateway.dev/controller",
			expectedGCs: map[string]*apiv1.GatewayClass{
				wellknown.GatewayClassName: {
					ObjectMeta: metav1.ObjectMeta{
						Name: wellknown.GatewayClassName,
					},
					Spec: apiv1.GatewayClassSpec{
						ControllerName: "kgateway.dev/controller",
					},
				},
				wellknown.WaypointClassName: {
					ObjectMeta: metav1.ObjectMeta{
						Name: wellknown.WaypointClassName,
						Annotations: map[string]string{
							"ambient.istio.io/waypoint-inbound-binding": "PROXY/15088",
						},
					},
					Spec: apiv1.GatewayClassSpec{
						ControllerName: "kgateway.dev/controller",
					},
				},
			},
		},
		{
			name: "skips existing gateway class with different controller",
			existingGCs: []client.Object{
				&apiv1.GatewayClass{
					ObjectMeta: metav1.ObjectMeta{
						Name: wellknown.GatewayClassName,
					},
					Spec: apiv1.GatewayClassSpec{
						ControllerName: "different.controller/name",
					},
				},
			},
			controllerName: "kgateway.dev/controller",
			expectedGCs: map[string]*apiv1.GatewayClass{
				wellknown.GatewayClassName: {
					ObjectMeta: metav1.ObjectMeta{
						Name: wellknown.GatewayClassName,
					},
					Spec: apiv1.GatewayClassSpec{
						ControllerName: "different.controller/name",
					},
				},
				wellknown.WaypointClassName: {
					ObjectMeta: metav1.ObjectMeta{
						Name: wellknown.WaypointClassName,
						Annotations: map[string]string{
							"ambient.istio.io/waypoint-inbound-binding": "PROXY/15088",
						},
					},
					Spec: apiv1.GatewayClassSpec{
						ControllerName: "kgateway.dev/controller",
					},
				},
			},
		},
		{
			name: "updates missing annotations",
			existingGCs: []client.Object{
				&apiv1.GatewayClass{
					ObjectMeta: metav1.ObjectMeta{
						Name: wellknown.WaypointClassName,
					},
					Spec: apiv1.GatewayClassSpec{
						ControllerName: "kgateway.dev/controller",
					},
				},
			},
			controllerName: "kgateway.dev/controller",
			expectedGCs: map[string]*apiv1.GatewayClass{
				wellknown.GatewayClassName: {
					ObjectMeta: metav1.ObjectMeta{
						Name: wellknown.GatewayClassName,
					},
					Spec: apiv1.GatewayClassSpec{
						ControllerName: "kgateway.dev/controller",
					},
				},
				wellknown.WaypointClassName: {
					ObjectMeta: metav1.ObjectMeta{
						Name: wellknown.WaypointClassName,
						Annotations: map[string]string{
							"ambient.istio.io/waypoint-inbound-binding": "PROXY/15088",
						},
					},
					Spec: apiv1.GatewayClassSpec{
						ControllerName: "kgateway.dev/controller",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithObjects(tt.existingGCs...).
				Build()

			r := &gatewayClassProvisioner{
				client:         client,
				controllerName: tt.controllerName,
			}

			ctx := context.Background()
			_, err := r.Reconcile(ctx, reconcile.Request{})
			require.NoError(t, err)

			// Verify expected GatewayClasses
			for name, expected := range tt.expectedGCs {
				actual := &apiv1.GatewayClass{}
				err := client.Get(ctx, types.NamespacedName{Name: name}, actual)
				require.NoError(t, err)

				assert.Equal(t, expected.Spec.ControllerName, actual.Spec.ControllerName)
				assert.Equal(t, expected.Annotations, actual.Annotations)
				assert.Equal(t, expected.Labels, actual.Labels)
			}
		})
	}
}

// TestGatewayClassProvisionerIntegration tests the controller with a real API server
func TestGatewayClassProvisionerIntegration(t *testing.T) {
	// Setup envtest
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			"../../../install/helm/kgateway/crds",
		},
	}

	cfg, err := testEnv.Start()
	require.NoError(t, err)
	defer func() { require.NoError(t, testEnv.Stop()) }()

	err = apiv1.AddToScheme(scheme.Scheme)
	require.NoError(t, err)

	// Create a new manager
	mgr, err := manager.New(cfg, manager.Options{})
	require.NoError(t, err)

	// Start the manager
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := mgr.Start(ctx)
		require.NoError(t, err)
	}()

	// Create the controller
	err = NewGatewayClassProvisioner(mgr, "kgateway.dev/controller")
	require.NoError(t, err)

	// Wait for controller to create GatewayClasses
	time.Sleep(time.Second)

	// Verify GatewayClasses were created
	client := mgr.GetClient()
	for _, name := range []string{wellknown.GatewayClassName, wellknown.WaypointClassName} {
		gc := &apiv1.GatewayClass{}
		err := client.Get(ctx, types.NamespacedName{Name: name}, gc)
		require.NoError(t, err)

		assert.Equal(t, apiv1.GatewayController("kgateway.dev/controller"), gc.Spec.ControllerName)
		if name == wellknown.WaypointClassName {
			assert.Equal(t, "PROXY/15088", gc.Annotations["ambient.istio.io/waypoint-inbound-binding"])
		}
	}
}
