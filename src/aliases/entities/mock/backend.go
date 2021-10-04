// Code generated by MockGen. DO NOT EDIT.
// Source: backend.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	gomock "github.com/golang/mock/gomock"
)

// MockAliasBackend is a mock of AliasBackend interface.
type MockAliasBackend struct {
	ctrl     *gomock.Controller
	recorder *MockAliasBackendMockRecorder
}

// MockAliasBackendMockRecorder is the mock recorder for MockAliasBackend.
type MockAliasBackendMockRecorder struct {
	mock *MockAliasBackend
}

// NewMockAliasBackend creates a new mock instance.
func NewMockAliasBackend(ctrl *gomock.Controller) *MockAliasBackend {
	mock := &MockAliasBackend{ctrl: ctrl}
	mock.recorder = &MockAliasBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAliasBackend) EXPECT() *MockAliasBackendMockRecorder {
	return m.recorder
}

// CreateAlias mocks base method.
func (m *MockAliasBackend) CreateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAlias", ctx, registry, alias)
	ret0, _ := ret[0].(*aliasent.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAlias indicates an expected call of CreateAlias.
func (mr *MockAliasBackendMockRecorder) CreateAlias(ctx, registry, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAlias", reflect.TypeOf((*MockAliasBackend)(nil).CreateAlias), ctx, registry, alias)
}

// DeleteAlias mocks base method.
func (m *MockAliasBackend) DeleteAlias(ctx context.Context, registry, aliasKey string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAlias", ctx, registry, aliasKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAlias indicates an expected call of DeleteAlias.
func (mr *MockAliasBackendMockRecorder) DeleteAlias(ctx, registry, aliasKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAlias", reflect.TypeOf((*MockAliasBackend)(nil).DeleteAlias), ctx, registry, aliasKey)
}

// DeleteRegistry mocks base method.
func (m *MockAliasBackend) DeleteRegistry(ctx context.Context, registry string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRegistry", ctx, registry)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRegistry indicates an expected call of DeleteRegistry.
func (mr *MockAliasBackendMockRecorder) DeleteRegistry(ctx, registry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRegistry", reflect.TypeOf((*MockAliasBackend)(nil).DeleteRegistry), ctx, registry)
}

// GetAlias mocks base method.
func (m *MockAliasBackend) GetAlias(ctx context.Context, registry, aliasKey string) (*aliasent.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlias", ctx, registry, aliasKey)
	ret0, _ := ret[0].(*aliasent.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlias indicates an expected call of GetAlias.
func (mr *MockAliasBackendMockRecorder) GetAlias(ctx, registry, aliasKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlias", reflect.TypeOf((*MockAliasBackend)(nil).GetAlias), ctx, registry, aliasKey)
}

// ListAliases mocks base method.
func (m *MockAliasBackend) ListAliases(ctx context.Context, registry string) ([]aliasent.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAliases", ctx, registry)
	ret0, _ := ret[0].([]aliasent.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAliases indicates an expected call of ListAliases.
func (mr *MockAliasBackendMockRecorder) ListAliases(ctx, registry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAliases", reflect.TypeOf((*MockAliasBackend)(nil).ListAliases), ctx, registry)
}

// UpdateAlias mocks base method.
func (m *MockAliasBackend) UpdateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAlias", ctx, registry, alias)
	ret0, _ := ret[0].(*aliasent.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAlias indicates an expected call of UpdateAlias.
func (mr *MockAliasBackendMockRecorder) UpdateAlias(ctx, registry, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAlias", reflect.TypeOf((*MockAliasBackend)(nil).UpdateAlias), ctx, registry, alias)
}
