package trafficpolicy

import (
	"context"

	ratev3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
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
	// TODO: rustformations, and ext auth/rate limit provider validation
	// Note: no need for buffer validation as it's a single int field, right?
	var validators []func() error
	if p.spec.AI != nil {
		if p.spec.AI.Transformation != nil {
			validators = append(validators, p.spec.AI.Transformation.Validate)
		}
		if p.spec.AI.Extproc != nil {
			validators = append(validators, p.spec.AI.Extproc.Validate)
		}
	}
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
	if p.spec.csrf != nil {
		validators = append(validators, p.spec.csrf.csrfPolicy.Validate)
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
	if p.spec.rateLimit != nil {
		// TODO: provider-based validation.
		if len(p.spec.rateLimit.rateLimitActions) > 0 {
			rateLimitPerRoute := &ratev3.RateLimitPerRoute{
				RateLimits: p.spec.rateLimit.rateLimitActions,
			}
			builder.AddFilterConfig(getRateLimitFilterName(p.spec.rateLimit.provider.ResourceName()), rateLimitPerRoute)
		}
	}
	if p.spec.ExtProc != nil {
		// TODO: provider-based validation.
		if p.spec.ExtProc.ExtProcPerRoute != nil {
			builder.AddFilterConfig(extProcFilterName(p.spec.ExtProc.provider.ResourceName()), p.spec.ExtProc.ExtProcPerRoute)
		}
	}
	if p.spec.extAuth != nil && p.spec.extAuth.provider != nil {
		// TODO: provider-based validation.
		if p.spec.extAuth.extauthPerRoute != nil {
			builder.AddFilterConfig(extAuthFilterName(p.spec.extAuth.provider.ResourceName()), p.spec.extAuth.extauthPerRoute)
		}
	}
	if p.spec.csrf != nil {
		builder.AddFilterConfig(csrfExtensionFilterName, p.spec.csrf.csrfPolicy)
	}
	if p.spec.AI != nil {
		if p.spec.AI.Transformation != nil {
			builder.AddFilterConfig(wellknown.AIPolicyTransformationFilterName, p.spec.AI.Transformation)
		}
		if p.spec.AI.Extproc != nil {
			builder.AddFilterConfig(wellknown.AIExtProcFilterName, p.spec.AI.Extproc)
		}
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
