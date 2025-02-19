package setup

import (
	"context"
	"net"
	"os"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xdsserver "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/go-logr/zapr"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	istiokube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/krt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	zaputil "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/admin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/controller"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils/krtutil"
	"github.com/kgateway-dev/kgateway/v2/internal/version"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/envutils"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/namespaces"
)

const (
	componentName = "kgateway"
)

type Option func(*Controller)

type Controller struct {
	leaderElection bool
	extraGwClasses []string
	extraPlugins   []extensionsplug.Plugin
}

func New(opts ...Option) *Controller {
	c := &Controller{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithLeaderElection() Option {
	return func(c *Controller) {
		c.leaderElection = true
	}
}

func WithExtraGwClasses(extraGwClasses []string) Option {
	return func(c *Controller) {
		c.extraGwClasses = extraGwClasses
	}
}

func WithExtraPlugins(extraPlugins []extensionsplug.Plugin) Option {
	return func(c *Controller) {
		c.extraPlugins = extraPlugins
	}
}

func (c *Controller) Start(ctx context.Context) error {
	setupLogging(ctx, componentName)
	restConfig, err := ctrl.GetConfig()
	if err != nil {
		return err
	}
	uniqueClientCallbacks, uccBuilder := krtcollections.NewUniquelyConnectedClients()
	cache, err := startControlPlane(ctx, uniqueClientCallbacks)
	if err != nil {
		return err
	}
	setupOpts := &controller.SetupOpts{
		Cache:               cache,
		KrtDebugger:         new(krt.DebugHandler),
		ExtraGatewayClasses: c.extraGwClasses,
		XdsHost:             GetControlPlaneXdsHost(),
		XdsPort:             9977,
		LeaderElection:      c.leaderElection,
	}
	return StartGGv2WithConfig(ctx, setupOpts, restConfig, uccBuilder, c.extraPlugins, c.extraGwClasses)
}

// GetControlPlaneXdsHost gets the xDS address from the gloo Service.
func GetControlPlaneXdsHost() string {
	return kubeutils.ServiceFQDN(metav1.ObjectMeta{
		Name:      kubeutils.GlooServiceName,
		Namespace: namespaces.GetPodNamespace(),
	})
}

func startControlPlane(ctx context.Context,
	callbacks xdsserver.Callbacks) (envoycache.SnapshotCache, error) {

	return NewControlPlane(ctx, &net.TCPAddr{IP: net.IPv4zero, Port: 9977}, callbacks)
}

func createIstioClient(restConfig *rest.Config) (istiokube.Client, error) {
	restCfg := istiokube.NewClientConfigForRestConfig(restConfig)
	client, err := istiokube.NewClient(restCfg, "")
	if err != nil {
		return nil, err
	}
	istiokube.EnableCrdWatcher(client)
	return client, nil
}

// TODO: Rename this function.
func StartGGv2WithConfig(
	ctx context.Context,
	setupOpts *controller.SetupOpts,
	restConfig *rest.Config,
	uccBuilder krtcollections.UniquelyConnectedClientsBulider,
	extraPlugins []extensionsplug.Plugin,
	extraGwClasses []string, // TODO: we can remove this and replace with something that watches all GW classes with our controller name
) error {
	ctx = contextutils.WithLogger(ctx, "k8s")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("starting %s", componentName)

	kubeClient, err := createIstioClient(restConfig)
	if err != nil {
		return err
	}

	logger.Info("creating krt collections")
	krtOpts := krtutil.NewKrtOptions(ctx.Done(), setupOpts.KrtDebugger)

	augmentedPods := krtcollections.NewPodsCollection(kubeClient, krtOpts)
	augmentedPodsForUcc := augmentedPods
	if envutils.IsEnvTruthy("DISABLE_POD_LOCALITY_XDS") {
		augmentedPodsForUcc = nil
	}

	ucc := uccBuilder(ctx, krtOpts, augmentedPodsForUcc)

	logger.Info("initializing controller")
	c, err := controller.NewControllerBuilder(ctx, controller.StartConfig{
		ExtraPlugins:  extraPlugins,
		RestConfig:    restConfig,
		SetupOpts:     setupOpts,
		Client:        kubeClient,
		AugmentedPods: augmentedPods,
		UniqueClients: ucc,
		// Dev flag may be useful for development purposes; not currently tied to any user-facing API
		Dev:        os.Getenv("LOG_LEVEL") == "debug",
		KrtOptions: krtOpts,
	})
	if err != nil {
		logger.Error("failed initializing controller: ", err)
		return err
	}
	/// no collections after this point

	logger.Info("waiting for cache sync")
	kubeClient.RunAndWait(ctx.Done())

	logger.Info("starting admin server")
	go admin.RunAdminServer(ctx, setupOpts)

	logger.Info("starting controller")
	return c.Start(ctx)
}

// setupLogging sets up controller-runtime logging
func setupLogging(ctx context.Context, loggerName string) {
	level := zapcore.InfoLevel
	// if log level is set in env, use that
	if envLogLevel := os.Getenv(contextutils.LogLevelEnvName); envLogLevel != "" {
		if err := (&level).Set(envLogLevel); err != nil {
			contextutils.LoggerFrom(ctx).Infof("Could not set log level from env %s=%s, available levels "+
				"can be found here: https://pkg.go.dev/go.uber.org/zap/zapcore?tab=doc#Level",
				contextutils.LogLevelEnvName,
				envLogLevel,
				zap.Error(err),
			)
		}
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	baseLogger := zaputil.NewRaw(
		zaputil.Level(&atomicLevel),
		zaputil.RawZapOpts(zap.Fields(zap.String("version", version.Version))),
	).Named(loggerName)

	// controller-runtime
	log.SetLogger(zapr.NewLogger(baseLogger))
}
