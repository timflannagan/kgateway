package controller_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	api "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
)

var _ = Describe("GwController", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	DescribeTable(
		"should add status to gateway",
		func(gwClass string) {
			same := api.NamespacesFromSame
			gwName := "gw-" + gwClass
			gw := api.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      gwName,
					Namespace: "default",
				},
				Spec: api.GatewaySpec{
					GatewayClassName: api.ObjectName(gwClass),
					Listeners: []api.Listener{{
						Protocol: "HTTP",
						Port:     80,
						AllowedRoutes: &api.AllowedRoutes{
							Namespaces: &api.RouteNamespaces{
								From: &same,
							},
						},
						Name: "listener",
					}},
				},
			}
			err := k8sClient.Create(ctx, &gw)
			Expect(err).NotTo(HaveOccurred())

			// Wait for service to be created
			var svc corev1.Service
			Eventually(func() bool {
				var createdServices corev1.ServiceList
				err := k8sClient.List(ctx, &createdServices)
				if err != nil {
					return false
				}
				for _, svc = range createdServices.Items {
					if len(svc.ObjectMeta.OwnerReferences) == 1 && svc.ObjectMeta.OwnerReferences[0].UID == gw.UID {
						return true
					}
				}
				return false
			}, timeout, interval).Should(BeTrue(), "service not created")
			Expect(svc.Spec.ClusterIP).NotTo(BeEmpty())

			// Need to update the status of the service
			svc.Status.LoadBalancer = corev1.LoadBalancerStatus{
				Ingress: []corev1.LoadBalancerIngress{{
					IP: "127.0.0.1",
				}},
			}
			Expect(k8sClient.Status().Update(ctx, &svc)).NotTo(HaveOccurred())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, client.ObjectKey{Name: gwName, Namespace: "default"}, &gw)
				if err != nil {
					return false
				}
				if len(gw.Status.Addresses) == 0 {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(gw.Status.Addresses).To(HaveLen(1))
			Expect(*gw.Status.Addresses[0].Type).To(Equal(api.IPAddressType))
			Expect(gw.Status.Addresses[0].Value).To(Equal("127.0.0.1"))
		},
		Entry("default gateway class", gatewayClassName),
		Entry("alternative gateway class", altGatewayClassName),
	)

	Context("invalid GatewayParameters", func() {
		var (
			gwp *v1alpha1.GatewayParameters
			gw  api.Gateway
		)
		BeforeEach(func() {
			gwp = &v1alpha1.GatewayParameters{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "invalid-gwp",
					Namespace: "default",
				},
				Spec: v1alpha1.GatewayParametersSpec{
					Kube: &v1alpha1.KubernetesProxyConfig{
						Deployment: &v1alpha1.ProxyDeployment{
							Replicas: ptr.To(uint32(2)),
						},
					},
				},
			}
			err := k8sClient.Create(ctx, gwp)
			Expect(err).NotTo(HaveOccurred())

			gw = api.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gw",
					Namespace: "default",
				},
				Spec: api.GatewaySpec{
					GatewayClassName: api.ObjectName(gatewayClassName),
					Infrastructure: &api.GatewayInfrastructure{
						ParametersRef: &api.LocalParametersReference{
							Group: v1alpha1.GroupName,
							Kind:  "InvalidKindName",
							Name:  gwp.Name,
						},
					},
					Listeners: []api.Listener{{
						Name:     "listener",
						Protocol: "HTTP",
						Port:     80,
					}},
				},
			}
			err = k8sClient.Create(ctx, &gw)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, gwp)).NotTo(HaveOccurred())
			Expect(k8sClient.Delete(ctx, &gw)).NotTo(HaveOccurred())
		})

		It("should eventually set status Accepted to false with the InvalidParameters reason", func() {
			Eventually(func() (*metav1.Condition, error) {
				err := k8sClient.Get(ctx, client.ObjectKeyFromObject(&gw), &gw)
				if err != nil {
					return nil, err
				}
				return meta.FindStatusCondition(gw.Status.Conditions, string(api.GatewayClassConditionStatusAccepted)), nil
			}).Should(And(
				Not(BeNil()),
				WithTransform(func(c *metav1.Condition) string { return c.Type }, Equal(string(api.GatewayClassConditionStatusAccepted))),
				WithTransform(func(c *metav1.Condition) bool { return c.Status == metav1.ConditionFalse }, BeTrue()),
				WithTransform(func(c *metav1.Condition) string { return c.Reason }, Equal(string(api.GatewayClassReasonInvalidParameters))),
				WithTransform(func(c *metav1.Condition) string { return c.Message }, ContainSubstring(`invalid kind`)),
			))
		})
	})
})
