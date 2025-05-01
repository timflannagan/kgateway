// pkg/pluginsdk/context.go
//
// Public context structs and lightweight helpers that plugin authors use
// inside Apply* methods.  They contain **no** internal logic, so exposing
// them does not leak implementation details.

package pluginsdk

import (
	"google.golang.org/protobuf/proto"
)

// -----------------------------------------------------------------------------
// Convenience map for per-route / per-cluster typed-config
// -----------------------------------------------------------------------------

type TypedFilterConfigMap map[string]proto.Message

// AddTypedConfig replaces (or sets) the proto for key.
func (m *TypedFilterConfigMap) AddTypedConfig(key string, v proto.Message) {
	if *m == nil {
		*m = make(TypedFilterConfigMap)
	}
	(*m)[key] = v
}

// GetTypedConfig returns the proto for key, nil if not present.
func (m *TypedFilterConfigMap) GetTypedConfig(key string) proto.Message {
	if m == nil || *m == nil {
		return nil
	}
	return (*m)[key]
}

// -----------------------------------------------------------------------------
// Per-stage context structs used by ProxyTranslationPass hooks
// -----------------------------------------------------------------------------

type ListenerContext struct {
	Policy PolicyIR
}

type RouteConfigContext struct{} // no per-policy data

type VirtualHostContext struct {
	Policy PolicyIR
}

type RouteBackendContext struct {
	FilterChainName   string
	Backend           *BackendObjectIR
	TypedFilterConfig TypedFilterConfigMap
}

type RouteContext struct {
	FilterChainName   string
	Policy            PolicyIR
	In                HttpRouteRuleMatchIR
	TypedFilterConfig TypedFilterConfigMap
}

type HcmContext struct {
	Policy PolicyIR
}
