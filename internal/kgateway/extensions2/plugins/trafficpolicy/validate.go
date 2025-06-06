package trafficpolicy

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/validator"
	"github.com/kgateway-dev/kgateway/v2/pkg/xds/bootstrap"
)

func (p *TrafficPolicy) Validate(ctx context.Context, v validator.Validator, policy *v1alpha1.TrafficPolicy) error {
	if shouldSkipValidation(policy) {
		logger.Info("skipping validation for policy", "policy", policy.Name)
		return nil
	}
	if err := p.validateProto(ctx); err != nil {
		return err
	}
	if err := p.validateXDS(ctx, v); err != nil {
		return err
	}
	return nil
}

// TODO: this is a bit of a mess.
func (p *TrafficPolicy) validateProto(ctx context.Context) error {
	var validators []func() error
	if p.spec.transform != nil {
		validators = append(validators, p.spec.transform.Validate)
	}
	// TODO: rustformations?
	if p.spec.localRateLimit != nil {
		validators = append(validators, p.spec.localRateLimit.Validate)
	}
	if p.spec.rateLimit != nil {
		if p.spec.rateLimit.provider != nil {
			validators = append(validators, p.spec.rateLimit.provider.Validate)
		}
		for _, rateLimit := range p.spec.rateLimit.rateLimitActions {
			validators = append(validators, rateLimit.Validate)
		}
	}
	if p.spec.ExtProc != nil {
		if p.spec.ExtProc.ExtProcPerRoute != nil {
			validators = append(validators, p.spec.ExtProc.ExtProcPerRoute.Validate)
		}
		if p.spec.ExtProc.provider != nil {
			validators = append(validators, p.spec.ExtProc.provider.Validate)
		}
	}
	if p.spec.extAuth != nil {
		if p.spec.extAuth.extauthPerRoute != nil {
			validators = append(validators, p.spec.extAuth.extauthPerRoute.Validate)
		}
		if p.spec.extAuth.provider != nil {
			validators = append(validators, p.spec.extAuth.provider.Validate)
		}
	}
	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
}

// validateXDS validates the xDS configuration.
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

func shouldSkipValidation(policy *v1alpha1.TrafficPolicy) bool {
	// TODO(tim): verify whether this is the right approach. less familiar with ancestors and
	// the implications of this approach if a policy attaches to multiple Gateways
	// TODO(tim): verify whether hardcoding the wellknown controller name is a safe assumption.
	for _, ancestor := range policy.Status.Ancestors {
		// not our controller, skip
		if ancestor.ControllerName != wellknown.GatewayControllerName {
			continue
		}
		// check for the Invalid condition reason at the current generation to
		// determine if we need to skip validation.
		cond := meta.FindStatusCondition(ancestor.Conditions, string(gwv1alpha2.PolicyConditionAccepted))
		if cond == nil {
			continue
		}
		if cond.Status != metav1.ConditionFalse {
			continue
		}
		if cond.Reason != string(gwv1alpha2.PolicyReasonInvalid) {
			continue
		}
		if cond.ObservedGeneration != policy.Generation {
			continue
		}
		return true
	}
	return false
}
