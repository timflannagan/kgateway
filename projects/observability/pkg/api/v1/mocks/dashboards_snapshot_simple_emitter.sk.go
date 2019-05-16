// Code generated by MockGen. DO NOT EDIT.
// Source: ./projects/observability/pkg/api/v1/dashboards_snapshot_simple_emitter.sk.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/solo-io/solo-projects/projects/observability/pkg/api/v1"
)

// MockDashboardsSimpleEmitter is a mock of DashboardsSimpleEmitter interface
type MockDashboardsSimpleEmitter struct {
	ctrl     *gomock.Controller
	recorder *MockDashboardsSimpleEmitterMockRecorder
}

// MockDashboardsSimpleEmitterMockRecorder is the mock recorder for MockDashboardsSimpleEmitter
type MockDashboardsSimpleEmitterMockRecorder struct {
	mock *MockDashboardsSimpleEmitter
}

// NewMockDashboardsSimpleEmitter creates a new mock instance
func NewMockDashboardsSimpleEmitter(ctrl *gomock.Controller) *MockDashboardsSimpleEmitter {
	mock := &MockDashboardsSimpleEmitter{ctrl: ctrl}
	mock.recorder = &MockDashboardsSimpleEmitterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDashboardsSimpleEmitter) EXPECT() *MockDashboardsSimpleEmitterMockRecorder {
	return m.recorder
}

// Snapshots mocks base method
func (m *MockDashboardsSimpleEmitter) Snapshots(ctx context.Context) (<-chan *v1.DashboardsSnapshot, <-chan error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Snapshots", ctx)
	ret0, _ := ret[0].(<-chan *v1.DashboardsSnapshot)
	ret1, _ := ret[1].(<-chan error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Snapshots indicates an expected call of Snapshots
func (mr *MockDashboardsSimpleEmitterMockRecorder) Snapshots(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Snapshots", reflect.TypeOf((*MockDashboardsSimpleEmitter)(nil).Snapshots), ctx)
}
