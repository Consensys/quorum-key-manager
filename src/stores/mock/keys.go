// Code generated by MockGen. DO NOT EDIT.
// Source: keys.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	entities "github.com/consensys/quorum-key-manager/src/stores/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockKeyStore is a mock of KeyStore interface
type MockKeyStore struct {
	ctrl     *gomock.Controller
	recorder *MockKeyStoreMockRecorder
}

// MockKeyStoreMockRecorder is the mock recorder for MockKeyStore
type MockKeyStoreMockRecorder struct {
	mock *MockKeyStore
}

// NewMockKeyStore creates a new mock instance
func NewMockKeyStore(ctrl *gomock.Controller) *MockKeyStore {
	mock := &MockKeyStore{ctrl: ctrl}
	mock.recorder = &MockKeyStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeyStore) EXPECT() *MockKeyStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockKeyStore) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, id, alg, attr)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockKeyStoreMockRecorder) Create(ctx, id, alg, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockKeyStore)(nil).Create), ctx, id, alg, attr)
}

// Import mocks base method
func (m *MockKeyStore) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Import", ctx, id, privKey, alg, attr)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Import indicates an expected call of Import
func (mr *MockKeyStoreMockRecorder) Import(ctx, id, privKey, alg, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Import", reflect.TypeOf((*MockKeyStore)(nil).Import), ctx, id, privKey, alg, attr)
}

// Get mocks base method
func (m *MockKeyStore) Get(ctx context.Context, id string) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockKeyStoreMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKeyStore)(nil).Get), ctx, id)
}

// List mocks base method
func (m *MockKeyStore) List(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockKeyStoreMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockKeyStore)(nil).List), ctx)
}

// Update mocks base method
func (m *MockKeyStore) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, attr)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockKeyStoreMockRecorder) Update(ctx, id, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockKeyStore)(nil).Update), ctx, id, attr)
}

// Delete mocks base method
func (m *MockKeyStore) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockKeyStoreMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockKeyStore)(nil).Delete), ctx, id)
}

// GetDeleted mocks base method
func (m *MockKeyStore) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeleted", ctx, id)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeleted indicates an expected call of GetDeleted
func (mr *MockKeyStoreMockRecorder) GetDeleted(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeleted", reflect.TypeOf((*MockKeyStore)(nil).GetDeleted), ctx, id)
}

// ListDeleted mocks base method
func (m *MockKeyStore) ListDeleted(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDeleted", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDeleted indicates an expected call of ListDeleted
func (mr *MockKeyStoreMockRecorder) ListDeleted(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeleted", reflect.TypeOf((*MockKeyStore)(nil).ListDeleted), ctx)
}

// Restore mocks base method
func (m *MockKeyStore) Restore(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore
func (mr *MockKeyStoreMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockKeyStore)(nil).Restore), ctx, id)
}

// Destroy mocks base method
func (m *MockKeyStore) Destroy(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockKeyStoreMockRecorder) Destroy(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockKeyStore)(nil).Destroy), ctx, id)
}

// Sign mocks base method
func (m *MockKeyStore) Sign(ctx context.Context, id string, data []byte, algo *entities.Algorithm) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", ctx, id, data, algo)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign
func (mr *MockKeyStoreMockRecorder) Sign(ctx, id, data, algo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockKeyStore)(nil).Sign), ctx, id, data, algo)
}

// Verify mocks base method
func (m *MockKeyStore) Verify(ctx context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", ctx, pubKey, data, sig, algo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify
func (mr *MockKeyStoreMockRecorder) Verify(ctx, pubKey, data, sig, algo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockKeyStore)(nil).Verify), ctx, pubKey, data, sig, algo)
}

// Encrypt mocks base method
func (m *MockKeyStore) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encrypt", ctx, id, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encrypt indicates an expected call of Encrypt
func (mr *MockKeyStoreMockRecorder) Encrypt(ctx, id, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encrypt", reflect.TypeOf((*MockKeyStore)(nil).Encrypt), ctx, id, data)
}

// Decrypt mocks base method
func (m *MockKeyStore) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decrypt", ctx, id, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decrypt indicates an expected call of Decrypt
func (mr *MockKeyStoreMockRecorder) Decrypt(ctx, id, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrypt", reflect.TypeOf((*MockKeyStore)(nil).Decrypt), ctx, id, data)
}
