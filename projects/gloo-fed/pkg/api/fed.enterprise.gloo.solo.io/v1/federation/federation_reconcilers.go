// Code generated by skv2. DO NOT EDIT.

// Definition for federated resource reconciler templates.
package federation

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/stringutils"
	"github.com/solo-io/skv2/pkg/reconcile"
	enterprise_gloo_solo_io_v1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	enterprise_gloo_solo_io_v1_sets "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1/sets"
	fed_enterprise_gloo_solo_io_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1/controller"
	fed_enterprise_gloo_solo_io_v1_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1/types"
	mc_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation/placement"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type federatedAuthConfigReconciler struct {
	ctx                  context.Context
	federatedAuthConfigs fed_enterprise_gloo_solo_io_v1.FederatedAuthConfigClient
	baseClients          enterprise_gloo_solo_io_v1.MulticlusterClientset
	statusBuilderFactory placement.StatusBuilderFactory
	clusterSet           multicluster.ClusterSet
}

func NewFederatedAuthConfigReconciler(
	ctx context.Context,
	federatedAuthConfigs fed_enterprise_gloo_solo_io_v1.FederatedAuthConfigClient,
	baseClients enterprise_gloo_solo_io_v1.MulticlusterClientset,
	statusBuilderFactory placement.StatusBuilderFactory,
	clusterSet multicluster.ClusterSet,
) controller.FederatedAuthConfigFinalizer {
	return &federatedAuthConfigReconciler{
		ctx:                  ctx,
		federatedAuthConfigs: federatedAuthConfigs,
		baseClients:          baseClients,
		statusBuilderFactory: statusBuilderFactory,
		clusterSet:           clusterSet,
	}
}

func (f *federatedAuthConfigReconciler) ReconcileFederatedAuthConfig(obj *fed_enterprise_gloo_solo_io_v1.FederatedAuthConfig) (reconcile.Result, error) {
	if !obj.NeedsReconcile() {
		return reconcile.Result{}, nil
	}

	contextutils.LoggerFrom(f.ctx).Debugw("processing federated authConfig", zap.Any("FederatedAuthConfig", obj))
	statusBuilder := f.statusBuilderFactory.GetBuilder()

	allClusters := f.clusterSet.ListClusters()

	// Validate resource
	if obj.Spec.GetPlacement() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.PlacementMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedAuthConfigs.UpdateFederatedAuthConfigStatus(f.ctx, obj)
	}
	for _, cluster := range obj.Spec.Placement.GetClusters() {
		if !stringutils.ContainsString(cluster, allClusters) {
			obj.Status.PlacementStatus = statusBuilder.
				UpdateUnprocessed(obj.Status.PlacementStatus, placement.ClusterNotRegistered(cluster), mc_types.PlacementStatus_INVALID).
				Eject(obj.GetGeneration())
			return reconcile.Result{}, f.federatedAuthConfigs.UpdateFederatedAuthConfigStatus(f.ctx, obj)
		}
	}
	if obj.Spec.Template.GetSpec() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.SpecTemplateMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedAuthConfigs.UpdateFederatedAuthConfigStatus(f.ctx, obj)
	}
	if obj.Spec.Template.GetMetadata() == nil {
		obj.Status.PlacementStatus = statusBuilder.
			UpdateUnprocessed(obj.Status.PlacementStatus, placement.MetaTemplateMissing, mc_types.PlacementStatus_INVALID).
			Eject(obj.GetGeneration())
		return reconcile.Result{}, f.federatedAuthConfigs.UpdateFederatedAuthConfigStatus(f.ctx, obj)
	}

	// ownerLabel is used to reference Federated resources via their children.
	ownerLabel := federation.GetOwnerLabel(obj)

	spec := obj.Spec.Template.GetSpec()
	meta := obj.Spec.Template.GetMetadata()
	labels := federation.Merge(meta.GetLabels(), ownerLabel)

	multiErr := &multierror.Error{}
	for _, cluster := range allClusters {
		clusterAuthConfigs := enterprise_gloo_solo_io_v1_sets.NewAuthConfigSet()
		if stringutils.ContainsString(cluster, obj.Spec.Placement.GetClusters()) {
			for _, namespace := range obj.Spec.Placement.GetNamespaces() {

				clusterAuthConfigs.Insert(&enterprise_gloo_solo_io_v1.AuthConfig{
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

		if err := f.ensureCluster(cluster, statusBuilder, clusterAuthConfigs, ownerLabel); err != nil {
			multiErr.Errors = append(multiErr.Errors, err)
		}
	}

	obj.Status = fed_enterprise_gloo_solo_io_v1_types.FederatedAuthConfigStatus{
		PlacementStatus: statusBuilder.Build(obj.GetGeneration()),
	}
	err := f.federatedAuthConfigs.UpdateFederatedAuthConfigStatus(f.ctx, obj)
	if err != nil {
		multiErr.Errors = append(multiErr.Errors, err)
		contextutils.LoggerFrom(f.ctx).Errorw("Failed to update status on federated authConfig", zap.Error(err))
	}

	return reconcile.Result{}, multiErr.ErrorOrNil()
}

func (f *federatedAuthConfigReconciler) FederatedAuthConfigFinalizerName() string {
	return federation.HubFinalizer
}

func (f *federatedAuthConfigReconciler) FinalizeFederatedAuthConfig(obj *fed_enterprise_gloo_solo_io_v1.FederatedAuthConfig) error {
	return f.deleteAll(federation.GetOwnerLabel(obj))
}

// ensureCluster upserts all desired resources on the given cluster.
// An error is returned only if a retry is expected to resolve the issue.
func (f *federatedAuthConfigReconciler) ensureCluster(cluster string, statusBuilder placement.StatusBuilder, desired enterprise_gloo_solo_io_v1_sets.AuthConfigSet, ownerLabel map[string]string) error {
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

	authConfigClient := clientset.AuthConfigs()

	existingList, err := authConfigClient.ListAuthConfig(f.ctx, client.MatchingLabels(ownerLabel))
	if err != nil {
		var namespaces []string
		for _, obj := range desired.List() {
			namespaces = append(namespaces, obj.GetNamespace())
		}

		contextutils.LoggerFrom(f.ctx).Errorw("Failed to list authConfigs", zap.Error(err))
		statusBuilder.AddDestinations([]string{cluster}, namespaces, mc_types.PlacementStatus_Namespace{
			State:   mc_types.PlacementStatus_FAILED,
			Message: placement.FailedToListResource("authConfig", cluster),
		})
		return err
	}

	existing := enterprise_gloo_solo_io_v1_sets.NewAuthConfigSet()
	for _, authConfig := range existingList.Items {
		authConfigPointer := authConfig
		existing.Insert(&authConfigPointer)
	}

	multiErr := &multierror.Error{}
	for _, desiredAuthConfig := range desired.List() {
		err := authConfigClient.UpsertAuthConfig(f.ctx, desiredAuthConfig)
		if err != nil && errors.IsConflict(err) {
			multiErr.Errors = append(multiErr.Errors, err)
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to upsert authConfig due to resource conflict", zap.Error(err))
			statusBuilder.AddDestination(cluster, desiredAuthConfig.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_FAILED,
				Message: placement.FailedToUpsertResourceDueToConflict("authConfig"),
			})
		} else if err != nil {
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to upsert authConfig", zap.Error(err))
			statusBuilder.AddDestination(cluster, desiredAuthConfig.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_FAILED,
				Message: placement.FailedToUpsertResource("authConfig"),
			})
		} else {
			statusBuilder.AddDestination(cluster, desiredAuthConfig.Namespace, mc_types.PlacementStatus_Namespace{
				State: mc_types.PlacementStatus_PLACED,
			})
		}
	}

	for _, staleAuthConfig := range existing.Difference(desired).List() {
		err := authConfigClient.DeleteAuthConfig(f.ctx, client.ObjectKey{
			Namespace: staleAuthConfig.Namespace,
			Name:      staleAuthConfig.Name,
		})
		if client.IgnoreNotFound(err) != nil {
			contextutils.LoggerFrom(f.ctx).Errorw("Failed to delete authConfig", zap.Error(err))
			statusBuilder.AddDestination(cluster, staleAuthConfig.Namespace, mc_types.PlacementStatus_Namespace{
				State:   mc_types.PlacementStatus_STALE,
				Message: placement.FailedToDeleteResource("authConfig"),
			})
		}
	}

	return multiErr.ErrorOrNil()
}

// Delete all authConfigs in the matching the ownerLabel on all clusters.
// Used to ensure that authConfigs generated by a FederatedAuthConfig are cleaned up on delete.
func (f *federatedAuthConfigReconciler) deleteAll(ownerLabel map[string]string) error {
	for _, cluster := range f.clusterSet.ListClusters() {
		clusterClient, err := f.baseClients.Cluster(cluster)
		if err != nil {
			return err
		}
		// TODO this requires permissions in all namespaces, we could restrict to to namespaces referenced by gloo instances
		list, err := clusterClient.AuthConfigs().ListAuthConfig(f.ctx, client.MatchingLabels(ownerLabel))
		if err != nil {
			return err
		}

		for _, authConfig := range list.Items {
			err = clusterClient.AuthConfigs().DeleteAuthConfig(f.ctx, client.ObjectKey{
				Namespace: authConfig.Namespace,
				Name:      authConfig.Name,
			})
			if client.IgnoreNotFound(err) != nil {
				return err
			}
		}
	}
	return nil
}
