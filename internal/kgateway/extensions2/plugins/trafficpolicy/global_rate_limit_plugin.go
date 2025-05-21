package trafficpolicy

import (
	"errors"
	"fmt"
	"time"

	routeconfv3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"google.golang.org/protobuf/proto"
	"istio.io/istio/pkg/kube/krt"

	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/pluginutils"
)

const (
	rateLimitFilterName     = "envoy.filters.http.ratelimit"
	rateLimitStatPrefix     = "http_rate_limit"
	defaultRateLimitTimeout = 2 * time.Second
)

// RateLimitIR represents the intermediate representation of a rate limit policy
type globalRateLimitIR struct {
	provider         *TrafficPolicyGatewayExtensionIR
	rateLimitActions []*routeconfv3.RateLimit
}

// Equals checks if two globalRateLimitIR objects are equal
func (r *globalRateLimitIR) Equals(other *globalRateLimitIR) bool {
	if r == nil && other == nil {
		return true
	}
	if r == nil || other == nil {
		return false
	}

	if len(r.rateLimitActions) != len(other.rateLimitActions) {
		return false
	}
	for i, action := range r.rateLimitActions {
		if !proto.Equal(action, other.rateLimitActions[i]) {
			return false
		}
	}
	if (r.provider == nil) != (other.provider == nil) {
		return false
	}
	if r.provider != nil && !r.provider.Equals(*other.provider) {
		return false
	}

	return true
}

func (r *globalRateLimitIR) Validate() error {
	// Implement me.
	return nil
}

// Add this function to handle the global rate limit configuration
func (b *TrafficPolicyBuilder) rateLimitForSpec(
	krtctx krt.HandlerContext,
	policy *v1alpha1.TrafficPolicy,
	out *trafficPolicySpecIr,
) []error {
	if policy.Spec.RateLimit == nil || policy.Spec.RateLimit.Global == nil {
		return nil
	}
	var errors []error
	globalPolicy := policy.Spec.RateLimit.Global

	// Create rate limit actions for the route or vhost
	actions, err := createRateLimitActions(globalPolicy.Descriptors)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to create rate limit actions: %w", err))
	}

	gwExtIR, err := b.FetchGatewayExtension(krtctx, globalPolicy.ExtensionRef, policy.GetNamespace())
	if err != nil {
		errors = append(errors, fmt.Errorf("ratelimit: %w", err))
		return errors
	}
	if gwExtIR.ExtType != v1alpha1.GatewayExtensionTypeRateLimit || gwExtIR.RateLimit == nil {
		errors = append(errors, pluginutils.ErrInvalidExtensionType(v1alpha1.GatewayExtensionTypeExtAuth, gwExtIR.ExtType))
	}

	if len(errors) > 0 {
		return errors
	}

	// Create route rate limits and store in the RateLimitIR struct
	out.globalRateLimit = &globalRateLimitIR{
		provider: gwExtIR,
		rateLimitActions: []*routev3.RateLimit{
			{
				Actions: actions,
			},
		},
	}
	return nil
}

// createRateLimitActions translates the API descriptors to Envoy route config rate limit actions
func createRateLimitActions(descriptors []v1alpha1.RateLimitDescriptor) ([]*routeconfv3.RateLimit_Action, error) {
	if len(descriptors) == 0 {
		return nil, errors.New("at least one descriptor is required for global rate limiting")
	}

	var result []*routeconfv3.RateLimit_Action

	// Process each descriptor
	for _, descriptor := range descriptors {
		// Each descriptor becomes a separate RateLimit in Envoy with its own set of actions
		// Create actions for each entry in the descriptor
		var actions []*routeconfv3.RateLimit_Action

		for _, entry := range descriptor.Entries {
			action := &routeconfv3.RateLimit_Action{}

			// Set the action specifier based on entry type
			switch entry.Type {
			case v1alpha1.RateLimitDescriptorEntryTypeGeneric:
				if entry.Generic == nil {
					return nil, fmt.Errorf("generic entry requires Generic field to be set")
				}
				action.ActionSpecifier = &routeconfv3.RateLimit_Action_GenericKey_{
					GenericKey: &routeconfv3.RateLimit_Action_GenericKey{
						DescriptorKey:   entry.Generic.Key,
						DescriptorValue: entry.Generic.Value,
					},
				}
			case v1alpha1.RateLimitDescriptorEntryTypeHeader:
				if entry.Header == "" {
					return nil, fmt.Errorf("header entry requires Header field to be set")
				}
				action.ActionSpecifier = &routeconfv3.RateLimit_Action_RequestHeaders_{
					RequestHeaders: &routeconfv3.RateLimit_Action_RequestHeaders{
						HeaderName:    entry.Header,
						DescriptorKey: entry.Header, // Use header name as key
					},
				}
			case v1alpha1.RateLimitDescriptorEntryTypeRemoteAddress:
				action.ActionSpecifier = &routeconfv3.RateLimit_Action_RemoteAddress_{
					RemoteAddress: &routeconfv3.RateLimit_Action_RemoteAddress{},
				}
			case v1alpha1.RateLimitDescriptorEntryTypePath:
				action.ActionSpecifier = &routeconfv3.RateLimit_Action_RequestHeaders_{
					RequestHeaders: &routeconfv3.RateLimit_Action_RequestHeaders{
						HeaderName:    ":path",
						DescriptorKey: "path",
					},
				}
			default:
				return nil, fmt.Errorf("unsupported entry type: %s", entry.Type)
			}

			actions = append(actions, action)
		}

		// If we have actions for this descriptor, add it
		if len(actions) > 0 {
			// In Envoy, a single RateLimit includes multiple Actions that together form a descriptor
			rateLimit := &routeconfv3.RateLimit{
				Actions: actions,
			}

			// The final result is a slice of complete RateLimit objects
			result = append(result, rateLimit.GetActions()...)
		}
	}

	return result, nil
}

func getRateLimitFilterName(name string) string {
	if name == "" {
		return rateLimitFilterNamePrefix
	}
	return fmt.Sprintf("%s/%s", rateLimitFilterNamePrefix, name)
}
