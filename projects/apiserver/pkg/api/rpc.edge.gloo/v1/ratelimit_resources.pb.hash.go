// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/ratelimit_resources.proto

package v1

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/mitchellh/hashstructure"
	safe_hasher "github.com/solo-io/protoc-gen-ext/pkg/hasher"
)

// ensure the imports are used
var (
	_ = errors.New("")
	_ = fmt.Print
	_ = binary.LittleEndian
	_ = new(hash.Hash64)
	_ = fnv.New64
	_ = hashstructure.Hash
	_ = new(safe_hasher.SafeHasher)
)

// Hash function
func (m *RateLimitConfig) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("rpc.edge.gloo.solo.io.github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1.RateLimitConfig")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetMetadata()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Metadata")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetMetadata(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Metadata")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if h, ok := interface{}(m.GetSpec()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Spec")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetSpec(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Spec")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if h, ok := interface{}(m.GetStatus()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Status")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetStatus(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Status")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if h, ok := interface{}(m.GetGlooInstance()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("GlooInstance")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetGlooInstance(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("GlooInstance")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// Hash function
func (m *ListRateLimitConfigsRequest) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("rpc.edge.gloo.solo.io.github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1.ListRateLimitConfigsRequest")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetGlooInstanceRef()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("GlooInstanceRef")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetGlooInstanceRef(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("GlooInstanceRef")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if h, ok := interface{}(m.GetPagination()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("Pagination")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetPagination(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("Pagination")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if _, err = hasher.Write([]byte(m.GetQueryString())); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetStatusFilter()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("StatusFilter")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetStatusFilter(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("StatusFilter")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	if h, ok := interface{}(m.GetSortOptions()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("SortOptions")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetSortOptions(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("SortOptions")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// Hash function
func (m *ListRateLimitConfigsResponse) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("rpc.edge.gloo.solo.io.github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1.ListRateLimitConfigsResponse")); err != nil {
		return 0, err
	}

	for _, v := range m.GetRateLimitConfigs() {

		if h, ok := interface{}(v).(safe_hasher.SafeHasher); ok {
			if _, err = hasher.Write([]byte("")); err != nil {
				return 0, err
			}
			if _, err = h.Hash(hasher); err != nil {
				return 0, err
			}
		} else {
			if fieldValue, err := hashstructure.Hash(v, nil); err != nil {
				return 0, err
			} else {
				if _, err = hasher.Write([]byte("")); err != nil {
					return 0, err
				}
				if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
					return 0, err
				}
			}
		}

	}

	err = binary.Write(hasher, binary.LittleEndian, m.GetTotal())
	if err != nil {
		return 0, err
	}

	return hasher.Sum64(), nil
}

// Hash function
func (m *GetRateLimitConfigYamlRequest) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("rpc.edge.gloo.solo.io.github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1.GetRateLimitConfigYamlRequest")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetRateLimitConfigRef()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("RateLimitConfigRef")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetRateLimitConfigRef(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("RateLimitConfigRef")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// Hash function
func (m *GetRateLimitConfigYamlResponse) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("rpc.edge.gloo.solo.io.github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1.GetRateLimitConfigYamlResponse")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetYamlData()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("YamlData")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetYamlData(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("YamlData")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// Hash function
func (m *GetRateLimitConfigDetailsRequest) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("rpc.edge.gloo.solo.io.github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1.GetRateLimitConfigDetailsRequest")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetRateLimitConfigRef()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("RateLimitConfigRef")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetRateLimitConfigRef(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("RateLimitConfigRef")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}

// Hash function
func (m *GetRateLimitConfigDetailsResponse) Hash(hasher hash.Hash64) (uint64, error) {
	if m == nil {
		return 0, nil
	}
	if hasher == nil {
		hasher = fnv.New64()
	}
	var err error
	if _, err = hasher.Write([]byte("rpc.edge.gloo.solo.io.github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1.GetRateLimitConfigDetailsResponse")); err != nil {
		return 0, err
	}

	if h, ok := interface{}(m.GetRateLimitConfig()).(safe_hasher.SafeHasher); ok {
		if _, err = hasher.Write([]byte("RateLimitConfig")); err != nil {
			return 0, err
		}
		if _, err = h.Hash(hasher); err != nil {
			return 0, err
		}
	} else {
		if fieldValue, err := hashstructure.Hash(m.GetRateLimitConfig(), nil); err != nil {
			return 0, err
		} else {
			if _, err = hasher.Write([]byte("RateLimitConfig")); err != nil {
				return 0, err
			}
			if err := binary.Write(hasher, binary.LittleEndian, fieldValue); err != nil {
				return 0, err
			}
		}
	}

	return hasher.Sum64(), nil
}
