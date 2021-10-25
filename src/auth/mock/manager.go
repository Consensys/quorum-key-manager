// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock is a generated GoMock package.
package mock

import (
	types "github.com/consensys/quorum-key-manager/src/auth/entities"
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

// Role mocks base method
func (m *MockManager) Role(name string) (*types.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Role", name)
	ret0, _ := ret[0].(*types.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Role indicates an expected call of Role
func (mr *MockManagerMockRecorder) Role(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Role", reflect.TypeOf((*MockManager)(nil).Role), name)
}

// Roles mocks base method
func (m *MockManager) Roles() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Roles")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Roles indicates an expected call of Roles
func (mr *MockManagerMockRecorder) Roles() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Roles", reflect.TypeOf((*MockManager)(nil).Roles))
}

// UserPermissions mocks base method
func (m *MockManager) UserPermissions(info *types.UserInfo) []types.Permission {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserPermissions", info)
	ret0, _ := ret[0].([]types.Permission)
	return ret0
}

// UserPermissions indicates an expected call of UserPermissions
func (mr *MockManagerMockRecorder) UserPermissions(info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserPermissions", reflect.TypeOf((*MockManager)(nil).UserPermissions), info)
}
