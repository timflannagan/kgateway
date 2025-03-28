package controller_test

import (
	"context"
	"time"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/controller"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	apiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// TODO: Handle edge case where CRD is not present on the cluster. Is this even possible
// given we wait until the CRD is present before starting the controller?

var _ = Describe("GatewayClass Status Controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	var (
		cancel           context.CancelFunc
		supportedVersion string
	)

	BeforeEach(func() {
		var err error
		cancel, err = createManager(ctx, nil, nil)
		Expect(err).NotTo(HaveOccurred())

		supportedVersion = controller.GetSupportedVersions().UnsortedList()[1]
	})

	AfterEach(func() {
		if cancel != nil {
			cancel()
		}
		// ensure goroutines cleanup
		Eventually(func() bool { return true }).WithTimeout(3 * time.Second).Should(BeTrue())
	})

	Context("GatewayClass reconciliation", func() {
		var (
			gc *apiv1.GatewayClass
		)
		BeforeEach(func() {
			gc = &apiv1.GatewayClass{}
			Eventually(func() bool {
				if err := k8sClient.Get(ctx, types.NamespacedName{Name: gatewayClassName}, gc); err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue(), "GatewayClass %s not found", gatewayClassName)
		})

		It("should set the Accepted=True condition type", func() {
			Expect(gc.Status.Conditions).Should(HaveLen(2))
			Expect(gc.Status.Conditions[0].Type).Should(Equal(string(apiv1.GatewayClassConditionStatusAccepted)))
			Expect(gc.Status.Conditions[0].Status).Should(Equal(metav1.ConditionTrue))
			Expect(gc.Status.Conditions[0].Reason).Should(Equal(string(apiv1.GatewayClassReasonAccepted)))
			Expect(gc.Status.Conditions[0].Message).Should(ContainSubstring(`accepted by kgateway controller`))
		})

		It("should set the SupportedVersion=True condition type", func() {
			Expect(gc.Status.Conditions).Should(HaveLen(2))
			Expect(gc.Status.Conditions[1].Type).Should(Equal(string(apiv1.GatewayClassConditionStatusSupportedVersion)))
			Expect(gc.Status.Conditions[1].Status).Should(Equal(metav1.ConditionTrue))
			Expect(gc.Status.Conditions[1].Reason).Should(Equal(string(apiv1.GatewayClassReasonSupportedVersion)))
			Expect(gc.Status.Conditions[1].Message).Should(ContainSubstring(`supported by kgateway controller`))
		})
	})

	Context("CRD reconciliation", func() {
		When("the CRD's bundle annotation is not set", func() {
			var (
				crd *apiextensionsv1.CustomResourceDefinition
			)
			BeforeEach(func() {
				crd = &apiextensionsv1.CustomResourceDefinition{}
				Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "gatewayclasses.gateway.networking.k8s.io"}, crd)).Should(Succeed())
				delete(crd.Annotations, "gateway.networking.k8s.io/bundle-version")
				Eventually(func() error {
					return k8sClient.Update(ctx, crd)
				}, timeout, interval).Should(Succeed())
			})
			AfterEach(func() {
				crd.Annotations["gateway.networking.k8s.io/bundle-version"] = supportedVersion
				Eventually(func() error {
					return k8sClient.Update(ctx, crd)
				}, timeout, interval).Should(Succeed())
			})

			It("should set the UnsupportedVersion condition type", func() {
				Eventually(func() (*metav1.Condition, error) {
					gc := &apiv1.GatewayClass{}
					if err := k8sClient.Get(ctx, types.NamespacedName{Name: gatewayClassName}, gc); err != nil {
						return nil, err
					}
					return meta.FindStatusCondition(gc.Status.Conditions, string(apiv1.GatewayClassConditionStatusSupportedVersion)), nil
				}, timeout, interval).Should(And(
					Not(BeNil()),
					WithTransform(func(c *metav1.Condition) string { return c.Type }, Equal(string(apiv1.GatewayClassConditionStatusSupportedVersion))),
					WithTransform(func(c *metav1.Condition) metav1.ConditionStatus { return c.Status }, Equal(metav1.ConditionFalse)),
					WithTransform(func(c *metav1.Condition) string { return c.Reason }, Equal(string(apiv1.GatewayClassReasonUnsupportedVersion))),
					WithTransform(func(c *metav1.Condition) string { return c.Message }, ContainSubstring(`not found`)),
				))
			})
		})

		When("an unsupported Gateway API version exists", func() {
			var (
				crd             *apiextensionsv1.CustomResourceDefinition
				previousVersion string
			)
			BeforeEach(func() {
				crd = &apiextensionsv1.CustomResourceDefinition{}
				Expect(k8sClient.Get(ctx, types.NamespacedName{Name: "gatewayclasses.gateway.networking.k8s.io"}, crd)).Should(Succeed())
				previousVersion = crd.Annotations["gateway.networking.k8s.io/bundle-version"]
				crd.Annotations["gateway.networking.k8s.io/bundle-version"] = "v1.0.0"
				Eventually(func() error {
					return k8sClient.Update(ctx, crd)
				}, timeout, interval).Should(Succeed())
			})
			AfterEach(func() {
				crd.Annotations["gateway.networking.k8s.io/bundle-version"] = previousVersion
				Eventually(func() error {
					return k8sClient.Update(ctx, crd)
				}, timeout, interval).Should(Succeed())
			})

			It("should reconcile CRDs correctly", func() {
				By("verifying the SupportedVersion=True condition type was set correctly")
				Eventually(func() (*metav1.Condition, error) {
					gc := &apiv1.GatewayClass{}
					if err := k8sClient.Get(ctx, types.NamespacedName{Name: altGatewayClassName}, gc); err != nil {
						return nil, err
					}
					return meta.FindStatusCondition(gc.Status.Conditions, string(apiv1.GatewayClassConditionStatusSupportedVersion)), nil
				}, timeout, interval).Should(And(
					Not(BeNil()),
					WithTransform(func(c *metav1.Condition) string { return c.Type }, Equal(string(apiv1.GatewayClassConditionStatusSupportedVersion))),
					WithTransform(func(c *metav1.Condition) metav1.ConditionStatus { return c.Status }, Equal(metav1.ConditionFalse)),
					WithTransform(func(c *metav1.Condition) string { return c.Reason }, Equal(string(apiv1.GatewayClassReasonUnsupportedVersion))),
					WithTransform(func(c *metav1.Condition) string { return c.Message }, ContainSubstring(`not supported by kgateway controller`)),
				))

				By("updating the CRD to a supported version")
				crd.Annotations["gateway.networking.k8s.io/bundle-version"] = supportedVersion
				Eventually(func() error {
					return k8sClient.Update(ctx, crd)
				}, timeout, interval).Should(Succeed())

				By("verifying the SupportedVersion=True condition type was set correctly")
				Eventually(func() (*metav1.Condition, error) {
					gc := &apiv1.GatewayClass{}
					if err := k8sClient.Get(ctx, types.NamespacedName{Name: altGatewayClassName}, gc); err != nil {
						return nil, err
					}
					return meta.FindStatusCondition(gc.Status.Conditions, string(apiv1.GatewayClassConditionStatusSupportedVersion)), nil
				}, timeout, interval).Should(And(
					Not(BeNil()),
					WithTransform(func(c *metav1.Condition) string { return c.Type }, Equal(string(apiv1.GatewayClassConditionStatusSupportedVersion))),
					WithTransform(func(c *metav1.Condition) metav1.ConditionStatus { return c.Status }, Equal(metav1.ConditionTrue)),
					WithTransform(func(c *metav1.Condition) string { return c.Reason }, Equal(string(apiv1.GatewayClassReasonSupportedVersion))),
					WithTransform(func(c *metav1.Condition) string { return c.Message }, ContainSubstring(`supported by kgateway controller`)),
				))
			})
		})
	})
})
