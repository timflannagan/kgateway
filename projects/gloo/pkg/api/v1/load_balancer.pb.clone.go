// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/v1/load_balancer.proto

package v1

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/solo-io/protoc-gen-ext/pkg/clone"
	"google.golang.org/protobuf/proto"

	github_com_golang_protobuf_ptypes_duration "github.com/golang/protobuf/ptypes/duration"

	github_com_golang_protobuf_ptypes_empty "github.com/golang/protobuf/ptypes/empty"

	github_com_golang_protobuf_ptypes_wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

// ensure the imports are used
var (
	_ = errors.New("")
	_ = fmt.Print
	_ = binary.LittleEndian
	_ = bytes.Compare
	_ = strings.Compare
	_ = clone.Cloner(nil)
	_ = proto.Message(nil)
)

// Clone function
func (m *LoadBalancerConfig) Clone() proto.Message {
	var target *LoadBalancerConfig
	if m == nil {
		return target
	}
	target = &LoadBalancerConfig{}

	if h, ok := interface{}(m.GetHealthyPanicThreshold()).(clone.Cloner); ok {
		target.HealthyPanicThreshold = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.DoubleValue)
	} else {
		target.HealthyPanicThreshold = proto.Clone(m.GetHealthyPanicThreshold()).(*github_com_golang_protobuf_ptypes_wrappers.DoubleValue)
	}

	if h, ok := interface{}(m.GetUpdateMergeWindow()).(clone.Cloner); ok {
		target.UpdateMergeWindow = h.Clone().(*github_com_golang_protobuf_ptypes_duration.Duration)
	} else {
		target.UpdateMergeWindow = proto.Clone(m.GetUpdateMergeWindow()).(*github_com_golang_protobuf_ptypes_duration.Duration)
	}

	if h, ok := interface{}(m.GetUseHostnameForHashing()).(clone.Cloner); ok {
		target.UseHostnameForHashing = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.BoolValue)
	} else {
		target.UseHostnameForHashing = proto.Clone(m.GetUseHostnameForHashing()).(*github_com_golang_protobuf_ptypes_wrappers.BoolValue)
	}

	switch m.Type.(type) {

	case *LoadBalancerConfig_RoundRobin_:

		if h, ok := interface{}(m.GetRoundRobin()).(clone.Cloner); ok {
			target.Type = &LoadBalancerConfig_RoundRobin_{
				RoundRobin: h.Clone().(*LoadBalancerConfig_RoundRobin),
			}
		} else {
			target.Type = &LoadBalancerConfig_RoundRobin_{
				RoundRobin: proto.Clone(m.GetRoundRobin()).(*LoadBalancerConfig_RoundRobin),
			}
		}

	case *LoadBalancerConfig_LeastRequest_:

		if h, ok := interface{}(m.GetLeastRequest()).(clone.Cloner); ok {
			target.Type = &LoadBalancerConfig_LeastRequest_{
				LeastRequest: h.Clone().(*LoadBalancerConfig_LeastRequest),
			}
		} else {
			target.Type = &LoadBalancerConfig_LeastRequest_{
				LeastRequest: proto.Clone(m.GetLeastRequest()).(*LoadBalancerConfig_LeastRequest),
			}
		}

	case *LoadBalancerConfig_Random_:

		if h, ok := interface{}(m.GetRandom()).(clone.Cloner); ok {
			target.Type = &LoadBalancerConfig_Random_{
				Random: h.Clone().(*LoadBalancerConfig_Random),
			}
		} else {
			target.Type = &LoadBalancerConfig_Random_{
				Random: proto.Clone(m.GetRandom()).(*LoadBalancerConfig_Random),
			}
		}

	case *LoadBalancerConfig_RingHash_:

		if h, ok := interface{}(m.GetRingHash()).(clone.Cloner); ok {
			target.Type = &LoadBalancerConfig_RingHash_{
				RingHash: h.Clone().(*LoadBalancerConfig_RingHash),
			}
		} else {
			target.Type = &LoadBalancerConfig_RingHash_{
				RingHash: proto.Clone(m.GetRingHash()).(*LoadBalancerConfig_RingHash),
			}
		}

	case *LoadBalancerConfig_Maglev_:

		if h, ok := interface{}(m.GetMaglev()).(clone.Cloner); ok {
			target.Type = &LoadBalancerConfig_Maglev_{
				Maglev: h.Clone().(*LoadBalancerConfig_Maglev),
			}
		} else {
			target.Type = &LoadBalancerConfig_Maglev_{
				Maglev: proto.Clone(m.GetMaglev()).(*LoadBalancerConfig_Maglev),
			}
		}

	}

	switch m.LocalityConfig.(type) {

	case *LoadBalancerConfig_LocalityWeightedLbConfig:

		if h, ok := interface{}(m.GetLocalityWeightedLbConfig()).(clone.Cloner); ok {
			target.LocalityConfig = &LoadBalancerConfig_LocalityWeightedLbConfig{
				LocalityWeightedLbConfig: h.Clone().(*github_com_golang_protobuf_ptypes_empty.Empty),
			}
		} else {
			target.LocalityConfig = &LoadBalancerConfig_LocalityWeightedLbConfig{
				LocalityWeightedLbConfig: proto.Clone(m.GetLocalityWeightedLbConfig()).(*github_com_golang_protobuf_ptypes_empty.Empty),
			}
		}

	}

	return target
}

// Clone function
func (m *LoadBalancerConfig_RoundRobin) Clone() proto.Message {
	var target *LoadBalancerConfig_RoundRobin
	if m == nil {
		return target
	}
	target = &LoadBalancerConfig_RoundRobin{}

	if h, ok := interface{}(m.GetSlowStartConfig()).(clone.Cloner); ok {
		target.SlowStartConfig = h.Clone().(*LoadBalancerConfig_SlowStartConfig)
	} else {
		target.SlowStartConfig = proto.Clone(m.GetSlowStartConfig()).(*LoadBalancerConfig_SlowStartConfig)
	}

	return target
}

// Clone function
func (m *LoadBalancerConfig_LeastRequest) Clone() proto.Message {
	var target *LoadBalancerConfig_LeastRequest
	if m == nil {
		return target
	}
	target = &LoadBalancerConfig_LeastRequest{}

	target.ChoiceCount = m.GetChoiceCount()

	if h, ok := interface{}(m.GetSlowStartConfig()).(clone.Cloner); ok {
		target.SlowStartConfig = h.Clone().(*LoadBalancerConfig_SlowStartConfig)
	} else {
		target.SlowStartConfig = proto.Clone(m.GetSlowStartConfig()).(*LoadBalancerConfig_SlowStartConfig)
	}

	return target
}

// Clone function
func (m *LoadBalancerConfig_Random) Clone() proto.Message {
	var target *LoadBalancerConfig_Random
	if m == nil {
		return target
	}
	target = &LoadBalancerConfig_Random{}

	return target
}

// Clone function
func (m *LoadBalancerConfig_RingHashConfig) Clone() proto.Message {
	var target *LoadBalancerConfig_RingHashConfig
	if m == nil {
		return target
	}
	target = &LoadBalancerConfig_RingHashConfig{}

	target.MinimumRingSize = m.GetMinimumRingSize()

	target.MaximumRingSize = m.GetMaximumRingSize()

	return target
}

// Clone function
func (m *LoadBalancerConfig_RingHash) Clone() proto.Message {
	var target *LoadBalancerConfig_RingHash
	if m == nil {
		return target
	}
	target = &LoadBalancerConfig_RingHash{}

	if h, ok := interface{}(m.GetRingHashConfig()).(clone.Cloner); ok {
		target.RingHashConfig = h.Clone().(*LoadBalancerConfig_RingHashConfig)
	} else {
		target.RingHashConfig = proto.Clone(m.GetRingHashConfig()).(*LoadBalancerConfig_RingHashConfig)
	}

	return target
}

// Clone function
func (m *LoadBalancerConfig_Maglev) Clone() proto.Message {
	var target *LoadBalancerConfig_Maglev
	if m == nil {
		return target
	}
	target = &LoadBalancerConfig_Maglev{}

	return target
}

// Clone function
func (m *LoadBalancerConfig_SlowStartConfig) Clone() proto.Message {
	var target *LoadBalancerConfig_SlowStartConfig
	if m == nil {
		return target
	}
	target = &LoadBalancerConfig_SlowStartConfig{}

	if h, ok := interface{}(m.GetSlowStartWindow()).(clone.Cloner); ok {
		target.SlowStartWindow = h.Clone().(*github_com_golang_protobuf_ptypes_duration.Duration)
	} else {
		target.SlowStartWindow = proto.Clone(m.GetSlowStartWindow()).(*github_com_golang_protobuf_ptypes_duration.Duration)
	}

	if h, ok := interface{}(m.GetAggression()).(clone.Cloner); ok {
		target.Aggression = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.DoubleValue)
	} else {
		target.Aggression = proto.Clone(m.GetAggression()).(*github_com_golang_protobuf_ptypes_wrappers.DoubleValue)
	}

	if h, ok := interface{}(m.GetMinWeightPercent()).(clone.Cloner); ok {
		target.MinWeightPercent = h.Clone().(*github_com_golang_protobuf_ptypes_wrappers.DoubleValue)
	} else {
		target.MinWeightPercent = proto.Clone(m.GetMinWeightPercent()).(*github_com_golang_protobuf_ptypes_wrappers.DoubleValue)
	}

	return target
}
