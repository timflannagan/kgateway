// Code generated by skv2. DO NOT EDIT.

package v1

import (
	mc_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
)

func (obj *FederatedUpstream) NeedsReconcile() bool {
	// If the FederatedUpstream has not been observed or is in some state where a retry is needed, it needs reconcile
	return obj.Generation != obj.Status.PlacementStatus.GetObservedGeneration() ||
		obj.Status.PlacementStatus.GetState() != mc_types.PlacementStatus_PLACED ||
		obj.Status.PlacementStatus.GetState() != mc_types.PlacementStatus_FAILED
}

func (obj *FederatedUpstreamGroup) NeedsReconcile() bool {
	// If the FederatedUpstreamGroup has not been observed or is in some state where a retry is needed, it needs reconcile
	return obj.Generation != obj.Status.PlacementStatus.GetObservedGeneration() ||
		obj.Status.PlacementStatus.GetState() != mc_types.PlacementStatus_PLACED ||
		obj.Status.PlacementStatus.GetState() != mc_types.PlacementStatus_FAILED
}

func (obj *FederatedSettings) NeedsReconcile() bool {
	// If the FederatedSettings has not been observed or is in some state where a retry is needed, it needs reconcile
	return obj.Generation != obj.Status.PlacementStatus.GetObservedGeneration() ||
		obj.Status.PlacementStatus.GetState() != mc_types.PlacementStatus_PLACED ||
		obj.Status.PlacementStatus.GetState() != mc_types.PlacementStatus_FAILED
}
