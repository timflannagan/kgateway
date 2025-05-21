package trafficpolicy

import (
	"fmt"

	envoy_ext_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/pluginutils"
)

// toEnvoyExtProc converts an ExtProcPolicy to an ExternalProcessor
func (b *TrafficPolicyBuilder) toEnvoyExtProc(
	krtctx krt.HandlerContext,
	trafficPolicy *v1alpha1.TrafficPolicy,
) (*ExtprocIR, error) {
	spec := trafficPolicy.Spec.ExtProc
	gatewayExtension, err := b.FetchGatewayExtension(krtctx, spec.ExtensionRef, trafficPolicy.GetNamespace())
	if err != nil {
		return nil, fmt.Errorf("extproc: %w", err)
	}
	if gatewayExtension.ExtType != v1alpha1.GatewayExtensionTypeExtProc || gatewayExtension.ExtProc == nil {
		return nil, pluginutils.ErrInvalidExtensionType(v1alpha1.GatewayExtensionTypeExtAuth, gatewayExtension.ExtType)
	}

	return &ExtprocIR{
		provider:        gatewayExtension,
		ExtProcPerRoute: translateExtProcPerFilterConfig(spec),
	}, nil
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
