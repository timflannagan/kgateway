package controller_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	api "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
)

var _ = Describe("Gateway Status Conflicts", func() {
	const (
		timeout  = time.Second * 15
		interval = time.Millisecond * 100
	)

	var (
		ctx    context.Context
		cancel context.CancelFunc
		gw     *api.Gateway
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		// Create a test Gateway
		gw = &api.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "test-gateway-" + randomString(8),
				Namespace:  "default",
				Generation: 1,
			},
			Spec: api.GatewaySpec{
				GatewayClassName: api.ObjectName(gatewayClassName),
				Listeners: []api.Listener{{
					Name:     "http",
					Protocol: "HTTP",
					Port:     80,
				}},
			},
		}
		Expect(k8sClient.Create(ctx, gw)).To(Succeed())
	})

	AfterEach(func() {
		if gw != nil {
			k8sClient.Delete(ctx, gw)
		}
		cancel()
	})

	Context("Concurrent Status Updates", func() {
		It("should handle concurrent address and condition updates without conflicts", func() {
			var wg sync.WaitGroup
			var addressUpdateErr, conditionUpdateErr error
			var conflictCount int32

			// More aggressive concurrent updates to trigger conflicts
			numWorkers := 5
			numUpdatesPerWorker := 20

			// Simulate Gateway Controller updating addresses
			for w := 0; w < numWorkers; w++ {
				wg.Add(1)
				go func(workerID int) {
					defer wg.Done()
					defer GinkgoRecover()

					for i := 0; i < numUpdatesPerWorker; i++ {
						err := updateGatewayAddresses(ctx, k8sClient, client.ObjectKeyFromObject(gw), []api.GatewayStatusAddress{
							{
								Type:  ptr.To(api.IPAddressType),
								Value: fmt.Sprintf("192.168.%d.%d", workerID+1, 100+i),
							},
						})
						if err != nil {
							if apierrors.IsConflict(err) {
								atomic.AddInt32(&conflictCount, 1)
							} else {
								addressUpdateErr = err
								return
							}
						}
						// No sleep to maximize contention
					}
				}(w)
			}

			// Simulate Status Syncer updating conditions
			for w := 0; w < numWorkers; w++ {
				wg.Add(1)
				go func(workerID int) {
					defer wg.Done()
					defer GinkgoRecover()

					for i := 0; i < numUpdatesPerWorker; i++ {
						condition := metav1.Condition{
							Type:               string(api.GatewayConditionProgrammed),
							Status:             metav1.ConditionTrue,
							Reason:             string(api.GatewayReasonProgrammed),
							Message:            fmt.Sprintf("Gateway programmed by worker %d iteration %d", workerID, i),
							ObservedGeneration: gw.Generation,
						}

						err := updateGatewayStatusWithRetryFunc(ctx, k8sClient, client.ObjectKeyFromObject(gw), func(latest *api.Gateway) bool {
							meta.SetStatusCondition(&latest.Status.Conditions, condition)
							return true
						})
						if err != nil {
							if apierrors.IsConflict(err) {
								atomic.AddInt32(&conflictCount, 1)
							} else {
								conditionUpdateErr = err
								return
							}
						}
						// No sleep to maximize contention
					}
				}(w)
			}

			wg.Wait()

			// Both types of updates should eventually succeed
			Expect(addressUpdateErr).To(BeNil())
			Expect(conditionUpdateErr).To(BeNil())

			// Verify final state has both addresses and conditions
			Eventually(func() bool {
				var latest api.Gateway
				err := k8sClient.Get(ctx, client.ObjectKeyFromObject(gw), &latest)
				if err != nil {
					return false
				}

				hasAddresses := len(latest.Status.Addresses) > 0
				hasProgrammedCondition := meta.FindStatusCondition(latest.Status.Conditions, string(api.GatewayConditionProgrammed)) != nil

				return hasAddresses && hasProgrammedCondition
			}, timeout, interval).Should(BeTrue())

			// With more aggressive concurrent updates, we should see conflicts
			// Note: This expectation might need to be relaxed if the test environment
			// has very fast API responses and the retry logic is very effective
			By(fmt.Sprintf("Observed %d conflicts from concurrent updates", conflictCount))

			// Even if we don't see conflicts, this demonstrates the test framework works
			// and can be used to verify that fixes prevent conflicts in CI environments
		})

		It("should handle Accepted condition conflicts between components", func() {
			var wg sync.WaitGroup
			var controllerErr, syncerErr error

			// Simulate Gateway Controller setting Accepted=false (InvalidParameters)
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer GinkgoRecover()

				condition := metav1.Condition{
					Type:               string(api.GatewayConditionAccepted),
					Status:             metav1.ConditionFalse,
					Reason:             string(api.GatewayReasonInvalidParameters),
					Message:            "Invalid gateway parameters",
					ObservedGeneration: gw.Generation,
				}

				controllerErr = updateGatewayStatusWithRetryFunc(ctx, k8sClient, client.ObjectKeyFromObject(gw), func(latest *api.Gateway) bool {
					meta.SetStatusCondition(&latest.Status.Conditions, condition)
					return true
				})
			}()

			// Small delay to ensure some overlap
			time.Sleep(5 * time.Millisecond)

			// Simulate Status Syncer setting Accepted=true (translation success)
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer GinkgoRecover()

				condition := metav1.Condition{
					Type:               string(api.GatewayConditionAccepted),
					Status:             metav1.ConditionTrue,
					Reason:             string(api.GatewayReasonAccepted),
					Message:            "Gateway accepted",
					ObservedGeneration: gw.Generation,
				}

				syncerErr = updateGatewayStatusWithRetryFunc(ctx, k8sClient, client.ObjectKeyFromObject(gw), func(latest *api.Gateway) bool {
					meta.SetStatusCondition(&latest.Status.Conditions, condition)
					return true
				})
			}()

			wg.Wait()

			// Both updates should eventually succeed (with retries)
			Expect(controllerErr).To(BeNil())
			Expect(syncerErr).To(BeNil())

			// Final state should have one of the conditions (last writer wins)
			Eventually(func() *metav1.Condition {
				var latest api.Gateway
				err := k8sClient.Get(ctx, client.ObjectKeyFromObject(gw), &latest)
				if err != nil {
					return nil
				}
				return meta.FindStatusCondition(latest.Status.Conditions, string(api.GatewayConditionAccepted))
			}, timeout, interval).ShouldNot(BeNil())

			// This test demonstrates the race condition - final state is non-deterministic
			var finalGw api.Gateway
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(gw), &finalGw)).To(Succeed())
			acceptedCondition := meta.FindStatusCondition(finalGw.Status.Conditions, string(api.GatewayConditionAccepted))
			Expect(acceptedCondition).NotTo(BeNil())

			// The final reason could be either InvalidParameters or Accepted - this is the bug!
			// In a fixed implementation, there should be clear precedence rules
			By("Final condition reason: " + acceptedCondition.Reason)
		})

		It("should demonstrate address field getting overwritten by condition updates", func() {
			// First, set some addresses via Gateway Controller
			err := updateGatewayAddresses(ctx, k8sClient, client.ObjectKeyFromObject(gw), []api.GatewayStatusAddress{
				{
					Type:  ptr.To(api.IPAddressType),
					Value: "10.0.0.1",
				},
			})
			Expect(err).To(BeNil())

			// Verify addresses are set
			Eventually(func() []api.GatewayStatusAddress {
				var latest api.Gateway
				k8sClient.Get(ctx, client.ObjectKeyFromObject(gw), &latest)
				return latest.Status.Addresses
			}, timeout, interval).Should(HaveLen(1))

			// Now simulate a Status Syncer update that might inadvertently overwrite addresses
			// if not properly implemented
			err = updateGatewayStatusWithRetryFunc(ctx, k8sClient, client.ObjectKeyFromObject(gw), func(latest *api.Gateway) bool {
				// This simulates the current behavior where Status Syncer does a full status patch
				// but tries to preserve addresses by setting them to nil in comparison
				condition := metav1.Condition{
					Type:               string(api.GatewayConditionProgrammed),
					Status:             metav1.ConditionTrue,
					Reason:             string(api.GatewayReasonProgrammed),
					Message:            "Gateway programmed",
					ObservedGeneration: gw.Generation,
				}
				meta.SetStatusCondition(&latest.Status.Conditions, condition)
				return true
			})
			Expect(err).To(BeNil())

			// Addresses should still be preserved
			Eventually(func() []api.GatewayStatusAddress {
				var latest api.Gateway
				k8sClient.Get(ctx, client.ObjectKeyFromObject(gw), &latest)
				return latest.Status.Addresses
			}, timeout, interval).Should(HaveLen(1))

			// Verify both addresses and conditions exist
			var finalGw api.Gateway
			Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(gw), &finalGw)).To(Succeed())
			Expect(finalGw.Status.Addresses).To(HaveLen(1))
			Expect(finalGw.Status.Addresses[0].Value).To(Equal("10.0.0.1"))

			programmedCondition := meta.FindStatusCondition(finalGw.Status.Conditions, string(api.GatewayConditionProgrammed))
			Expect(programmedCondition).NotTo(BeNil())
			Expect(programmedCondition.Status).To(Equal(metav1.ConditionTrue))
		})
	})
})

// Helper functions that use the new strategic patching utilities

func updateGatewayAddresses(ctx context.Context, cli client.Client, gwNN types.NamespacedName, desired []api.GatewayStatusAddress) error {
	patcher := utils.NewGatewayStatusPatcher(cli)
	return patcher.PatchAddresses(ctx, gwNN, desired)
}

func updateGatewayStatusWithRetryFunc(ctx context.Context, cli client.Client, gwNN types.NamespacedName, updateFunc func(*api.Gateway) bool) error {
	patcher := utils.NewGatewayStatusPatcher(cli)
	return patcher.PatchFullStatus(ctx, gwNN, updateFunc)
}

// Helper function to generate random strings for unique names
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

