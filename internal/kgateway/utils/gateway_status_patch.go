package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilretry "k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	api "sigs.k8s.io/gateway-api/apis/v1"
)

// GatewayStatusPatcher provides strategic field-level patching for Gateway status
// to avoid race conditions between different components updating different fields.
type GatewayStatusPatcher struct {
	client client.Client
}

// NewGatewayStatusPatcher creates a new strategic status patcher
func NewGatewayStatusPatcher(client client.Client) *GatewayStatusPatcher {
	return &GatewayStatusPatcher{client: client}
}

// PatchAddresses updates only the addresses field of a Gateway status
// This is used by the Gateway Controller for infrastructure-related address updates
func (p *GatewayStatusPatcher) PatchAddresses(
	ctx context.Context,
	gwNN types.NamespacedName,
	addresses []api.GatewayStatusAddress,
) error {
	return utilretry.RetryOnConflict(utilretry.DefaultRetry, func() error {
		// Create strategic patch for addresses only
		patch := map[string]interface{}{
			"status": map[string]interface{}{
				"addresses": addresses,
			},
		}

		patchBytes, err := json.Marshal(patch)
		if err != nil {
			return fmt.Errorf("failed to marshal addresses patch: %w", err)
		}

		var gw api.Gateway
		gw.SetNamespace(gwNN.Namespace)
		gw.SetName(gwNN.Name)
		return p.client.Status().Patch(ctx, &gw, client.RawPatch(types.MergePatchType, patchBytes))
	})
}

// PatchCondition updates only a specific condition in the Gateway status
// This is used by the Status Syncer for translation/xDS-related condition updates
func (p *GatewayStatusPatcher) PatchCondition(
	ctx context.Context,
	gwNN types.NamespacedName,
	condition metav1.Condition,
) error {
	return utilretry.RetryOnConflict(utilretry.DefaultRetry, func() error {
		// First, get the current Gateway to read existing conditions
		var gw api.Gateway
		if err := p.client.Get(ctx, gwNN, &gw); err != nil {
			return err
		}

		// Update conditions in memory
		meta.SetStatusCondition(&gw.Status.Conditions, condition)

		// Create strategic patch for conditions only
		patch := map[string]interface{}{
			"status": map[string]interface{}{
				"conditions": gw.Status.Conditions,
			},
		}

		patchBytes, err := json.Marshal(patch)
		if err != nil {
			return fmt.Errorf("failed to marshal condition patch: %w", err)
		}

		return p.client.Status().Patch(ctx, &gw, client.RawPatch(types.MergePatchType, patchBytes))
	})
}

// PatchAcceptedConditionAndListeners updates only the Accepted condition and listener status
// This is used by the Status Syncer when updating Gateway status from translation results
// It specifically preserves the Programmed condition which is owned by Gateway Controller
func (p *GatewayStatusPatcher) PatchAcceptedConditionAndListeners(
	ctx context.Context,
	gwNN types.NamespacedName,
	newStatus *api.GatewayStatus,
) error {
	return utilretry.RetryOnConflict(utilretry.DefaultRetry, func() error {
		// First, get the current Gateway to read existing conditions
		var gw api.Gateway
		if err := p.client.Get(ctx, gwNN, &gw); err != nil {
			return err
		}

		// Only patch conditions that the Status Syncer actually built.
		// The newStatus should only contain conditions the Status Syncer owns.
		var ownedConditions []metav1.Condition
		for _, condition := range newStatus.Conditions {
			// Double-check: only patch conditions the Status Syncer should manage
			if condition.Type == string(api.GatewayConditionAccepted) ||
			   condition.Type == "AttachedListenerSets" {
				ownedConditions = append(ownedConditions, condition)
			}
		}

		// Build JSON patches for each owned condition and listeners
		var patches []map[string]interface{}

		// Add patches for each condition we own
		for _, condition := range ownedConditions {
			// Find the index of this condition type in the current conditions
			conditionIndex := -1
			for i, existingCond := range gw.Status.Conditions {
				if existingCond.Type == condition.Type {
					conditionIndex = i
					break
				}
			}

			var patch map[string]interface{}
			if conditionIndex >= 0 {
				// Update existing condition
				patch = map[string]interface{}{
					"op":    "replace",
					"path":  fmt.Sprintf("/status/conditions/%d", conditionIndex),
					"value": condition,
				}
			} else {
				// Add new condition
				patch = map[string]interface{}{
					"op":    "add",
					"path":  "/status/conditions/-",
					"value": condition,
				}
			}
			patches = append(patches, patch)
		}

		// Add patch for listeners
		patches = append(patches, map[string]interface{}{
			"op":    "replace",
			"path":  "/status/listeners",
			"value": newStatus.Listeners,
		})

		patchBytes, err := json.Marshal(patches)
		if err != nil {
			return fmt.Errorf("failed to marshal status syncer patch: %w", err)
		}

		return p.client.Status().Patch(ctx, &gw, client.RawPatch(types.JSONPatchType, patchBytes))
	})
}

// PatchFullStatus performs a full status update using strategic patching
// This preserves the existing behavior but with better conflict resolution
func (p *GatewayStatusPatcher) PatchFullStatus(
	ctx context.Context,
	gwNN types.NamespacedName,
	updateFunc func(*api.Gateway) bool,
) error {
	return utilretry.RetryOnConflict(utilretry.DefaultRetry, func() error {
		var gw api.Gateway
		if err := p.client.Get(ctx, gwNN, &gw); err != nil {
			return err
		}

		original := gw.DeepCopy()
		if !updateFunc(&gw) {
			return nil // No update needed
		}

		return p.client.Status().Patch(ctx, &gw, client.MergeFrom(original))
	})
}