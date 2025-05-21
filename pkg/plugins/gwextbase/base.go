package gwextbase

import (
	"context"

	"istio.io/istio/pkg/kube/krt"

	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugins/trafficpolicy"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
)

type TrafficPolicy = trafficpolicy.TrafficPolicy
type TrafficPolicyBuilder = trafficpolicy.TrafficPolicyBuilder

// Create a traffic policy builder. This converts a traffic policy into its IR form.
func NewTrafficPolicyBuilder(
	ctx context.Context,
	commoncol *common.CommonCollections,
	fetch func(krtctx krt.HandlerContext, extType v1alpha1.GatewayExtensionType) *ir.GatewayExtension,
) *trafficpolicy.TrafficPolicyBuilder {
	return trafficpolicy.NewTrafficPolicyBuilder(ctx, commoncol)
}

func NewGatewayTranslationPass(ctx context.Context, tctx ir.GwTranslationCtx, reporter reports.Reporter) ir.ProxyTranslationPass {
	return trafficpolicy.NewGatewayTranslationPass(ctx, tctx, reporter)
}

// ResolveExtGrpcService resolves a gateway extension gRPC service by looking up the backend reference
// and constructing an Envoy gRPC service configuration. It takes the following parameters:
//   - krtctx: The KRT context
//   - backends: The backend index collection
//   - disableExtensionRefValidation: Whether to skip reference grant validation
//   - objectSource: The source object making the request
//   - grpcService: The gRPC service configuration to resolve
//
// Returns:
//   - *envoy_core_v3.GrpcService: The resolved Envoy gRPC service configuration
//   - error: Any error that occurred during resolution

func ResolveExtGrpcService(krtctx krt.HandlerContext, backends *krtcollections.BackendIndex, disableExtensionRefValidation bool, objectSource ir.ObjectSource, grpcService *v1alpha1.ExtGrpcService) (*envoy_core_v3.GrpcService, error) {
	return trafficpolicy.ResolveExtGrpcService(krtctx, backends, disableExtensionRefValidation, objectSource, grpcService)
}
