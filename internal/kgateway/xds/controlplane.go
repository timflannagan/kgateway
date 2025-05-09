package xds

import (
	"context"
	"errors"
	"net"

	envoy_service_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	envoy_service_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	envoy_service_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	envoy_service_route_v3 "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Config is the configuration for the xDS server.
type Config struct {
	Listener  net.Listener
	Callbacks server.Callbacks
	Delta     bool        // ADS=true | SOTW=false
	Logger    *zap.Logger // must not be nil

}

// ControlPlane is a wrapper around the xDS server. It is used to start the
// xDS server and stop it when the parent context is cancelled.
type ControlPlane struct {
	cache  cachev3.SnapshotCache
	server *grpc.Server
	lis    net.Listener
	errCh  chan error
}

// New creates a new xDS server.
func New(ctx context.Context, cfg *Config) (*ControlPlane, error) {
	if cfg.Logger == nil {
		cfg.Logger = zap.NewNop()
	}
	cache := cachev3.NewSnapshotCache(!cfg.Delta, NewNodeRoleHasher(), cfg.Logger.Sugar())
	srv := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(cfg.Logger),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				cfg.Logger.Debug("gRPC call", zap.String("method", info.FullMethod))
				return handler(srv, ss)
			}),
	))
	cp := &ControlPlane{
		cache:  cache,
		server: srv,
		lis:    cfg.Listener,
		errCh:  make(chan error, 1),
	}
	xdsSrv := serverv3.NewServer(ctx, cache, cfg.Callbacks)
	registerAllServices(srv, xdsSrv)
	return cp, nil
}

// Start is non-blocking; it returns the snapshot cache and an error channel.
// The error channel will be closed when the server is stopped.
func (cp *ControlPlane) Start(ctx context.Context) (cachev3.SnapshotCache, <-chan error) {
	// Shut the server down when the parent context is cancelled.
	go func() {
		<-ctx.Done()
		cp.server.GracefulStop()
	}()
	// Run Serve in its own goroutine and forward real failures.
	go func() {
		if err := cp.server.Serve(cp.lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			cp.errCh <- err
		}
		close(cp.errCh)
	}()

	return cp.cache, cp.errCh
}

// Stop stops the xDS server.
func (cp *ControlPlane) Stop() { cp.server.GracefulStop() }

// registerAllServices registers all xDS services with the gRPC server.
func registerAllServices(srv *grpc.Server, xdsSrv serverv3.Server) {
	envoy_service_endpoint_v3.RegisterEndpointDiscoveryServiceServer(srv, xdsSrv)
	envoy_service_cluster_v3.RegisterClusterDiscoveryServiceServer(srv, xdsSrv)
	envoy_service_route_v3.RegisterRouteDiscoveryServiceServer(srv, xdsSrv)
	envoy_service_listener_v3.RegisterListenerDiscoveryServiceServer(srv, xdsSrv)
	envoy_service_discovery_v3.RegisterAggregatedDiscoveryServiceServer(srv, xdsSrv)
}
