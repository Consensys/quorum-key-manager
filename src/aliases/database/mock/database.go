// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	entities "github.com/consensys/quorum-key-manager/src/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRegistry is a mock of Registry interface
type MockRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockRegistryMockRecorder
}

// MockRegistryMockRecorder is the mock recorder for MockRegistry
type MockRegistryMockRecorder struct {
	mock *MockRegistry
}

// NewMockRegistry creates a new mock instance
func NewMockRegistry(ctrl *gomock.Controller) *MockRegistry {
	mock := &MockRegistry{ctrl: ctrl}
	mock.recorder = &MockRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRegistry) EXPECT() *MockRegistryMockRecorder {
	return m.recorder
}

// Insert mocks base method
func (m *MockRegistry) Insert(ctx context.Context, registry *entities.AliasRegistry) (*entities.AliasRegistry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, registry)
	ret0, _ := ret[0].(*entities.AliasRegistry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert
func (mr *MockRegistryMockRecorder) Insert(ctx, registry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockRegistry)(nil).Insert), ctx, registry)
}

// FindOne mocks base method
func (m *MockRegistry) FindOne(ctx context.Context, name, tenant string) (*entities.AliasRegistry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", ctx, name, tenant)
	ret0, _ := ret[0].(*entities.AliasRegistry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne
func (mr *MockRegistryMockRecorder) FindOne(ctx, name, tenant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockRegistry)(nil).FindOne), ctx, name, tenant)
}

// Delete mocks base method
func (m *MockRegistry) Delete(ctx context.Context, name, tenant string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, name, tenant)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRegistryMockRecorder) Delete(ctx, name, tenant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRegistry)(nil).Delete), ctx, name, tenant)
}

// MockAlias is a mock of Alias interface
type MockAlias struct {
	ctrl     *gomock.Controller
	recorder *MockAliasMockRecorder
}

// MockAliasMockRecorder is the mock recorder for MockAlias
type MockAliasMockRecorder struct {
	mock *MockAlias
}

// NewMockAlias creates a new mock instance
func NewMockAlias(ctrl *gomock.Controller) *MockAlias {
	mock := &MockAlias{ctrl: ctrl}
	mock.recorder = &MockAliasMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAlias) EXPECT() *MockAliasMockRecorder {
	return m.recorder
}

// Insert mocks base method
func (m *MockAlias) Insert(ctx context.Context, alias *entities.Alias) (*entities.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, alias)
	ret0, _ := ret[0].(*entities.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert
func (mr *MockAliasMockRecorder) Insert(ctx, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockAlias)(nil).Insert), ctx, alias)
}

// FindOne mocks base method
func (m *MockAlias) FindOne(ctx context.Context, registry, key, tenant string) (*entities.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", ctx, registry, key, tenant)
	ret0, _ := ret[0].(*entities.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne
func (mr *MockAliasMockRecorder) FindOne(ctx, registry, key, tenant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockAlias)(nil).FindOne), ctx, registry, key, tenant)
}

// Update mocks base method
func (m *MockAlias) Update(ctx context.Context, alias *entities.Alias) (*entities.Alias, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, alias)
	ret0, _ := ret[0].(*entities.Alias)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockAliasMockRecorder) Update(ctx, alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAlias)(nil).Update), ctx, alias)
}

// Delete mocks base method
func (m *MockAlias) Delete(ctx context.Context, registry, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, registry, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockAliasMockRecorder) Delete(ctx, registry, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAlias)(nil).Delete), ctx, registry, key)
}