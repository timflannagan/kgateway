package trafficpolicy

import (
	"fmt"

	envoy_ext_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/pluginutils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	pluginsdkir "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/cmputils"
)

type ExtprocIR struct {
	provider *TrafficPolicyGatewayExtensionIR
	perRoute *envoy_ext_proc_v3.ExtProcPerRoute
}

var _ PolicySubIR = &ExtprocIR{}

func (e *ExtprocIR) Equals(other PolicySubIR) bool {
	otherExtproc, ok := other.(*ExtprocIR)
	if !ok {
		return false
	}
	if e == nil && otherExtproc == nil {
		return true
	}
	if e == nil || otherExtproc == nil {
		return false
	}
	if !proto.Equal(e.perRoute, otherExtproc.perRoute) {
		return false
	}
	if !cmputils.CompareWithNils(e.provider, otherExtproc.provider, func(a, b *TrafficPolicyGatewayExtensionIR) bool {
		return a.Equals(*b)
	}) {
		return false
	}
	return true
}

// Validate performs PGV-based validation on the extproc components
func (e *ExtprocIR) Validate() error {
	if e == nil {
		return nil
	}
	if e.perRoute != nil {
		if err := e.perRoute.ValidateAll(); err != nil {
			return err
		}
	}
	if e.provider != nil {
		if err := e.provider.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// MergeInto handles merging extproc policy from p2 into p1
func (e *ExtprocIR) MergeInto(
	p1, p2 *TrafficPolicy,
	p2Ref *pluginsdkir.AttachedPolicyRef,
	opts policy.MergeOptions,
	mergeOrigins pluginsdkir.MergeOrigins,
) {
	if !policy.IsMergeable(p1.spec.extProc, p2.spec.extProc, opts) {
		return
	}

	switch opts.Strategy {
	case policy.AugmentedDeepMerge, policy.OverridableDeepMerge:
		if p1.spec.extProc != nil {
			return
		}
		fallthrough // can override p1 if it is unset

	case policy.AugmentedShallowMerge, policy.OverridableShallowMerge:
		p1.spec.extProc = p2.spec.extProc
		mergeOrigins.SetOne("extProc", p2Ref)

	default:
		logger.Warn("unsupported merge strategy for extProc policy", "strategy", opts.Strategy, "policy", p2Ref)
	}
}

// toEnvoyExtProc converts an ExtProcPolicy to an ExternalProcessor
func extProcForSpec(
	krtctx krt.HandlerContext,
	in *v1alpha1.TrafficPolicy,
	out *trafficPolicySpecIr,
	fetchGatewayExtension FetchGatewayExtensionFunc,
) error {
	if in.Spec.ExtProc == nil {
		return nil
	}
	gatewayExtension, err := fetchGatewayExtension(krtctx, in.Spec.ExtProc.ExtensionRef, in.GetNamespace())
	if err != nil {
		return fmt.Errorf("extproc: %w", err)
	}
	if gatewayExtension.ExtType != v1alpha1.GatewayExtensionTypeExtProc || gatewayExtension.ExtProc == nil {
		return pluginutils.ErrInvalidExtensionType(v1alpha1.GatewayExtensionTypeExtAuth, gatewayExtension.ExtType)
	}
	out.extProc = &ExtprocIR{
		provider: gatewayExtension,
		perRoute: translateExtProcPerFilterConfig(in.Spec.ExtProc),
	}
	return nil
}

func translateExtProcPerFilterConfig(extProc *v1alpha1.ExtProcPolicy) *envoy_ext_proc_v3.ExtProcPerRoute {
	overrides := &envoy_ext_proc_v3.ExtProcOverrides{}
	if extProc.ProcessingMode != nil {
		overrides.ProcessingMode = toEnvoyProcessingMode(extProc.ProcessingMode)
	}

	return &envoy_ext_proc_v3.ExtProcPerRoute{
		Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{
			Overrides: overrides,
		},
	}
}

// headerSendModeFromString converts a string to envoy HeaderSendMode
func headerSendModeFromString(mode *string) envoy_ext_proc_v3.ProcessingMode_HeaderSendMode {
	if mode == nil {
		return envoy_ext_proc_v3.ProcessingMode_DEFAULT
	}
	switch *mode {
	case "SEND":
		return envoy_ext_proc_v3.ProcessingMode_SEND
	case "SKIP":
		return envoy_ext_proc_v3.ProcessingMode_SKIP
	default:
		return envoy_ext_proc_v3.ProcessingMode_DEFAULT
	}
}

// bodySendModeFromString converts a string to envoy BodySendMode
func bodySendModeFromString(mode *string) envoy_ext_proc_v3.ProcessingMode_BodySendMode {
	if mode == nil {
		return envoy_ext_proc_v3.ProcessingMode_NONE
	}
	switch *mode {
	case "STREAMED":
		return envoy_ext_proc_v3.ProcessingMode_STREAMED
	case "BUFFERED":
		return envoy_ext_proc_v3.ProcessingMode_BUFFERED
	case "BUFFERED_PARTIAL":
		return envoy_ext_proc_v3.ProcessingMode_BUFFERED_PARTIAL
	case "FULL_DUPLEX_STREAMED":
		return envoy_ext_proc_v3.ProcessingMode_FULL_DUPLEX_STREAMED
	default:
		return envoy_ext_proc_v3.ProcessingMode_NONE
	}
}

// toEnvoyProcessingMode converts our ProcessingMode to envoy's ProcessingMode
func toEnvoyProcessingMode(p *v1alpha1.ProcessingMode) *envoy_ext_proc_v3.ProcessingMode {
	if p == nil {
		return nil
	}

	return &envoy_ext_proc_v3.ProcessingMode{
		RequestHeaderMode:   headerSendModeFromString(p.RequestHeaderMode),
		ResponseHeaderMode:  headerSendModeFromString(p.ResponseHeaderMode),
		RequestBodyMode:     bodySendModeFromString(p.RequestBodyMode),
		ResponseBodyMode:    bodySendModeFromString(p.ResponseBodyMode),
		RequestTrailerMode:  headerSendModeFromString(p.RequestTrailerMode),
		ResponseTrailerMode: headerSendModeFromString(p.ResponseTrailerMode),
	}
}

// FIXME: Using the wrong filter name prefix when the name is empty?
func extProcFilterName(name string) string {
	if name == "" {
		return extauthFilterNamePrefix
	}
	return fmt.Sprintf("%s/%s", "ext_proc", name)
}

func (p *trafficPolicyPluginGwPass) handleExtProc(fcn string, pCtxTypedFilterConfig *ir.TypedFilterConfigMap, extProc *ExtprocIR) {
	if extProc == nil || extProc.provider == nil {
		return
	}
	providerName := extProc.provider.ResourceName()
	// Handle the enablement state

	if extProc.perRoute != nil {
		pCtxTypedFilterConfig.AddTypedConfig(extProcFilterName(providerName),
			extProc.perRoute,
		)
	} else {
		// if you are on a route and not trying to disable it then we need to override the top level disable on the filter chain
		pCtxTypedFilterConfig.AddTypedConfig(extProcFilterName(providerName),
			&envoy_ext_proc_v3.ExtProcPerRoute{Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{Overrides: &envoy_ext_proc_v3.ExtProcOverrides{}}},
		)
	}
	p.extProcPerProvider.Add(fcn, providerName, extProc.provider)
}
