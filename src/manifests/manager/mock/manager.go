// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock is a generated GoMock package.
package mock

import (
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Load mocks base method
func (m *MockManager) Load(kinds []manifest.Kind) ([]manifest.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", kinds)
	ret0, _ := ret[0].([]manifest.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load
func (mr *MockManagerMockRecorder) Load(kinds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockManager)(nil).Load), kinds)
}
