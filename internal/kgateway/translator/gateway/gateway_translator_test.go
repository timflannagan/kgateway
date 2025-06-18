package gateway_test

import (
	"context"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/reporter"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/fsutils"
	translatortest "github.com/kgateway-dev/kgateway/v2/test/translator"
)

type translatorTestCase struct {
	inputFile     string
	outputFile    string
	gwNN          types.NamespacedName
	assertReports translatortest.AssertReports
}

var _ = DescribeTable("Basic GatewayTranslator Tests",
	func(in translatorTestCase) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		dir := fsutils.MustGetThisDir()

		inputFiles := []string{filepath.Join(dir, "testutils/inputs/", in.inputFile)}
		expectedProxyFile := filepath.Join(dir, "testutils/outputs/", in.outputFile)
		translatortest.TestTranslation(GinkgoT(), ctx, inputFiles, expectedProxyFile, in.gwNN, in.assertReports)
	},
	Entry(
		"http gateway with basic routing",
		translatorTestCase{
			inputFile:  "http-routing",
			outputFile: "http-routing-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
		}),
	Entry(
		"http gateway with custom class",
		translatorTestCase{
			inputFile:  "custom-gateway-class",
			outputFile: "custom-gateway-class.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
		}),
	Entry(
		"https gateway with basic routing",
		translatorTestCase{
			inputFile:  "https-routing",
			outputFile: "https-routing-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
		}),
	Entry(
		"http gateway with multiple listeners on the same port",
		translatorTestCase{
			inputFile:  "multiple-listeners-http-routing",
			outputFile: "multiple-listeners-http-routing-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "http",
			},
		}),
	Entry(
		"https gateway with multiple listeners on the same port",
		translatorTestCase{
			inputFile:  "multiple-listeners-https-routing",
			outputFile: "multiple-listeners-https-routing-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "http",
			},
		}),
	Entry(
		"http gateway with multiple routing rules and HeaderModifier filter",
		translatorTestCase{
			inputFile:  "http-with-header-modifier",
			outputFile: "http-with-header-modifier-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "gw",
			},
		}),
	XEntry(
		"http gateway with azure destination",
		translatorTestCase{
			inputFile:  "http-with-azure-destination",
			outputFile: "http-with-azure-destination-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "gw",
			},
		}),
	Entry(
		"gateway with correctly sorted routes",
		translatorTestCase{
			inputFile:  "route-sort.yaml",
			outputFile: "route-sort.yaml",
			gwNN: types.NamespacedName{
				Namespace: "infra",
				Name:      "example-gateway",
			},
		}),
	Entry(
		"httproute with missing backend reports correctly",
		translatorTestCase{
			inputFile:  "http-routing-missing-backend",
			outputFile: "http-routing-missing-backend.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Reason).To(Equal(string(gwv1.RouteReasonBackendNotFound)))
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
				Expect(resolvedRefs.Message).To(Equal(`Service "example-svc" not found`))
				Expect(resolvedRefs.ObservedGeneration).To(Equal(int64(0)))
			},
		}),
	Entry(
		"httproute with invalid backend reports correctly",
		translatorTestCase{
			inputFile:  "http-routing-invalid-backend",
			outputFile: "http-routing-invalid-backend.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Reason).To(Equal(string(gwv1.RouteReasonInvalidKind)))
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
				Expect(resolvedRefs.Message).To(Equal(`unknown backend kind`))
				Expect(resolvedRefs.ObservedGeneration).To(Equal(int64(0)))
			},
		}),
	Entry(
		"httproute with invalid prefix match reports correctly",
		translatorTestCase{
			inputFile:  "http-routing-invalid-prefix/manifest.yaml",
			outputFile: "http-routing-invalid-prefix/manifest.yaml",
			gwNN: types.NamespacedName{
				Namespace: "gwtest",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "invalid-traffic-policy-route",
						Namespace:  "gwtest",
						Generation: 1,
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				partiallyInvalid := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionPartiallyInvalid))
				Expect(partiallyInvalid).NotTo(BeNil())
				Expect(partiallyInvalid.Status).To(Equal(metav1.ConditionTrue))
				Expect(partiallyInvalid.Reason).To(Equal(string(gwv1.RouteReasonUnsupportedValue)))
				Expect(partiallyInvalid.Message).To(ContainSubstring(`Dropped Rule`))
				Expect(partiallyInvalid.ObservedGeneration).To(Equal(int64(1)))
			},
		}),
	Entry(
		"TrafficPolicy merging",
		translatorTestCase{
			inputFile:  "traffic-policy/merge.yaml",
			outputFile: "traffic-policy/merge.yaml",
			gwNN: types.NamespacedName{
				Namespace: "infra",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				expectedPolicies := []reports.PolicyKey{
					{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "infra", Name: "policy-with-section-name"},
					{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "infra", Name: "policy-without-section-name"},
				}
				assertAcceptedPolicyStatus(reportsMap, expectedPolicies)
			},
		}),
	Entry(
		"TrafficPolicy with targetSelectors",
		translatorTestCase{
			inputFile:  "traffic-policy/label_based.yaml",
			outputFile: "traffic-policy/label_based.yaml",
			gwNN: types.NamespacedName{
				Namespace: "infra",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				expectedPolicies := []reports.PolicyKey{
					{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "infra", Name: "transform"},
					{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "infra", Name: "rate-limit"},
				}
				assertAcceptedPolicyStatus(reportsMap, expectedPolicies)
			},
		}),
	Entry(
		"TrafficPolicy edge cases",
		translatorTestCase{
			inputFile:  "traffic-policy/extauth.yaml",
			outputFile: "traffic-policy/extauth.yaml",
			gwNN: types.NamespacedName{
				Namespace: "infra",
				Name:      "example-gateway",
			},
		}),
	Entry(
		"tcp gateway with basic routing",
		translatorTestCase{
			inputFile:  "tcp-routing/basic.yaml",
			outputFile: "tcp-routing/basic-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1alpha2.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-tcp-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionTrue))
				Expect(resolvedRefs.Reason).To(Equal(string(gwv1.RouteReasonResolvedRefs)))
			},
		}),
	Entry(
		"tcproute with missing backend reports correctly",
		translatorTestCase{
			inputFile:  "tcp-routing/missing-backend.yaml",
			outputFile: "tcp-routing/missing-backend.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1alpha2.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-tcp-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
				Expect(resolvedRefs.Message).To(Equal("Service \"example-tcp-svc\" not found"))
			},
		}),
	Entry(
		"tcproute with invalid backend reports correctly",
		translatorTestCase{
			inputFile:  "tcp-routing/invalid-backend.yaml",
			outputFile: "tcp-routing/invalid-backend.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1alpha2.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-tcp-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
				Expect(resolvedRefs.Message).To(Equal("unknown backend kind"))
			},
		}),
	Entry(
		"tcp gateway with multiple backend services",
		translatorTestCase{
			inputFile:  "tcp-routing/multi-backend.yaml",
			outputFile: "tcp-routing/multi-backend-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-tcp-gateway",
			},
		}),
	Entry(
		"tls gateway with basic routing",
		translatorTestCase{
			inputFile:  "tls-routing/basic.yaml",
			outputFile: "tls-routing/basic-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1alpha2.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-tls-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionTrue))
				Expect(resolvedRefs.Reason).To(Equal(string(gwv1.RouteReasonResolvedRefs)))
			},
		}),
	Entry(
		"tlsroute with missing backend reports correctly",
		translatorTestCase{
			inputFile:  "tls-routing/missing-backend.yaml",
			outputFile: "tls-routing/missing-backend.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1alpha2.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-tls-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
				Expect(resolvedRefs.Message).To(Equal("Service \"example-tls-svc\" not found"))
			},
		}),
	Entry(
		"tlsroute with invalid backend reports correctly",
		translatorTestCase{
			inputFile:  "tls-routing/invalid-backend.yaml",
			outputFile: "tls-routing/invalid-backend.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1alpha2.TLSRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-tls-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
				Expect(resolvedRefs.Message).To(Equal("unknown backend kind"))
			},
		}),
	Entry(
		"tls gateway with multiple backend services",
		translatorTestCase{
			inputFile:  "tls-routing/multi-backend.yaml",
			outputFile: "tls-routing/multi-backend-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
		}),
	Entry(
		"grpc gateway with basic routing",
		translatorTestCase{
			inputFile:  "grpc-routing/basic.yaml",
			outputFile: "grpc-routing/basic-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1.GRPCRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-grpc-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionTrue))
				Expect(resolvedRefs.Reason).To(Equal(string(gwv1.RouteReasonResolvedRefs)))
			},
		}),
	Entry(
		"grpcroute with missing backend reports correctly",
		translatorTestCase{
			inputFile:  "grpc-routing/missing-backend.yaml",
			outputFile: "grpc-routing/missing-backend.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1.GRPCRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-grpc-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
				Expect(resolvedRefs.Message).To(Equal(`Service "example-grpc-svc" not found`))
			},
		}),
	Entry(
		"grpcroute with invalid backend reports correctly",
		translatorTestCase{
			inputFile:  "grpc-routing/invalid-backend.yaml",
			outputFile: "grpc-routing/invalid-backend.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				route := &gwv1.GRPCRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example-grpc-route",
						Namespace: "default",
					},
				}
				routeStatus := reportsMap.BuildRouteStatus(context.Background(), route, wellknown.GatewayControllerName)
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				resolvedRefs := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs).NotTo(BeNil())
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
				Expect(resolvedRefs.Message).To(Equal("unknown backend kind"))
			},
		}),
	Entry(
		"grpc gateway with multiple backend services",
		translatorTestCase{
			inputFile:  "grpc-routing/multi-backend.yaml",
			outputFile: "grpc-routing/multi-backend-proxy.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-grpc-gateway",
			},
		}),
	Entry("Basic service backend", translatorTestCase{
		inputFile:  "backends/basic.yaml",
		outputFile: "backends/basic.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("AWS Lambda backend", translatorTestCase{
		inputFile:  "backends/aws_lambda.yaml",
		outputFile: "backends/aws_lambda.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("DFP Backend with TLS", translatorTestCase{
		inputFile:  "dfp/tls.yaml",
		outputFile: "dfp/tls.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("DFP Backend with simple", translatorTestCase{
		inputFile:  "dfp/simple.yaml",
		outputFile: "dfp/simple.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Backend TLS Policy", translatorTestCase{
		inputFile:  "backendtlspolicy/tls.yaml",
		outputFile: "backendtlspolicy/tls.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Backend TLS Policy with SAN", translatorTestCase{
		inputFile:  "backendtlspolicy/tls-san.yaml",
		outputFile: "backendtlspolicy/tls-san.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Proxy with no routes", translatorTestCase{
		inputFile:  "edge-cases/no_route.yaml",
		outputFile: "no_route.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Direct response", translatorTestCase{
		inputFile:  "directresponse/manifest.yaml",
		outputFile: "directresponse.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("HTTPRoutes with timeout and retry", translatorTestCase{
		inputFile:  "httproute-timeout-retry/manifest.yaml",
		outputFile: "httproute-timeout-retry-proxy.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry(
		"http gateway with session persistence (cookie)",
		translatorTestCase{
			inputFile:  "session-persistence/cookie.yaml",
			outputFile: "session-persistence/cookie.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
		}),
	Entry(
		"http gateway with session persistence (header)",
		translatorTestCase{
			inputFile:  "session-persistence/header.yaml",
			outputFile: "session-persistence/header.yaml",
			gwNN: types.NamespacedName{
				Namespace: "default",
				Name:      "example-gateway",
			},
		}),
	Entry("HTTPListenerPolicy with upgrades", translatorTestCase{
		inputFile:  "https-listener-pol/upgrades.yaml",
		outputFile: "https-listener-pol/upgrades.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Service with appProtocol=kubernetes.io/h2c", translatorTestCase{
		inputFile:  "backend-protocol/svc-h2c.yaml",
		outputFile: "backend-protocol/svc-h2c.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Service with appProtocol=kubernetes.io/ws", translatorTestCase{
		inputFile:  "backend-protocol/svc-ws.yaml",
		outputFile: "backend-protocol/svc-ws.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Service with appProtocol=anything", translatorTestCase{
		inputFile:  "backend-protocol/svc-default.yaml",
		outputFile: "backend-protocol/svc-default.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Static Backend with appProtocol=kubernetes.io/h2c", translatorTestCase{
		inputFile:  "backend-protocol/backend-h2c.yaml",
		outputFile: "backend-protocol/backend-h2c.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Static Backend with appProtocol=kubernetes.io/ws", translatorTestCase{
		inputFile:  "backend-protocol/backend-ws.yaml",
		outputFile: "backend-protocol/backend-ws.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Static Backend with no appProtocol", translatorTestCase{
		inputFile:  "backend-protocol/backend-default.yaml",
		outputFile: "backend-protocol/backend-default.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Backend Config Policy with LB Config", translatorTestCase{
		inputFile:  "backendconfigpolicy/lb-config.yaml",
		outputFile: "backendconfigpolicy/lb-config.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Backend Config Policy with Common HTTP Protocol - HTTP backend", translatorTestCase{
		inputFile:  "backendconfigpolicy/commonhttpprotocol-httpbackend.yaml",
		outputFile: "backendconfigpolicy/commonhttpprotocol-httpbackend.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry("Backend Config Policy with Common HTTP Protocol - HTTP2 backend", translatorTestCase{
		inputFile:  "backendconfigpolicy/commonhttpprotocol-http2backend.yaml",
		outputFile: "backendconfigpolicy/commonhttpprotocol-http2backend.yaml",
		gwNN: types.NamespacedName{
			Namespace: "default",
			Name:      "example-gateway",
		},
	}),
	Entry(
		"TrafficPolicy with explicit generation",
		translatorTestCase{
			inputFile:  "traffic-policy/generation.yaml",
			outputFile: "traffic-policy/generation.yaml",
			gwNN: types.NamespacedName{
				Namespace: "infra",
				Name:      "example-gateway",
			},
			assertReports: func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				expectedPolicies := []reports.PolicyKey{
					{Group: "gateway.kgateway.dev", Kind: "TrafficPolicy", Namespace: "infra", Name: "test-policy"},
				}
				assertPolicyStatusWithGeneration(reportsMap, expectedPolicies, 42)
			},
		}),
	// TODO: Add this once istio adds support for listener sets
	// Entry(
	// 	"listener sets",
	// 	translatorTestCase{
	// 		inputFile:  "listener-sets/manifest.yaml",
	// 		outputFile: "listener-sets-proxy.yaml",
	// 		gwNN: types.NamespacedName{
	// 			Namespace: "default",
	// 			Name:      "example-gateway",
	// 		},
	// 	}),
)

var _ = DescribeTable("Route Delegation translator",
	func(inputFile string, errdesc string) {
		dir := fsutils.MustGetThisDir()
		translatortest.TestTranslation(
			GinkgoT(),
			context.Background(),
			[]string{
				filepath.Join(dir, "testutils/inputs/delegation/common.yaml"),
				filepath.Join(dir, "testutils/inputs/delegation", inputFile),
			},
			filepath.Join(dir, "testutils/outputs/delegation", inputFile),
			types.NamespacedName{
				Namespace: "infra",
				Name:      "example-gateway",
			},
			func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				if errdesc == "" {
					Expect(translatortest.AreReportsSuccess(gwNN, reportsMap)).NotTo(HaveOccurred())
				} else {
					Expect(translatortest.AreReportsSuccess(gwNN, reportsMap)).To(MatchError(ContainSubstring(errdesc)))
				}
			},
		)
	},
	Entry("Basic config", "basic.yaml", ""),
	Entry("Child matches parent via parentRefs", "basic_parentref_match.yaml", ""),
	Entry("Child doesn't match parent via parentRefs", "basic_parentref_mismatch.yaml", "BackendNotFound unresolved reference gateway.networking.k8s.io/HTTPRoute/a/*"),
	Entry("Children using parentRefs and inherit-parent-matcher", "inherit_parentref.yaml", ""),
	Entry("Parent delegates to multiple chidren", "multiple_children.yaml", ""),
	Entry("Child is invalid as it is delegatee and specifies hostnames", "basic_invalid_hostname.yaml", "spec.hostnames must be unset on a delegatee route as they are inherited from the parent route"),
	Entry("Multi-level recursive delegation", "recursive.yaml", ""),
	Entry("Cyclic child route", "cyclic1.yaml", "cyclic reference detected while evaluating delegated routes"),
	Entry("Multi-level cyclic child route", "cyclic2.yaml", "cyclic reference detected while evaluating delegated routes"),
	Entry("Child rule matcher", "child_rule_matcher.yaml", ""),
	Entry("Child with multiple parents", "multiple_parents.yaml", "BackendNotFound unresolved reference gateway.networking.k8s.io/HTTPRoute/b/*"),
	Entry("Child can be an invalid delegatee but valid standalone", "invalid_child_valid_standalone.yaml", "spec.hostnames must be unset on a delegatee route as they are inherited from the parent route"),
	Entry("Relative paths", "relative_paths.yaml", ""),
	Entry("Nested absolute and relative path inheritance", "nested_absolute_relative.yaml", ""),
	Entry("Child route matcher does not match parent", "discard_invalid_child_matches.yaml", ""),
	Entry("Multi-level multiple parents delegation", "multi_level_multiple_parents.yaml", ""),
	Entry("TrafficPolicy only on child", "traffic_policy.yaml", ""),
	Entry("TrafficPolicy inheritance from parent", "traffic_policy_inheritance.yaml", ""),
	Entry("TrafficPolicy ignore child override on conflict", "traffic_policy_inheritance_child_override_ignore.yaml", ""),
	Entry("TrafficPolicy merge child override on no conflict", "traffic_policy_inheritance_child_override_ok.yaml", ""),
	Entry("TrafficPolicy multi level inheritance with child override disabled", "traffic_policy_multi_level_inheritance_override_disabled.yaml", ""),
	Entry("TrafficPolicy multi level inheritance with child override enabled", "traffic_policy_multi_level_inheritance_override_enabled.yaml", ""),
	Entry("TrafficPolicy filter override merge", "traffic_policy_filter_override_merge.yaml", ""),
	Entry("Built-in rule inheritance", "builtin_rule_inheritance.yaml", ""),
	Entry("Label based delegation", "label_based.yaml", ""),
)

var _ = DescribeTable("Discovery Namespace Selector",
	func(cfgJSON string, inputFile string, outputFile string, errdesc string) {
		dir := fsutils.MustGetThisDir()
		translatortest.TestTranslation(
			GinkgoT(),
			context.Background(),
			[]string{
				filepath.Join(dir, "testutils/inputs/discovery-namespace-selector", inputFile),
			},
			filepath.Join(dir, "testutils/outputs/discovery-namespace-selector", outputFile),
			types.NamespacedName{
				Namespace: "infra",
				Name:      "example-gateway",
			},
			func(gwNN types.NamespacedName, reportsMap reports.ReportMap) {
				if errdesc == "" {
					Expect(translatortest.AreReportsSuccess(gwNN, reportsMap)).NotTo(HaveOccurred())
				} else {
					Expect(translatortest.AreReportsSuccess(gwNN, reportsMap)).To(MatchError(ContainSubstring(errdesc)))
				}
			},
			translatortest.SettingsWithDiscoveryNamespaceSelectors(cfgJSON),
		)
	},
	Entry("Select all resources",
		`[
  {
    "matchExpressions": [
      {
        "key": "kubernetes.io/metadata.name",
        "operator": "In",
        "values": [
          "infra"
        ]
      }
    ]
  },
	{
		"matchLabels": {
			"app": "a"
		}
	}
]`,
		"base.yaml", "base_select_all.yaml", ""),
	Entry("Select all resources; AND matchExpressions and matchLabels",
		`[
  {
    "matchExpressions": [
      {
        "key": "kubernetes.io/metadata.name",
        "operator": "In",
        "values": [
          "infra"
        ]
      }
    ]
  },
	{
    "matchExpressions": [
      {
        "key": "kubernetes.io/metadata.name",
        "operator": "In",
        "values": [
          "a"
        ]
      }
    ],
		"matchLabels": {
			"app": "a"
		}
	}
]`,
		"base.yaml", "base_select_all.yaml", ""),
	Entry("Select only namespace infra",
		`[
  {
    "matchExpressions": [
      {
        "key": "kubernetes.io/metadata.name",
        "operator": "In",
        "values": [
          "infra"
        ]
      }
    ]
  }
]`,
		"base.yaml", "base_select_infra.yaml", "condition error for httproute: infra/example-route"),
)

// assertPolicyStatusWithGeneration is a helper function to verify policy status conditions with a specific generation
func assertPolicyStatusWithGeneration(reportsMap reports.ReportMap, policies []reports.PolicyKey, expectedGeneration int64) {
	var currentStatus gwv1alpha2.PolicyStatus

	for _, policy := range policies {
		// Validate each policy's status
		status := reportsMap.BuildPolicyStatus(context.Background(), policy, wellknown.GatewayControllerName, currentStatus)
		Expect(status).NotTo(BeNil(), "status missing for policy %v", policy)
		Expect(status.Ancestors).To(HaveLen(1), "ancestor missing for policy %v", policy) // 1 Gateway(ancestor)

		acceptedCondition := meta.FindStatusCondition(status.Ancestors[0].Conditions, string(gwv1alpha2.PolicyConditionAccepted))
		Expect(acceptedCondition).NotTo(BeNil())
		Expect(acceptedCondition.Status).To(Equal(metav1.ConditionTrue))
		Expect(acceptedCondition.Reason).To(Equal(string(gwv1alpha2.PolicyReasonAccepted)))
		Expect(acceptedCondition.Message).To(Equal(reporter.PolicyAcceptedMsg))
		Expect(acceptedCondition.ObservedGeneration).To(Equal(expectedGeneration))
	}
}

// assertAcceptedPolicyStatus is a helper function to verify policy status conditions
func assertAcceptedPolicyStatus(reportsMap reports.ReportMap, policies []reports.PolicyKey) {
	assertPolicyStatusWithGeneration(reportsMap, policies, 0)
}
