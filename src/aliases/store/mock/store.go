// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	aliasmodels "github.com/consensys/quorum-key-manager/src/aliases/models"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateAlias mocks base method.
func (m *MockStore) CreateAlias(ctx context.Context, registry aliasmodels.RegistryName, alias aliasmodels.Alias) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAlias", ctx, registry, alias)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAlias indicates an expected call of CreateAlias.
func (mr *MockStoreMockRecorder) CreateAlias(ctx, registry, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAlias", reflect.TypeOf((*MockStore)(nil).CreateAlias), ctx, registry, alias)
}

// DeleteAlias mocks base method.
func (m *MockStore) DeleteAlias(ctx context.Context, registry aliasmodels.RegistryName, aliasKey aliasmodels.AliasKey) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAlias", ctx, registry, aliasKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAlias indicates an expected call of DeleteAlias.
func (mr *MockStoreMockRecorder) DeleteAlias(ctx, registry, aliasKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAlias", reflect.TypeOf((*MockStore)(nil).DeleteAlias), ctx, registry, aliasKey)
}

// DeleteRegistry mocks base method.
func (m *MockStore) DeleteRegistry(ctx context.Context, registry aliasmodels.RegistryName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRegistry", ctx, registry)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRegistry indicates an expected call of DeleteRegistry.
func (mr *MockStoreMockRecorder) DeleteRegistry(ctx, registry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRegistry", reflect.TypeOf((*MockStore)(nil).DeleteRegistry), ctx, registry)
}

// GetAlias mocks base method.
func (m *MockStore) GetAlias(ctx context.Context, registry aliasmodels.RegistryName, aliasKey aliasmodels.AliasKey) (*aliasmodels.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAlias", ctx, registry, aliasKey)
	ret0, _ := ret[0].(*aliasmodels.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAlias indicates an expected call of GetAlias.
func (mr *MockStoreMockRecorder) GetAlias(ctx, registry, aliasKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAlias", reflect.TypeOf((*MockStore)(nil).GetAlias), ctx, registry, aliasKey)
}

// ListAliases mocks base method.
func (m *MockStore) ListAliases(ctx context.Context, registry aliasmodels.RegistryName) ([]aliasmodels.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAliases", ctx, registry)
	ret0, _ := ret[0].([]aliasmodels.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAliases indicates an expected call of ListAliases.
func (mr *MockStoreMockRecorder) ListAliases(ctx, registry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAliases", reflect.TypeOf((*MockStore)(nil).ListAliases), ctx, registry)
}

// UpdateAlias mocks base method.
func (m *MockStore) UpdateAlias(ctx context.Context, registry aliasmodels.RegistryName, alias aliasmodels.Alias) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAlias", ctx, registry, alias)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAlias indicates an expected call of UpdateAlias.
func (mr *MockStoreMockRecorder) UpdateAlias(ctx, registry, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAlias", reflect.TypeOf((*MockStore)(nil).UpdateAlias), ctx, registry, alias)
}
