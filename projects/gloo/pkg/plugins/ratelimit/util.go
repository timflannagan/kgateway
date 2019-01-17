package ratelimit

import (
	"errors"

	envoyvhostratelimit "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/rate_limit/v2"

	"github.com/gogo/protobuf/types"
	v1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
)

/*
translate virtual hosts
save them
then translate get rate limit configs
*/

func TranslateUserConfigToRateLimitServerConfig(vhostname string, ingressRl ratelimit.IngressRateLimit) (*v1.Constraint, error) {

	vhostConstraint := &v1.Constraint{
		Key:         genericKey,
		Value:       vhostname,
		Constraints: []*v1.Constraint{},
	}

	if ingressRl.AnonymousLimits != nil {

		if ingressRl.AnonymousLimits.Unit == ratelimit.RateLimit_UNKNOWN {
			return nil, errors.New("unknown unit for anonymous config")
		}

		c := &v1.Constraint{
			Key:   headerMatch,
			Value: anonymous,
			Constraints: []*v1.Constraint{
				{
					Key:       remoteAddress,
					RateLimit: ingressRl.AnonymousLimits,
				},
			},
		}

		vhostConstraint.Constraints = append(vhostConstraint.Constraints, c)
	}

	if ingressRl.AuthorizedLimits != nil {

		if ingressRl.AuthorizedLimits.Unit == ratelimit.RateLimit_UNKNOWN {
			return nil, errors.New("unknown unit for authenticated config")
		}

		c := &v1.Constraint{
			Key:   headerMatch,
			Value: authenticated,
			Constraints: []*v1.Constraint{
				{
					Key:       userid,
					RateLimit: ingressRl.AuthorizedLimits,
				},
			},
		}
		vhostConstraint.Constraints = append(vhostConstraint.Constraints, c)
	}

	return vhostConstraint, nil
}

func generateEnvoyConfigForFilter() *envoyratelimit.RateLimit {
	timeout := timeout
	envoyrl := envoyratelimit.RateLimit{
		Domain:      IngressDomain,
		Stage:       stage,
		RequestType: requestType,
		Timeout:     &timeout,
	}
	return &envoyrl
}

func generateEnvoyConfigForVhost(vhostname, headername string) []*envoyvhostratelimit.RateLimit {
	// the filter config, virtual host config are always the same:

	empty := headername == ""
	if empty {
		// TODO(yuval-k): fix this hack
		headername = "not-a-header"
	}

	vhostAction := getPerVhostRateLimit(vhostname)

	vhostrl := []*envoyvhostratelimit.RateLimit{
		{
			Stage: &types.UInt32Value{Value: stage},
			Actions: []*envoyvhostratelimit.RateLimit_Action{
				vhostAction,
				getAuthHeaderRateLimit(headername, true),
				getUserIdRateLimit(headername),
			},
		},
		{
			Stage: &types.UInt32Value{Value: stage},
			Actions: []*envoyvhostratelimit.RateLimit_Action{
				vhostAction,
				getAuthHeaderRateLimit(vhostname, false),
				getPerIpRateLimit(),
			},
		},
	}
	return vhostrl
}

func getPerVhostRateLimit(vhostname string) *envoyvhostratelimit.RateLimit_Action {
	return &envoyvhostratelimit.RateLimit_Action{
		ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_GenericKey_{
			GenericKey: &envoyvhostratelimit.RateLimit_Action_GenericKey{
				DescriptorValue: vhostname,
			},
		},
	}
}

func getAuthHeaderRateLimit(headername string, match bool) *envoyvhostratelimit.RateLimit_Action {

	headersmatcher := []*envoyvhostratelimit.HeaderMatcher{{
		Name:                 headername,
		HeaderMatchSpecifier: &envoyvhostratelimit.HeaderMatcher_PresentMatch{PresentMatch: true},
	}}

	var value string
	if match {
		value = authenticated
	} else {
		value = anonymous
	}

	return &envoyvhostratelimit.RateLimit_Action{
		ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_HeaderValueMatch_{
			HeaderValueMatch: &envoyvhostratelimit.RateLimit_Action_HeaderValueMatch{

				DescriptorValue: value,
				ExpectMatch:     &types.BoolValue{Value: match},
				Headers:         headersmatcher,
			},
		},
	}
}

func getUserIdRateLimit(headername string) *envoyvhostratelimit.RateLimit_Action {
	return &envoyvhostratelimit.RateLimit_Action{
		ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_RequestHeaders_{
			RequestHeaders: &envoyvhostratelimit.RateLimit_Action_RequestHeaders{
				DescriptorKey: userid,
				HeaderName:    headername,
			},
		},
	}
}

func getPerIpRateLimit() *envoyvhostratelimit.RateLimit_Action {
	return &envoyvhostratelimit.RateLimit_Action{
		ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_RemoteAddress_{
			RemoteAddress: &envoyvhostratelimit.RateLimit_Action_RemoteAddress{},
		},
	}
}
