package backendconfigpolicy

import (
	"context"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/validator"
	"github.com/kgateway-dev/kgateway/v2/pkg/xds/bootstrap"
)

func (p *BackendConfigPolicyIR) Validate(
	ctx context.Context,
	v validator.Validator,
	_ *v1alpha1.BackendConfigPolicy,
) error {
	if err := p.validateProto(ctx); err != nil {
		return err
	}
	if err := p.validateXDS(ctx, v); err != nil {
		return err
	}
	return nil
}

func (p *BackendConfigPolicyIR) validateProto(ctx context.Context) error {
	var validators []func() error
	if p.tcpKeepalive != nil {
		validators = append(validators, p.tcpKeepalive.Validate)
	}
	if p.commonHttpProtocolOptions != nil {
		validators = append(validators, p.commonHttpProtocolOptions.Validate)
	}
	if p.http1ProtocolOptions != nil {
		validators = append(validators, p.http1ProtocolOptions.Validate)
	}
	if p.tlsConfig != nil {
		validators = append(validators, p.tlsConfig.Validate)
	}
	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
}

func (p *BackendConfigPolicyIR) validateXDS(ctx context.Context, v validator.Validator) error {
	testCluster := buildValidationCluster(ctx, "validation-cluster", p)

	b := bootstrap.New()
	b.AddClusterConfig(testCluster)
	b.AddRouteConfig(&routev3.Route{
		Name: "validation-route",
		Match: &routev3.RouteMatch{
			PathSpecifier: &routev3.RouteMatch_Prefix{Prefix: "/"},
		},
		Action: &routev3.Route_Route{
			Route: &routev3.RouteAction{
				ClusterSpecifier: &routev3.RouteAction_Cluster{
					Cluster: testCluster.GetName(),
				},
			},
		},
	})
	bootstrap, err := b.Build()
	if err != nil {
		return err
	}
	raw, err := protojson.Marshal(bootstrap)
	if err != nil {
		return err
	}

	return v.Validate(ctx, string(raw))
}

func buildValidationCluster(ctx context.Context, name string, p *BackendConfigPolicyIR) *clusterv3.Cluster {
	cluster := &clusterv3.Cluster{
		Name: name,
	}
	var empty ir.BackendObjectIR
	processBackend(ctx, p, empty, cluster)
	return cluster
}
