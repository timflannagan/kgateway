package deployer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"istio.io/istio/pkg/kube/krt/krttest"
	"istio.io/istio/pkg/test"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	infextv1a2 "sigs.k8s.io/gateway-api-inference-extension/api/v1alpha2"
	api "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"

	gw2_v1alpha1 "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	common "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/collections"
	"github.com/kgateway-dev/kgateway/v2/pkg/schemes"
)

const (
	defaultNamespace = "default"
)

type testHelmValuesGenerator struct{}

func (thv *testHelmValuesGenerator) GetValues(ctx context.Context, gw *api.Gateway, inputs *Inputs) (map[string]any, error) {
	return map[string]any{
		"testHelmValuesGenerator": struct{}{},
	}, nil
}

func TestShouldUseDefaultGatewayParameters(t *testing.T) {
	gwc := defaultGatewayClass()
	gwParams := emptyGatewayParameters()

	gw := &api.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: defaultNamespace,
			UID:       "1235",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Gateway",
			APIVersion: "gateway.networking.k8s.io",
		},
		Spec: api.GatewaySpec{
			GatewayClassName: wellknown.GatewayClassName,
		},
	}

	gwp := NewGatewayParameters(newFakeClientWithObjs(gwc, gwParams))
	vals, err := gwp.GetValues(context.Background(), gw, defaultInputs(t, gwc, gw))

	assert.NoError(t, err)
	assert.Contains(t, vals, "gateway")
}

func TestShouldUseExtendedGatewayParameters(t *testing.T) {
	gwc := defaultGatewayClass()
	gwParams := emptyGatewayParameters()
	extraGwParams := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Namespace: defaultNamespace}}

	gw := &api.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: defaultNamespace,
			UID:       "1235",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Gateway",
			APIVersion: "gateway.networking.k8s.io",
		},
		Spec: api.GatewaySpec{
			Infrastructure: &api.GatewayInfrastructure{
				ParametersRef: &api.LocalParametersReference{
					Group: "v1",
					Kind:  "ConfigMap",
					Name:  "testing",
				},
			},
			GatewayClassName: wellknown.GatewayClassName,
		},
	}

	gwp := NewGatewayParameters(newFakeClientWithObjs(gwc, gwParams, extraGwParams)).
		WithExtraGatewayParameters(ExtraGatewayParameters{Group: "v1", Kind: "ConfigMap", Object: extraGwParams, Generator: &testHelmValuesGenerator{}})
	vals, err := gwp.GetValues(context.Background(), gw, defaultInputs(t, gwc, gw))

	assert.NoError(t, err)
	assert.Contains(t, vals, "testHelmValuesGenerator")
}

func defaultGatewayClass() *api.GatewayClass {
	return &api.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: wellknown.GatewayClassName,
		},
		Spec: api.GatewayClassSpec{
			ControllerName: wellknown.GatewayControllerName,
			ParametersRef: &api.ParametersReference{
				Group:     gw2_v1alpha1.GroupName,
				Kind:      api.Kind(wellknown.GatewayParametersGVK.Kind),
				Name:      wellknown.DefaultGatewayParametersName,
				Namespace: ptr.To(api.Namespace(defaultNamespace)),
			},
		},
	}
}

func emptyGatewayParameters() *gw2_v1alpha1.GatewayParameters {
	return &gw2_v1alpha1.GatewayParameters{
		TypeMeta: metav1.TypeMeta{
			Kind: wellknown.GatewayParametersGVK.Kind,
			// The parsing expects GROUP/VERSION format in this field
			APIVersion: gw2_v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      wellknown.DefaultGatewayParametersName,
			Namespace: defaultNamespace,
			UID:       "1237",
		},
	}
}

func defaultInputs(t *testing.T, gwc *api.GatewayClass, gw *api.Gateway) *Inputs {
	return &Inputs{
		CommonCollections: newCommonCols(t, gwc, gw),
		ControllerName:    wellknown.GatewayControllerName,
		Dev:               false,
		ControlPlane: ControlPlaneInfo{
			XdsHost: "something.cluster.local",
			XdsPort: 1234,
		},
		ImageInfo: &ImageInfo{
			Registry: "foo",
			Tag:      "bar",
		},
	}
}

// initialize a fake controller-runtime client with the given list of objects
func newFakeClientWithObjs(objs ...client.Object) client.Client {
	scheme := schemes.GatewayScheme()

	// Ensure the rbac types are registered.
	if err := rbacv1.AddToScheme(scheme); err != nil {
		panic(fmt.Sprintf("failed to add rbacv1 scheme: %v", err))
	}

	// Check if any object is an InferencePool, and add its scheme if needed.
	for _, obj := range objs {
		gvk := obj.GetObjectKind().GroupVersionKind()
		if gvk.Kind == wellknown.InferencePoolKind {
			if err := infextv1a2.AddToScheme(scheme); err != nil {
				panic(fmt.Sprintf("failed to add InferenceExtension scheme: %v", err))
			}
			break
		}
	}

	return fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		Build()
}

func newCommonCols(t test.Failer, initObjs ...client.Object) *common.CommonCollections {
	ctx := context.Background()
	var anys []any
	for _, obj := range initObjs {
		anys = append(anys, obj)
	}
	mock := krttest.NewMock(t, anys)

	policies := krtcollections.NewPolicyIndex(krtutil.KrtOptions{}, extensionsplug.ContributesPolicies{})
	kubeRawGateways := krttest.GetMockCollection[*api.Gateway](mock)
	kubeRawListenerSets := krttest.GetMockCollection[*apixv1a1.XListenerSet](mock)
	gatewayClasses := krttest.GetMockCollection[*api.GatewayClass](mock)
	nsCol := krtcollections.NewNamespaceCollectionFromCol(ctx, krttest.GetMockCollection[*corev1.Namespace](mock), krtutil.KrtOptions{})

	krtopts := krtutil.NewKrtOptions(ctx.Done(), nil)
	gateways := krtcollections.NewGatewayIndex(krtopts, wellknown.GatewayControllerName, policies, kubeRawGateways, kubeRawListenerSets, gatewayClasses, nsCol)

	commonCols := &common.CommonCollections{
		GatewayIndex: gateways,
	}

	for !kubeRawGateways.HasSynced() || !kubeRawListenerSets.HasSynced() || !gatewayClasses.HasSynced() {
		time.Sleep(time.Second / 10)
	}

	gateways.Gateways.WaitUntilSynced(ctx.Done())
	return commonCols
}
