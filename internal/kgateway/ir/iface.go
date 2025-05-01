// internal/kgateway/ir/iface.go
//
// Internal-facing aliases and helper structs that allow the core codebase
// to keep compiling while the authoritative interface & context definitions
// now live in the public `pkg/pluginsdk` package.
//
// ⚠️  Do **not** add new logic here; treat this file as a shim layer that
//     forwards to the SDK types.  All downstream plugins should import
//     github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk directly.

package ir

import (
	"context"
	"fmt"
	"time"

	envoy_cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	intplugins "github.com/kgateway-dev/kgateway/v2/internal/kgateway/plugins"
	sdk "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk"
)

/* -------------------------------------------------------------------------
   Public-SDK type aliases
   --------------------------------------------------------------------- */

type (
	ListenerContext     = sdk.ListenerContext
	RouteConfigContext  = sdk.RouteConfigContext
	VirtualHostContext  = sdk.VirtualHostContext
	RouteBackendContext = sdk.RouteBackendContext
	RouteContext        = sdk.RouteContext
	HcmContext          = sdk.HcmContext

	TypedFilterConfigMap = sdk.TypedFilterConfigMap

	ProxyTranslationPass              = sdk.ProxyTranslationPass
	UnimplementedProxyTranslationPass = sdk.UnimplementedProxyTranslationPass
)

/* -------------------------------------------------------------------------
   Local helpers & wrappers (unchanged from previous implementation)
   --------------------------------------------------------------------- */

type Resources struct {
	Clusters []*envoy_cluster.Cluster
}

// Kept as a plain struct; not part of the public SDK.
type GwTranslationCtx struct{}

/* ---------- Policy IR plumbing ------------------------------------- */

type PolicyIR interface {
	CreationTime() time.Time
	Equals(other any) bool
}

type PolicyWrapper struct {
	// Original policy resource identity
	ObjectSource `json:",inline"`

	// Raw K8s object
	Policy metav1.Object

	// Pure-policy errors (prior to attachment validation)
	Errors []error

	// Fully-translated in-memory representation
	PolicyIR PolicyIR

	// Where the policy should attach
	TargetRefs []PolicyRef
}

func (w PolicyWrapper) ResourceName() string { return w.ObjectSource.ResourceName() }

func (w PolicyWrapper) Equals(other PolicyWrapper) bool {
	if w.ObjectSource != other.ObjectSource {
		return false
	}
	return versionEquals(w.Policy, other.Policy) && w.PolicyIR.Equals(other.PolicyIR)
}

func versionEquals(a, b metav1.Object) bool {
	genEqual := a.GetGeneration() != 0 && b.GetGeneration() != 0 &&
		a.GetGeneration() == b.GetGeneration()
	resEqual := a.GetResourceVersion() == b.GetResourceVersion()
	return (genEqual || resEqual) && a.GetUID() == b.GetUID()
}

var ErrNotAttachable = fmt.Errorf("policy is not attachable to this object")

/* ---------- Execution-time hooks for plugin collections ------------ */

type PolicyRun interface {
	NewGatewayTranslationPass(ctx context.Context, tctx GwTranslationCtx) ProxyTranslationPass
	ProcessBackend(ctx context.Context, in BackendObjectIR, out *envoy_cluster.Cluster) error
}

/* ---------- Opaque aliases used elsewhere in internal packages ----- */

type (
	BackendObjectIR      = any
	HttpBackend          = any
	FilterChainCommon    = any
	StagedHttpFilter     = intplugins.StagedHttpFilter
	StagedNetworkFilter  = intplugins.StagedNetworkFilter
	HttpRouteRuleMatchIR = any
)
