package extensions2

import sdk "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk"

// **Re-export** so existing core imports keep compiling.
type (
	Plugin              = sdk.Plugin
	PolicyPlugin        = sdk.PolicyPlugin
	BackendPlugin       = sdk.BackendPlugin
	ContributesPolicies = sdk.ContributesPolicies
	ContributesBackends = sdk.ContributesBackends
)

var Register = sdk.Register
