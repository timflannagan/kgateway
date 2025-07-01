package trafficpolicy

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/pkg/validator"
	"github.com/kgateway-dev/kgateway/v2/pkg/xds/bootstrap"
)

func (p *TrafficPolicy) Validate(ctx context.Context, v validator.Validator, policy *v1alpha1.TrafficPolicy) error {
	if err := p.validateProto(ctx); err != nil {
		return err
	}
	if err := p.validateXDS(ctx, v); err != nil {
		return err
	}
	return nil
}

func (p *TrafficPolicy) validateProto(ctx context.Context) error {
	// TODO: rustformations.
	// TODO: ext auth & rate limit provider validation
	var validators []func() error
	if p.spec.transform != nil {
		validators = append(validators, p.spec.transform.Validate)
	}
	if p.spec.localRateLimit != nil {
		validators = append(validators, p.spec.localRateLimit.Validate)
	}
	if p.spec.rateLimit != nil {
		for _, rateLimit := range p.spec.rateLimit.rateLimitActions {
			validators = append(validators, rateLimit.Validate)
		}
	}
	if p.spec.ExtProc != nil {
		if p.spec.ExtProc.ExtProcPerRoute != nil {
			validators = append(validators, p.spec.ExtProc.ExtProcPerRoute.Validate)
		}
	}
	if p.spec.extAuth != nil {
		if p.spec.extAuth.extauthPerRoute != nil {
			validators = append(validators, p.spec.extAuth.extauthPerRoute.Validate)
		}
	}
	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
}

// validateXDS builds a partial bootstrap config and validates it via
// envoy validate mode.
func (p *TrafficPolicy) validateXDS(ctx context.Context, v validator.Validator) error {
	builder := bootstrap.New()
	if p.spec.transform != nil {
		builder.AddFilterConfig(transformationFilterNamePrefix, p.spec.transform)
	}
	if p.spec.rustformation != nil {
		builder.AddFilterConfig(rustformationFilterNamePrefix, p.spec.rustformation)
	}
	if p.spec.localRateLimit != nil {
		builder.AddFilterConfig(localRateLimitFilterNamePrefix, p.spec.localRateLimit)
	}
	if p.spec.rateLimit != nil && p.spec.rateLimit.provider != nil {
		builder.AddFilterConfig(getRateLimitFilterName(p.spec.rateLimit.provider.ResourceName()), p.spec.rateLimit.provider.RateLimit)
	}
	if p.spec.ExtProc != nil && p.spec.ExtProc.provider != nil {
		builder.AddFilterConfig(extProcFilterName(p.spec.ExtProc.provider.ResourceName()), p.spec.ExtProc.provider.ExtProc)
	}
	if p.spec.extAuth != nil && p.spec.extAuth.provider != nil {
		builder.AddFilterConfig(extAuthFilterName(p.spec.extAuth.provider.ResourceName()), p.spec.extAuth.provider.ExtAuth)
	}

	bootstrap, err := builder.Build()
	if err != nil {
		return err
	}
	data, err := protojson.Marshal(bootstrap)
	if err != nil {
		return err
	}
	if err := v.Validate(ctx, string(data)); err != nil {
		return err
	}
	return nil
}
