// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package devportal

import (
	"context"

	"github.com/solo-io/dev-portal/pkg/admin/grpc/portal"
	"github.com/solo-io/dev-portal/pkg/assets"
	"github.com/solo-io/go-utils/envutils"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"
)

// Injectors from wire.go:

func InitDevPortal(ctx context.Context) (Registrar, error) {
	config, err := setup.NewKubeConfig()
	if err != nil {
		return nil, err
	}
	string2 := envutils.MustGetPodNamespace(ctx)
	manager, err := NewManager(ctx, config, string2)
	if err != nil {
		return nil, err
	}
	client := NewDynamicClient(manager)
	portalClient := NewPortalClient(client)
	configMapClient := NewConfigMapClient(client)
	configMapStorage := assets.NewConfigMapStorage(configMapClient)
	grpcService := portal.NewPortalGrpcService(portalClient, configMapStorage)
	registrar := NewRegistrar(grpcService)
	return registrar, nil
}
