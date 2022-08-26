// Code generated by skv2. DO NOT EDIT.

// Definition for federated resource reconciler templates.
package federation

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/stringutils"
	"github.com/solo-io/skv2/pkg/reconcile"
	gateway_solo_io_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gateway_solo_io_v1_sets "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1/sets"
	fed_gateway_solo_io_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1/controller"
	fed_gateway_solo_io_v1_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1/types"
	mc_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation/placement"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type federatedGatewayReconciler struct {
	ctx                  context.Context
	federatedGateways    fed_gateway_solo_io_v1.FederatedGatewayClient
	baseClients          gateway_solo_io_v1.MulticlusterClientset
	statusBuilderFactory placement.StatusBuilderFactory
	clusterSet           multicluster.ClusterSet
}

func NewFederatedGatewayReconciler(
	ctx context.Context,
	federatedGateways fed_gateway_solo_io_v1.FederatedGatewayClient,
	baseClients gateway_solo_io_v1.MulticlusterClientset,
	statusBuilderFactory placement.StatusBuilderFactory,
	clusterSet multicluster.ClusterSet,
) controller.FederatedGatewayFinalizer {
	return &federatedGatewayReconciler{
		ctx:                  ctx,
		federatedGateways:    federatedGateways,
		baseClients:          baseClients,
		statusBuilderFactory: statusBuilderFactory,
		clusterSet:           clusterSet,
	}
}

func (f *federatedGatewayReconciler) ReconcileFederatedGateway(obj *fed_gateway_solo_io_v1.FederatedGateway) (reconcile.Result, error) {
	if !obj.NeedsReconcile() {
		return reconcile.Result{}, nil
	}

	contextutils.LoggerFrom(f.ctx).Debugw("processing federated gateway", zap.Any("FederatedGateway", obj))
	statusBuilder := f.statusBuilderFactory.GetBuilder()

	allClusters := f.clusterSet.ListClusters()

	// Validate resource
	if obj.Spec.GetPlacement() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.PlacementMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedGateways.UpdateFederatedGatewayStatus(f.ctx, obj)
	}
	for _, cluster := range obj.Spec.Placement.GetClusters() {
		if !stringutils.ContainsString(cluster, allClusters) {
			obj.Status.PlacementStatus = statusBuilder.
				UpdateUnprocessed(obj.Status.PlacementStatus, placement.ClusterNotRegistered(cluster), mc_types.PlacementStatus_INVALID).
				Eject(obj.GetGeneration())
			return reconcile.Result{}, f.federatedGateways.UpdateFederatedGatewayStatus(f.ctx, obj)
		}
	}
	if obj.Spec.Template.GetSpec() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.SpecTemplateMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedGateways.UpdateFederatedGatewayStatus(f.ctx, obj)
	}
	if obj.Spec.Template.GetMetadata() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.MetaTemplateMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedGateways.UpdateFederatedGatewayStatus(f.ctx, obj)
	}

	// ownerLabel is used to reference Federated resources via their children.
	ownerLabel := federation.GetOwnerLabel(obj)

	spec := obj.Spec.Template.GetSpec()
	meta := obj.Spec.Template.GetMetadata()
	labels := federation.Merge(meta.GetLabels(), ownerLabel)

	multiErr := &multierror.Error{}
	for _, cluster := range allClusters {
		clusterGateways := gateway_solo_io_v1_sets.NewGatewaySet()
		if stringutils.ContainsString(cluster, obj.Spec.Placement.GetClusters()) {
			for _, namespace := range obj.Spec.Placement.GetNamespaces() {

				clusterGateways.Insert(&gateway_solo_io_v1.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:   namespace,
						Name:        meta.GetName(),
						Labels:      labels,
						Annotations: meta.GetAnnotations(),
					},
					Spec: *spec,
				})
			}
		}

		if err := f.ensureCluster(cluster, statusBuilder, clusterGateways, ownerLabel); err != nil {
			multiErr.Errors = append(multiErr.Errors, err)
		}
	}

	obj.Status = fed_gateway_solo_io_v1_types.FederatedGatewayStatus{
		PlacementStatus: statusBuilder.Build(obj.GetGeneration()),
	}
	err := f.federatedGateways.UpdateFederatedGatewayStatus(f.ctx, obj)
	if err != nil {
		multiErr.Errors = append(multiErr.Errors, err)
		contextutils.LoggerFrom(f.ctx).Errorw("Failed to update status on federated gateway", zap.Error(err))
	}

	return reconcile.Result{}, multiErr.ErrorOrNil()
}

func (f *federatedGatewayReconciler) FederatedGatewayFinalizerName() string {
	return federation.HubFinalizer
}

func (f *federatedGatewayReconciler) FinalizeFederatedGateway(obj *fed_gateway_solo_io_v1.FederatedGateway) error {
	return f.deleteAll(federation.GetOwnerLabel(obj))
}

// ensureCluster upserts all desired resources on the given cluster.
// An error is returned only if a retry is expected to resolve the issue.
func (f *federatedGatewayReconciler) ensureCluster(cluster string, statusBuilder placement.StatusBuilder, desired gateway_solo_io_v1_sets.GatewaySet, ownerLabel map[string]string) error {
	clientset, err := f.baseClients.Cluster(cluster)
	if err != nil {
		var namespaces []string
		for _, obj := range desired.List() {
			namespaces = append(namespaces, obj.GetNamespace())
		}

		contextutils.LoggerFrom(f.ctx).Errorw("Failed to get clientset", zap.String("cluster", cluster), zap.Error(err))
		statusBuilder.AddDestinations([]string{cluster}, namespaces, mc_types.PlacementStatus_Namespace{
			State:   mc_types.PlacementStatus_FAILED,
			Message: placement.FailedToCreateClientForCluster(cluster),
		})
		return err
	}

	gatewayClient := clientset.Gateways()

	existingList, err := gatewayClient.ListGateway(f.ctx, client.MatchingLabels(ownerLabel))
	if err != nil {
		var namespaces []string
		for _, obj := range desired.List() {
			namespaces = append(namespaces, obj.GetNamespace())
		}

		contextutils.LoggerFrom(f.ctx).Errorw("Failed to list gateways", zap.Error(err))
		statusBuilder.AddDestinations([]string{cluster}, namespaces, mc_types.PlacementStatus_Namespace{
			State:   mc_types.PlacementStatus_FAILED,
			Message: placement.FailedToListResource("gateway", cluster),
		})
		return err
	}

	existing := gateway_solo_io_v1_sets.NewGatewaySet()
	for _, gateway := range existingList.Items {
		gatewayPointer := gateway
		existing.Insert(&gatewayPointer)
	}

	multiErr := &multierror.Error{}
	for _, desiredGateway := range desired.List() {
		err := gatewayClient.UpsertGateway(f.ctx, desiredGateway)
		if err != nil && errors.IsConflict(err) {
			multiErr.Errors = append(multiErr.Errors, err)
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to upsert gateway due to resource conflict", zap.Error(err))
			statusBuilder.AddDestination(cluster, desiredGateway.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_FAILED,
				Message: placement.FailedToUpsertResourceDueToConflict("gateway"),
			})
		} else if err != nil {
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to upsert gateway", zap.Error(err))
			statusBuilder.AddDestination(cluster, desiredGateway.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_FAILED,
				Message: placement.FailedToUpsertResource("gateway"),
			})
		} else {
			statusBuilder.AddDestination(cluster, desiredGateway.Namespace, mc_types.PlacementStatus_Namespace{
				State: mc_types.PlacementStatus_PLACED,
			})
		}
	}

	for _, staleGateway := range existing.Difference(desired).List() {
		err := gatewayClient.DeleteGateway(f.ctx, client.ObjectKey{
			Namespace: staleGateway.Namespace,
			Name:      staleGateway.Name,
		})
		if client.IgnoreNotFound(err) != nil {
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to delete gateway", zap.Error(err))
			statusBuilder.AddDestination(cluster, staleGateway.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_STALE,
				Message: placement.FailedToDeleteResource("gateway"),
			})
		}
	}

	return multiErr.ErrorOrNil()
}

// Delete all gateways in the matching the ownerLabel on all clusters.
// Used to ensure that gateways generated by a FederatedGateway are cleaned up on delete.
func (f *federatedGatewayReconciler) deleteAll(ownerLabel map[string]string) error {
	for _, cluster := range f.clusterSet.ListClusters() {
		clusterClient, err := f.baseClients.Cluster(cluster)
		if err != nil {
			return err
		}
		// TODO this requires permissions in all namespaces, we could restrict to to namespaces referenced by gloo instances
		list, err := clusterClient.Gateways().ListGateway(f.ctx, client.MatchingLabels(ownerLabel))
		if err != nil {
			return err
		}

		for _, gateway := range list.Items {
			err = clusterClient.Gateways().DeleteGateway(f.ctx, client.ObjectKey{
				Namespace: gateway.Namespace,
				Name:      gateway.Name,
			})
			if client.IgnoreNotFound(err) != nil {
				return err
			}
		}
	}
	return nil
}

type federatedVirtualServiceReconciler struct {
	ctx                      context.Context
	federatedVirtualServices fed_gateway_solo_io_v1.FederatedVirtualServiceClient
	baseClients              gateway_solo_io_v1.MulticlusterClientset
	statusBuilderFactory     placement.StatusBuilderFactory
	clusterSet               multicluster.ClusterSet
}

func NewFederatedVirtualServiceReconciler(
	ctx context.Context,
	federatedVirtualServices fed_gateway_solo_io_v1.FederatedVirtualServiceClient,
	baseClients gateway_solo_io_v1.MulticlusterClientset,
	statusBuilderFactory placement.StatusBuilderFactory,
	clusterSet multicluster.ClusterSet,
) controller.FederatedVirtualServiceFinalizer {
	return &federatedVirtualServiceReconciler{
		ctx:                      ctx,
		federatedVirtualServices: federatedVirtualServices,
		baseClients:              baseClients,
		statusBuilderFactory:     statusBuilderFactory,
		clusterSet:               clusterSet,
	}
}

func (f *federatedVirtualServiceReconciler) ReconcileFederatedVirtualService(obj *fed_gateway_solo_io_v1.FederatedVirtualService) (reconcile.Result, error) {
	if !obj.NeedsReconcile() {
		return reconcile.Result{}, nil
	}

	contextutils.LoggerFrom(f.ctx).Debugw("processing federated virtualService", zap.Any("FederatedVirtualService", obj))
	statusBuilder := f.statusBuilderFactory.GetBuilder()

	allClusters := f.clusterSet.ListClusters()

	// Validate resource
	if obj.Spec.GetPlacement() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.PlacementMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedVirtualServices.UpdateFederatedVirtualServiceStatus(f.ctx, obj)
	}
	for _, cluster := range obj.Spec.Placement.GetClusters() {
		if !stringutils.ContainsString(cluster, allClusters) {
			obj.Status.PlacementStatus = statusBuilder.
				UpdateUnprocessed(obj.Status.PlacementStatus, placement.ClusterNotRegistered(cluster), mc_types.PlacementStatus_INVALID).
				Eject(obj.GetGeneration())
			return reconcile.Result{}, f.federatedVirtualServices.UpdateFederatedVirtualServiceStatus(f.ctx, obj)
		}
	}
	if obj.Spec.Template.GetSpec() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.SpecTemplateMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedVirtualServices.UpdateFederatedVirtualServiceStatus(f.ctx, obj)
	}
	if obj.Spec.Template.GetMetadata() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.MetaTemplateMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedVirtualServices.UpdateFederatedVirtualServiceStatus(f.ctx, obj)
	}

	// ownerLabel is used to reference Federated resources via their children.
	ownerLabel := federation.GetOwnerLabel(obj)

	spec := obj.Spec.Template.GetSpec()
	meta := obj.Spec.Template.GetMetadata()
	labels := federation.Merge(meta.GetLabels(), ownerLabel)

	multiErr := &multierror.Error{}
	for _, cluster := range allClusters {
		clusterVirtualServices := gateway_solo_io_v1_sets.NewVirtualServiceSet()
		if stringutils.ContainsString(cluster, obj.Spec.Placement.GetClusters()) {
			for _, namespace := range obj.Spec.Placement.GetNamespaces() {

				clusterVirtualServices.Insert(&gateway_solo_io_v1.VirtualService{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:   namespace,
						Name:        meta.GetName(),
						Labels:      labels,
						Annotations: meta.GetAnnotations(),
					},
					Spec: *spec,
				})
			}
		}

		if err := f.ensureCluster(cluster, statusBuilder, clusterVirtualServices, ownerLabel); err != nil {
			multiErr.Errors = append(multiErr.Errors, err)
		}
	}

	obj.Status = fed_gateway_solo_io_v1_types.FederatedVirtualServiceStatus{
		PlacementStatus: statusBuilder.Build(obj.GetGeneration()),
	}
	err := f.federatedVirtualServices.UpdateFederatedVirtualServiceStatus(f.ctx, obj)
	if err != nil {
		multiErr.Errors = append(multiErr.Errors, err)
		contextutils.LoggerFrom(f.ctx).Errorw("Failed to update status on federated virtualService", zap.Error(err))
	}

	return reconcile.Result{}, multiErr.ErrorOrNil()
}

func (f *federatedVirtualServiceReconciler) FederatedVirtualServiceFinalizerName() string {
	return federation.HubFinalizer
}

func (f *federatedVirtualServiceReconciler) FinalizeFederatedVirtualService(obj *fed_gateway_solo_io_v1.FederatedVirtualService) error {
	return f.deleteAll(federation.GetOwnerLabel(obj))
}

// ensureCluster upserts all desired resources on the given cluster.
// An error is returned only if a retry is expected to resolve the issue.
func (f *federatedVirtualServiceReconciler) ensureCluster(cluster string, statusBuilder placement.StatusBuilder, desired gateway_solo_io_v1_sets.VirtualServiceSet, ownerLabel map[string]string) error {
	clientset, err := f.baseClients.Cluster(cluster)
	if err != nil {
		var namespaces []string
		for _, obj := range desired.List() {
			namespaces = append(namespaces, obj.GetNamespace())
		}

		contextutils.LoggerFrom(f.ctx).Errorw("Failed to get clientset", zap.String("cluster", cluster), zap.Error(err))
		statusBuilder.AddDestinations([]string{cluster}, namespaces, mc_types.PlacementStatus_Namespace{
			State:   mc_types.PlacementStatus_FAILED,
			Message: placement.FailedToCreateClientForCluster(cluster),
		})
		return err
	}

	virtualServiceClient := clientset.VirtualServices()

	existingList, err := virtualServiceClient.ListVirtualService(f.ctx, client.MatchingLabels(ownerLabel))
	if err != nil {
		var namespaces []string
		for _, obj := range desired.List() {
			namespaces = append(namespaces, obj.GetNamespace())
		}

		contextutils.LoggerFrom(f.ctx).Errorw("Failed to list virtualServices", zap.Error(err))
		statusBuilder.AddDestinations([]string{cluster}, namespaces, mc_types.PlacementStatus_Namespace{
			State:   mc_types.PlacementStatus_FAILED,
			Message: placement.FailedToListResource("virtualService", cluster),
		})
		return err
	}

	existing := gateway_solo_io_v1_sets.NewVirtualServiceSet()
	for _, virtualService := range existingList.Items {
		virtualServicePointer := virtualService
		existing.Insert(&virtualServicePointer)
	}

	multiErr := &multierror.Error{}
	for _, desiredVirtualService := range desired.List() {
		err := virtualServiceClient.UpsertVirtualService(f.ctx, desiredVirtualService)
		if err != nil && errors.IsConflict(err) {
			multiErr.Errors = append(multiErr.Errors, err)
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to upsert virtualService due to resource conflict", zap.Error(err))
			statusBuilder.AddDestination(cluster, desiredVirtualService.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_FAILED,
				Message: placement.FailedToUpsertResourceDueToConflict("virtualService"),
			})
		} else if err != nil {
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to upsert virtualService", zap.Error(err))
			statusBuilder.AddDestination(cluster, desiredVirtualService.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_FAILED,
				Message: placement.FailedToUpsertResource("virtualService"),
			})
		} else {
			statusBuilder.AddDestination(cluster, desiredVirtualService.Namespace, mc_types.PlacementStatus_Namespace{
				State: mc_types.PlacementStatus_PLACED,
			})
		}
	}

	for _, staleVirtualService := range existing.Difference(desired).List() {
		err := virtualServiceClient.DeleteVirtualService(f.ctx, client.ObjectKey{
			Namespace: staleVirtualService.Namespace,
			Name:      staleVirtualService.Name,
		})
		if client.IgnoreNotFound(err) != nil {
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to delete virtualService", zap.Error(err))
			statusBuilder.AddDestination(cluster, staleVirtualService.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_STALE,
				Message: placement.FailedToDeleteResource("virtualService"),
			})
		}
	}

	return multiErr.ErrorOrNil()
}

// Delete all virtualServices in the matching the ownerLabel on all clusters.
// Used to ensure that virtualServices generated by a FederatedVirtualService are cleaned up on delete.
func (f *federatedVirtualServiceReconciler) deleteAll(ownerLabel map[string]string) error {
	for _, cluster := range f.clusterSet.ListClusters() {
		clusterClient, err := f.baseClients.Cluster(cluster)
		if err != nil {
			return err
		}
		// TODO this requires permissions in all namespaces, we could restrict to to namespaces referenced by gloo instances
		list, err := clusterClient.VirtualServices().ListVirtualService(f.ctx, client.MatchingLabels(ownerLabel))
		if err != nil {
			return err
		}

		for _, virtualService := range list.Items {
			err = clusterClient.VirtualServices().DeleteVirtualService(f.ctx, client.ObjectKey{
				Namespace: virtualService.Namespace,
				Name:      virtualService.Name,
			})
			if client.IgnoreNotFound(err) != nil {
				return err
			}
		}
	}
	return nil
}

type federatedRouteTableReconciler struct {
	ctx                  context.Context
	federatedRouteTables fed_gateway_solo_io_v1.FederatedRouteTableClient
	baseClients          gateway_solo_io_v1.MulticlusterClientset
	statusBuilderFactory placement.StatusBuilderFactory
	clusterSet           multicluster.ClusterSet
}

func NewFederatedRouteTableReconciler(
	ctx context.Context,
	federatedRouteTables fed_gateway_solo_io_v1.FederatedRouteTableClient,
	baseClients gateway_solo_io_v1.MulticlusterClientset,
	statusBuilderFactory placement.StatusBuilderFactory,
	clusterSet multicluster.ClusterSet,
) controller.FederatedRouteTableFinalizer {
	return &federatedRouteTableReconciler{
		ctx:                  ctx,
		federatedRouteTables: federatedRouteTables,
		baseClients:          baseClients,
		statusBuilderFactory: statusBuilderFactory,
		clusterSet:           clusterSet,
	}
}

func (f *federatedRouteTableReconciler) ReconcileFederatedRouteTable(obj *fed_gateway_solo_io_v1.FederatedRouteTable) (reconcile.Result, error) {
	if !obj.NeedsReconcile() {
		return reconcile.Result{}, nil
	}

	contextutils.LoggerFrom(f.ctx).Debugw("processing federated routeTable", zap.Any("FederatedRouteTable", obj))
	statusBuilder := f.statusBuilderFactory.GetBuilder()

	allClusters := f.clusterSet.ListClusters()

	// Validate resource
	if obj.Spec.GetPlacement() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.PlacementMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedRouteTables.UpdateFederatedRouteTableStatus(f.ctx, obj)
	}
	for _, cluster := range obj.Spec.Placement.GetClusters() {
		if !stringutils.ContainsString(cluster, allClusters) {
			obj.Status.PlacementStatus = statusBuilder.
				UpdateUnprocessed(obj.Status.PlacementStatus, placement.ClusterNotRegistered(cluster), mc_types.PlacementStatus_INVALID).
				Eject(obj.GetGeneration())
			return reconcile.Result{}, f.federatedRouteTables.UpdateFederatedRouteTableStatus(f.ctx, obj)
		}
	}
	if obj.Spec.Template.GetSpec() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.SpecTemplateMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedRouteTables.UpdateFederatedRouteTableStatus(f.ctx, obj)
	}
	if obj.Spec.Template.GetMetadata() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.MetaTemplateMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedRouteTables.UpdateFederatedRouteTableStatus(f.ctx, obj)
	}

	// ownerLabel is used to reference Federated resources via their children.
	ownerLabel := federation.GetOwnerLabel(obj)

	spec := obj.Spec.Template.GetSpec()
	meta := obj.Spec.Template.GetMetadata()
	labels := federation.Merge(meta.GetLabels(), ownerLabel)

	multiErr := &multierror.Error{}
	for _, cluster := range allClusters {
		clusterRouteTables := gateway_solo_io_v1_sets.NewRouteTableSet()
		if stringutils.ContainsString(cluster, obj.Spec.Placement.GetClusters()) {
			for _, namespace := range obj.Spec.Placement.GetNamespaces() {

				clusterRouteTables.Insert(&gateway_solo_io_v1.RouteTable{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:   namespace,
						Name:        meta.GetName(),
						Labels:      labels,
						Annotations: meta.GetAnnotations(),
					},
					Spec: *spec,
				})
			}
		}

		if err := f.ensureCluster(cluster, statusBuilder, clusterRouteTables, ownerLabel); err != nil {
			multiErr.Errors = append(multiErr.Errors, err)
		}
	}

	obj.Status = fed_gateway_solo_io_v1_types.FederatedRouteTableStatus{
		PlacementStatus: statusBuilder.Build(obj.GetGeneration()),
	}
	err := f.federatedRouteTables.UpdateFederatedRouteTableStatus(f.ctx, obj)
	if err != nil {
		multiErr.Errors = append(multiErr.Errors, err)
		contextutils.LoggerFrom(f.ctx).Errorw("Failed to update status on federated routeTable", zap.Error(err))
	}

	return reconcile.Result{}, multiErr.ErrorOrNil()
}

func (f *federatedRouteTableReconciler) FederatedRouteTableFinalizerName() string {
	return federation.HubFinalizer
}

func (f *federatedRouteTableReconciler) FinalizeFederatedRouteTable(obj *fed_gateway_solo_io_v1.FederatedRouteTable) error {
	return f.deleteAll(federation.GetOwnerLabel(obj))
}

// ensureCluster upserts all desired resources on the given cluster.
// An error is returned only if a retry is expected to resolve the issue.
func (f *federatedRouteTableReconciler) ensureCluster(cluster string, statusBuilder placement.StatusBuilder, desired gateway_solo_io_v1_sets.RouteTableSet, ownerLabel map[string]string) error {
	clientset, err := f.baseClients.Cluster(cluster)
	if err != nil {
		var namespaces []string
		for _, obj := range desired.List() {
			namespaces = append(namespaces, obj.GetNamespace())
		}

		contextutils.LoggerFrom(f.ctx).Errorw("Failed to get clientset", zap.String("cluster", cluster), zap.Error(err))
		statusBuilder.AddDestinations([]string{cluster}, namespaces, mc_types.PlacementStatus_Namespace{
			State:   mc_types.PlacementStatus_FAILED,
			Message: placement.FailedToCreateClientForCluster(cluster),
		})
		return err
	}

	routeTableClient := clientset.RouteTables()

	existingList, err := routeTableClient.ListRouteTable(f.ctx, client.MatchingLabels(ownerLabel))
	if err != nil {
		var namespaces []string
		for _, obj := range desired.List() {
			namespaces = append(namespaces, obj.GetNamespace())
		}

		contextutils.LoggerFrom(f.ctx).Errorw("Failed to list routeTables", zap.Error(err))
		statusBuilder.AddDestinations([]string{cluster}, namespaces, mc_types.PlacementStatus_Namespace{
			State:   mc_types.PlacementStatus_FAILED,
			Message: placement.FailedToListResource("routeTable", cluster),
		})
		return err
	}

	existing := gateway_solo_io_v1_sets.NewRouteTableSet()
	for _, routeTable := range existingList.Items {
		routeTablePointer := routeTable
		existing.Insert(&routeTablePointer)
	}

	multiErr := &multierror.Error{}
	for _, desiredRouteTable := range desired.List() {
		err := routeTableClient.UpsertRouteTable(f.ctx, desiredRouteTable)
		if err != nil && errors.IsConflict(err) {
			multiErr.Errors = append(multiErr.Errors, err)
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to upsert routeTable due to resource conflict", zap.Error(err))
			statusBuilder.AddDestination(cluster, desiredRouteTable.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_FAILED,
				Message: placement.FailedToUpsertResourceDueToConflict("routeTable"),
			})
		} else if err != nil {
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to upsert routeTable", zap.Error(err))
			statusBuilder.AddDestination(cluster, desiredRouteTable.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_FAILED,
				Message: placement.FailedToUpsertResource("routeTable"),
			})
		} else {
			statusBuilder.AddDestination(cluster, desiredRouteTable.Namespace, mc_types.PlacementStatus_Namespace{
				State: mc_types.PlacementStatus_PLACED,
			})
		}
	}

	for _, staleRouteTable := range existing.Difference(desired).List() {
		err := routeTableClient.DeleteRouteTable(f.ctx, client.ObjectKey{
			Namespace: staleRouteTable.Namespace,
			Name:      staleRouteTable.Name,
		})
		if client.IgnoreNotFound(err) != nil {
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to delete routeTable", zap.Error(err))
			statusBuilder.AddDestination(cluster, staleRouteTable.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_STALE,
				Message: placement.FailedToDeleteResource("routeTable"),
			})
		}
	}

	return multiErr.ErrorOrNil()
}

// Delete all routeTables in the matching the ownerLabel on all clusters.
// Used to ensure that routeTables generated by a FederatedRouteTable are cleaned up on delete.
func (f *federatedRouteTableReconciler) deleteAll(ownerLabel map[string]string) error {
	for _, cluster := range f.clusterSet.ListClusters() {
		clusterClient, err := f.baseClients.Cluster(cluster)
		if err != nil {
			return err
		}
		// TODO this requires permissions in all namespaces, we could restrict to to namespaces referenced by gloo instances
		list, err := clusterClient.RouteTables().ListRouteTable(f.ctx, client.MatchingLabels(ownerLabel))
		if err != nil {
			return err
		}

		for _, routeTable := range list.Items {
			err = clusterClient.RouteTables().DeleteRouteTable(f.ctx, client.ObjectKey{
				Namespace: routeTable.Namespace,
				Name:      routeTable.Name,
			})
			if client.IgnoreNotFound(err) != nil {
				return err
			}
		}
	}
	return nil
}
