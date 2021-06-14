// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	node "github.com/consensysquorum/quorum-key-manager/src/nodes/node"
	gomock "github.com/golang/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockManager) List(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockManagerMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockManager)(nil).List), ctx)
}

// Node mocks base method.
func (m *MockManager) Node(ctx context.Context, name string) (node.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Node", ctx, name)
	ret0, _ := ret[0].(node.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Node indicates an expected call of Node.
func (mr *MockManagerMockRecorder) Node(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Node", reflect.TypeOf((*MockManager)(nil).Node), ctx, name)
}
