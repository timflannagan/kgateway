package setup

import (
	core "github.com/kgateway-dev/kgateway/v2/internal/kgateway/setup"
	sdk "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk"
)

type Options struct {
	ExtraPlugins []sdk.Plugin
}

func New(opts Options) core.Server {
	// internal setup already accepted functional-options; we wrap only extras.
	return core.New(core.WithExtraPlugins(opts.ExtraPlugins...))
}
