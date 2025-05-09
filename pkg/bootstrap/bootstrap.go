package bootstrap

import (
	"context"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	sdk "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk"
)

type bootstrap struct {
	logger      logr.Logger
	xdsOptions  XDSOptions
	MetricsAddr string
	HealthAddr  string
	// TODO: add leader election

	extraPlugins   func(ctx context.Context, commoncol *common.CommonCollections) []sdk.Plugin
	extraRunnables []manager.Runnable
}

type XDSOptions struct {
	GRPCAddr      string
	TLSCertPath   string
	TLSKeyPath    string
	SnapshotCache cache.SnapshotCache
}

type Option func(*bootstrap)

func WithLogger(logger logr.Logger) Option {
	return func(b *bootstrap) {
		b.logger = logger
	}
}

func WithExtraPlugins(extraPlugins func(ctx context.Context, commoncol *common.CommonCollections) []sdk.Plugin) Option {
	return func(b *bootstrap) {
		b.extraPlugins = extraPlugins
	}
}

func WithXDSOptions(xdsOptions XDSOptions) Option {
	return func(b *bootstrap) {
		b.xdsOptions = xdsOptions
	}
}

func WithMetricsAddr(metricsAddr string) Option {
	return func(b *bootstrap) {
		b.MetricsAddr = metricsAddr
	}
}

func WithHealthAddr(healthAddr string) Option {
	return func(b *bootstrap) {
		b.HealthAddr = healthAddr
	}
}

func WithExtraRunnables(extraRunnables ...manager.Runnable) Option {
	return func(b *bootstrap) {
		b.extraRunnables = extraRunnables
	}
}

func New(opts ...Option) *bootstrap {
	b := &bootstrap{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *bootstrap) Start(ctx context.Context) error {
	return nil
}
