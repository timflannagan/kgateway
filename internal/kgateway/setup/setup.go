package setup

import (
	"context"
	"fmt"
	"net"
	"os"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/go-logr/zapr"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	istiokube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/krt"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaputil "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/admin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/controller"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/settings"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/xds"
	"github.com/kgateway-dev/kgateway/v2/internal/version"
	sdk "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/envutils"
)

const (
	componentName = "kgateway"
)

// Server is the interface for the kgateway server.
type Server interface {
	// Start starts the kgateway server.
	Start(ctx context.Context) error
}

// setup is the internal implementation of the kgateway server.
type setup struct {
	logger       *zap.Logger
	extraPlugins func(ctx context.Context, commoncol *common.CommonCollections) []sdk.Plugin

	// previously hardcoded in the implementation. TODO add to CLI.
	metricsAddr     string
	healthProbeAddr string
	pprofAddr       string
	restConfig      *rest.Config // useful in integration tests
	logLevel        string       // let CLI decide env vs flag
	controllerName  string
	xdsOptions      *xds.Config
}

var _ Server = &setup{}

// New creates a new kgateway server.
func New(opts ...func(*setup)) *setup {
	s := &setup{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithExtraPlugins extends the kgateway server with extra plugins.
func WithExtraPlugins(extraPlugins func(ctx context.Context, commoncol *common.CommonCollections) []sdk.Plugin) func(*setup) {
	return func(s *setup) {
		s.extraPlugins = extraPlugins
	}
}

// WithMetricsAddr sets the metrics address.
func WithMetricsAddr(metricsAddr string) func(*setup) {
	return func(s *setup) {
		s.metricsAddr = metricsAddr
	}
}

// WithHealthProbeAddr sets the health probe address.
func WithHealthProbeAddr(healthProbeAddr string) func(*setup) {
	return func(s *setup) {
		s.healthProbeAddr = healthProbeAddr
	}
}

// WithPprofAddr sets the pprof address.
func WithPprofAddr(pprofAddr string) func(*setup) {
	return func(s *setup) {
		s.pprofAddr = pprofAddr
	}
}

// WithRestConfig sets the rest config.
func WithRestConfig(restConfig *rest.Config) func(*setup) {
	return func(s *setup) {
		s.restConfig = restConfig
	}
}

// WithLogLevel sets the log level.
func WithLogLevel(logLevel string) func(*setup) {
	return func(s *setup) {
		s.logLevel = logLevel
	}
}

// WithControllerName sets the controller name.
func WithControllerName(controllerName string) func(*setup) {
	return func(s *setup) {
		s.controllerName = controllerName
	}
}

// WithXdsOptions sets the xDS options.
func WithXdsOptions(xdsOptions *xds.Config) func(*setup) {
	return func(s *setup) {
		s.xdsOptions = xdsOptions
	}
}

// Start starts the kgateway server. It will block until the server is stopped.
func (s *setup) Start(ctx context.Context) error {
	logger, err := setupLogging(ctx, componentName)
	if err != nil {
		return fmt.Errorf("failed to setup logging: %w", err)
	}
	s.logger = logger

	// load global settings
	st, err := loadGlobalSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to load global settings: %w", err)
	}
	s.logger.Info("loaded global settings", zap.Any("settings", st))

	// start the xDS server
	cache, uccBuilder, err := startXds(ctx, logger, st)
	if err != nil {
		return fmt.Errorf("failed to start xDS server: %w", err)
	}

	// setup the controller and run it
	return s.run(ctx, uccBuilder, &controller.SetupOpts{
		Cache:                  cache,
		KrtDebugger:            new(krt.DebugHandler),
		GlobalSettings:         st,
		PprofBindAddress:       s.pprofAddr,
		HealthProbeBindAddress: s.healthProbeAddr,
		MetricsBindAddress:     s.metricsAddr,
	})
}

// loadGlobalSettings loads the global settings from the environment. Any
// environment variable prefixed with "KGW_" will be used to set a setting.
func loadGlobalSettings(ctx context.Context) (*settings.Settings, error) {
	st, err := settings.BuildSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to build settings: %w", err)
	}
	return st, nil
}

// startXds starts the xDS server. It returns a cache and a builder for the
// uniquely connected clients that will be used by the controller.
func startXds(
	ctx context.Context,
	logger *zap.Logger,
	st *settings.Settings,
) (cache envoycache.SnapshotCache, uccBuilder krtcollections.UniquelyConnectedClientsBulider, err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", st.XdsServicePort))
	if err != nil {
		return nil, nil, err
	}

	cb, uccBuilder := krtcollections.NewUniquelyConnectedClients()
	cp, err := xds.New(ctx, &xds.Config{
		Listener:  lis,
		Callbacks: cb,
		Delta:     true,
		Logger:    logger,
	})
	if err != nil {
		return nil, nil, err
	}

	cache, errCh := cp.Start(ctx)
	go func() {
		if err := <-errCh; err != nil {
			logger.Error("xDS serve failed", zap.Error(err))
		}
	}()
	return cache, uccBuilder, nil
}

// run is responsible for running the kgateway server. It will block until the
// server is stopped. This will setup the controller-runtime manager, it's
// controllers, KRT collections, and the admin server.
//
// TODO(tim): this still seems very awkward.
func (s *setup) run(
	ctx context.Context,
	uccBuilder krtcollections.UniquelyConnectedClientsBulider,
	setupOpts *controller.SetupOpts,
) error {
	s.logger.Info("starting the kgateway server")

	restConfig := ctrl.GetConfigOrDie()
	if s.restConfig == nil {
		s.restConfig = restConfig
	}

	kubeClient, err := createKubeClient(s.restConfig)
	if err != nil {
		return err
	}

	s.logger.Info("creating krt collections")
	krtOpts := krtutil.NewKrtOptions(ctx.Done(), setupOpts.KrtDebugger)

	augmentedPods := krtcollections.NewPodsCollection(kubeClient, krtOpts)
	augmentedPodsForUcc := augmentedPods
	// Should this be a Setting or CLI flag?
	if envutils.IsEnvTruthy("DISABLE_POD_LOCALITY_XDS") {
		augmentedPodsForUcc = nil
	}

	// use the controller name from the CLI flag if provided, otherwise use the
	// default wellknown controller name.
	controllerName := s.controllerName
	if controllerName == "" {
		controllerName = wellknown.GatewayControllerName
	}

	s.logger.Info("initializing controller")
	c, err := controller.NewControllerBuilder(ctx, controller.StartConfig{
		ControllerName: controllerName,
		ExtraPlugins:   s.extraPlugins,
		RestConfig:     s.restConfig,
		SetupOpts:      setupOpts,
		Client:         kubeClient,
		AugmentedPods:  augmentedPods,
		UniqueClients:  uccBuilder(ctx, krtOpts, augmentedPodsForUcc),
		Dev:            s.logLevel == "debug",
		KrtOptions:     krtOpts,
	})
	if err != nil {
		s.logger.Error("failed initializing controller: ", zap.Error(err))
		return err
	}
	/// no collections after this point

	s.logger.Info("waiting for cache sync")
	kubeClient.RunAndWait(ctx.Done())

	s.logger.Info("starting admin server")
	go admin.RunAdminServer(ctx, setupOpts)

	s.logger.Info("starting controller")
	return c.Start(ctx)
}

func buildLogger(logLevel string) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if logLevel != "" {
		if err := (&level).Set(logLevel); err != nil {
			return nil, fmt.Errorf("failed to set log level from env %s=%s: %w",
				contextutils.LogLevelEnvName,
				logLevel,
				err,
			)
		}
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	return zaputil.NewRaw(
		zaputil.Level(&atomicLevel),
		zaputil.RawZapOpts(zap.Fields(zap.String("version", version.Version))),
	).Named(componentName), nil
}

// setupLogging sets up controller-runtime logging
func setupLogging(ctx context.Context, loggerName string) (*zap.Logger, error) {
	baseLogger, err := buildLogger(os.Getenv(contextutils.LogLevelEnvName))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}
	// controller-runtime
	log.SetLogger(zapr.NewLogger(baseLogger))
	return baseLogger, nil
}

// createKubeClient creates a new Istio Kubernetes client.
func createKubeClient(restConfig *rest.Config) (istiokube.Client, error) {
	restCfg := istiokube.NewClientConfigForRestConfig(restConfig)
	client, err := istiokube.NewClient(restCfg, "")
	if err != nil {
		return nil, err
	}
	istiokube.EnableCrdWatcher(client)
	return client, nil
}
