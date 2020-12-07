// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache (interfaces: SnapshotCache)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	envoy_service_discovery_v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	gomock "github.com/golang/mock/gomock"
	cache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
)

// MockSnapshotCache is a mock of SnapshotCache interface
type MockSnapshotCache struct {
	ctrl     *gomock.Controller
	recorder *MockSnapshotCacheMockRecorder
}

// MockSnapshotCacheMockRecorder is the mock recorder for MockSnapshotCache
type MockSnapshotCacheMockRecorder struct {
	mock *MockSnapshotCache
}

// NewMockSnapshotCache creates a new mock instance
func NewMockSnapshotCache(ctrl *gomock.Controller) *MockSnapshotCache {
	mock := &MockSnapshotCache{ctrl: ctrl}
	mock.recorder = &MockSnapshotCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSnapshotCache) EXPECT() *MockSnapshotCacheMockRecorder {
	return m.recorder
}

// ClearSnapshot mocks base method
func (m *MockSnapshotCache) ClearSnapshot(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearSnapshot", arg0)
}

// ClearSnapshot indicates an expected call of ClearSnapshot
func (mr *MockSnapshotCacheMockRecorder) ClearSnapshot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearSnapshot", reflect.TypeOf((*MockSnapshotCache)(nil).ClearSnapshot), arg0)
}

// CreateWatch mocks base method
func (m *MockSnapshotCache) CreateWatch(arg0 envoy_service_discovery_v3.DiscoveryRequest) (chan cache.Response, func()) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWatch", arg0)
	ret0, _ := ret[0].(chan cache.Response)
	ret1, _ := ret[1].(func())
	return ret0, ret1
}

// CreateWatch indicates an expected call of CreateWatch
func (mr *MockSnapshotCacheMockRecorder) CreateWatch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWatch", reflect.TypeOf((*MockSnapshotCache)(nil).CreateWatch), arg0)
}

// Fetch mocks base method
func (m *MockSnapshotCache) Fetch(arg0 context.Context, arg1 envoy_service_discovery_v3.DiscoveryRequest) (*cache.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0, arg1)
	ret0, _ := ret[0].(*cache.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch
func (mr *MockSnapshotCacheMockRecorder) Fetch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockSnapshotCache)(nil).Fetch), arg0, arg1)
}

// GetSnapshot mocks base method
func (m *MockSnapshotCache) GetSnapshot(arg0 string) (cache.Snapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSnapshot", arg0)
	ret0, _ := ret[0].(cache.Snapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSnapshot indicates an expected call of GetSnapshot
func (mr *MockSnapshotCacheMockRecorder) GetSnapshot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSnapshot", reflect.TypeOf((*MockSnapshotCache)(nil).GetSnapshot), arg0)
}

// GetStatusInfo mocks base method
func (m *MockSnapshotCache) GetStatusInfo(arg0 string) cache.StatusInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatusInfo", arg0)
	ret0, _ := ret[0].(cache.StatusInfo)
	return ret0
}

// GetStatusInfo indicates an expected call of GetStatusInfo
func (mr *MockSnapshotCacheMockRecorder) GetStatusInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatusInfo", reflect.TypeOf((*MockSnapshotCache)(nil).GetStatusInfo), arg0)
}

// GetStatusKeys mocks base method
func (m *MockSnapshotCache) GetStatusKeys() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatusKeys")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetStatusKeys indicates an expected call of GetStatusKeys
func (mr *MockSnapshotCacheMockRecorder) GetStatusKeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatusKeys", reflect.TypeOf((*MockSnapshotCache)(nil).GetStatusKeys))
}

// SetSnapshot mocks base method
func (m *MockSnapshotCache) SetSnapshot(arg0 string, arg1 cache.Snapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSnapshot", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSnapshot indicates an expected call of SetSnapshot
func (mr *MockSnapshotCacheMockRecorder) SetSnapshot(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSnapshot", reflect.TypeOf((*MockSnapshotCache)(nil).SetSnapshot), arg0, arg1)
}
