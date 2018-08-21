package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	"github.com/solo-io/solo-kit/pkg/api/v1/propagator"
)

func Setup(namespace string, inputResourceOpts factory.ResourceClientFactoryOpts, opts clients.WatchOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "setup")
	inputFactory := factory.NewResourceClientFactory(inputResourceOpts)

	gatewayClient, err := v1.NewGatewayClient(inputFactory)
	if err != nil {
		return err
	}

	virtualServicesClient, err := v1.NewVirtualServiceClient(inputFactory)
	if err != nil {
		return err
	}

	proxyClient, err := gloov1.NewProxyClient(inputFactory)
	if err != nil {
		return err
	}

	cache := v1.NewCache(gatewayClient, virtualServicesClient)

	rpt := reporter.NewReporter("gateway", gatewayClient.BaseClient(), virtualServicesClient.BaseClient())

	sync := syncer.NewSyncer(namespace, proxyClient, rpt)

	eventLoop := v1.NewEventLoop(cache, sync)
	eventLoop.Run(namespace, opts)

	errs := make(chan error)

	logger := contextutils.LoggerFrom(opts.Ctx)

	for {
		select {
		case err := <-errs:
			logger.Errorf("error: %v", err)
		case <-opts.Ctx.Done():
			close(errs)
			return nil
		}
	}
}
