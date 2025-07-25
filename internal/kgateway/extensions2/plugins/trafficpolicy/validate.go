package trafficpolicy

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/settings"
	"github.com/kgateway-dev/kgateway/v2/pkg/validator"
	"github.com/kgateway-dev/kgateway/v2/pkg/xds/bootstrap"
)

func (p *TrafficPolicy) Validate(ctx context.Context, v validator.Validator, mode settings.RouteReplacementMode) error {
	switch mode {
	case settings.RouteReplacementStandard:
		return p.validateStandard()
	case settings.RouteReplacementStrict:
		return p.validateStrict(ctx, v)
	}
	return nil
}

// validateStandard performs basic proto validation that runs in STANDARD mode
func (p *TrafficPolicy) validateStandard() error {
	return p.validateProto()
}

// validateStrict performs both proto and xDS validation that runs in STRICT mode
func (p *TrafficPolicy) validateStrict(ctx context.Context, v validator.Validator) error {
	if err := p.validateStandard(); err != nil {
		return err
	}
	return p.validateXDS(ctx, v)
}

// validateProto performs basic proto validation that runs in STANDARD mode.
// Each policy sub-IR maintains its own Validate() method following the sub-IR pattern.
func (p *TrafficPolicy) validateProto() error {
	var validators []func() error
	// Collect validation functions from each policy sub-IR
	if p.spec.ai != nil {
		validators = append(validators, p.spec.ai.Validate)
	}
	if p.spec.transformation != nil {
		validators = append(validators, p.spec.transformation.Validate)
	}
	if p.spec.rustformation != nil {
		validators = append(validators, p.spec.rustformation.Validate)
	}
	if p.spec.localRateLimit != nil {
		validators = append(validators, p.spec.localRateLimit.Validate)
	}
	if p.spec.rateLimit != nil {
		validators = append(validators, p.spec.rateLimit.Validate)
	}
	if p.spec.extProc != nil {
		validators = append(validators, p.spec.extProc.Validate)
	}
	if p.spec.extAuth != nil {
		validators = append(validators, p.spec.extAuth.Validate)
	}
	if p.spec.csrf != nil {
		validators = append(validators, p.spec.csrf.Validate)
	}
	if p.spec.cors != nil {
		validators = append(validators, p.spec.cors.Validate)
	}
	if p.spec.buffer != nil {
		validators = append(validators, p.spec.buffer.Validate)
	}
	if p.spec.hashPolicies != nil {
		validators = append(validators, p.spec.hashPolicies.Validate)
	}
	if p.spec.autoHostRewrite != nil {
		validators = append(validators, p.spec.autoHostRewrite.Validate)
	}
	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
}

// validateXDS builds a partial bootstrap config and validates it via envoy
// validate mode. It re-uses the ApplyForRoute method to ensure that the translation
// and validation logic go through the same code path as normal.
func (p *TrafficPolicy) validateXDS(ctx context.Context, v validator.Validator) error {
	// use a fake translation pass to ensure we have the desired typed filter config
	// on the placeholder vhost.
	typedPerFilterConfig := ir.TypedFilterConfigMap(map[string]proto.Message{})
	fakePass := NewGatewayTranslationPass(ctx, ir.GwTranslationCtx{}, nil)
	if err := fakePass.ApplyForRoute(ctx, &ir.RouteContext{
		Policy:            p,
		TypedFilterConfig: typedPerFilterConfig,
	}, nil); err != nil {
		return err
	}

	// build a partial bootstrap config with the typed filter config applied.
	builder := bootstrap.New()
	for name, config := range typedPerFilterConfig {
		builder.AddFilterConfig(name, config)
	}
	bootstrap, err := builder.Build()
	if err != nil {
		return err
	}
	data, err := protojson.Marshal(bootstrap)
	if err != nil {
		return err
	}

	// shell out to envoy to validate the partial bootstrap config.
	return v.Validate(ctx, string(data))
}
