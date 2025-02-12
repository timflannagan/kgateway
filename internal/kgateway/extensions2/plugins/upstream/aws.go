package upstream

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_lambda_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/aws_lambda/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/pluginutils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
)

const (
	ParameterGroup = "gateway.kgateway.dev"
	ParameterKind  = "Parameter"

	AccessKey    = "accessKey"
	SessionToken = "sessionToken"
	SecretKey    = "secretKey"
	FilterName   = "envoy.filters.http.aws_lambda"
)

func processAws(_ context.Context, in *v1alpha1.AwsUpstream, ir *UpstreamIr, out *envoy_config_cluster_v3.Cluster) error {
	if ir.AwsSecret == nil {
		return fmt.Errorf("AwsSecret is nil")
	}

	// Configure cluster for Lambda
	out.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
		Type: envoy_config_cluster_v3.Cluster_LOGICAL_DNS,
	}
	out.DnsLookupFamily = envoy_config_cluster_v3.Cluster_V4_ONLY

	lambdaHostname := getLambdaHostname(in)
	pluginutils.EnvoySingleEndpointLoadAssignment(out, lambdaHostname, 443)

	typedConfig, err := anypb.New(&envoyauth.UpstreamTlsContext{
		Sni: lambdaHostname,
	})
	if err != nil {
		return fmt.Errorf("failed to create tls context: %w", err)
	}
	out.TransportSocket = &envoy_config_core_v3.TransportSocket{
		Name: wellknown.TransportSocketTls,
		ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{
			TypedConfig: typedConfig,
		},
	}
	out.Metadata = &envoy_core_v3.Metadata{
		FilterMetadata: map[string]*structpb.Struct{
			"com.amazonaws.lambda": {
				Fields: map[string]*structpb.Value{
					// For whatever reason, this was needed. Saw in the trace-level logs that the filter was skipped otherwise.
					"egress_gateway": {Kind: &structpb.Value_BoolValue{
						BoolValue: true,
					}},
				},
			},
		},
	}
	return nil
}

func getLambdaHostname(in *v1alpha1.AwsUpstream) string {
	return fmt.Sprintf("lambda.%s.amazonaws.com", in.Region)
}

func processEndpointsAws(_ *v1alpha1.AwsUpstream) *ir.EndpointsForUpstream {
	return nil
}

func (p *plugin2) processBackendAws(
	_ context.Context,
	pCtx *ir.RouteBackendContext,
	dest *upstreamDestination,
) error {
	if p.needFilter == nil {
		p.needFilter = make(map[string]bool)
	}
	p.needFilter[pCtx.FilterChainName] = true

	if dest.FunctionName == "" {
		return fmt.Errorf("internal error: function name is not set")
	}

	lambdaRouteFunc := &envoy_lambda_v3.PerRouteConfig{
		InvokeConfig: &envoy_lambda_v3.Config{
			Arn: dest.FunctionName,
		},
	}
	lambdaRouteFuncAny, err := anypb.New(lambdaRouteFunc)
	if err != nil {
		return err
	}
	pCtx.AddTypedConfig(FilterName, lambdaRouteFuncAny)

	return nil
}

type staticSecretDerivation struct {
	access, session, secret string
}

// deriveStaticSecret from ingest if we are using a kubernetes secretref
// Named returns with the derived string contents or an error due to retrieval or format.
func deriveStaticSecret(awsSecrets *ir.Secret) (staticSecretDerivation, error) {
	derived := staticSecretDerivation{
		access:  string(awsSecrets.Data[AccessKey]),
		session: string(awsSecrets.Data[SessionToken]),
		secret:  string(awsSecrets.Data[SecretKey]),
	}

	var errs []error
	if derived.access == "" || !utf8.Valid([]byte(derived.access)) {
		errs = append(errs, errors.New("access_key is not a valid string"))
	}
	if derived.secret == "" || !utf8.Valid([]byte(derived.secret)) {
		errs = append(errs, errors.New("secret_key is not a valid string"))
	}
	if derived.session != "" && !utf8.Valid([]byte(derived.session)) {
		errs = append(errs, errors.New("session_key is not a valid string"))
	}
	return derived, errors.Join(errs...)
}
