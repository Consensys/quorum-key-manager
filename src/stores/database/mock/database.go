// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	database "github.com/consensys/quorum-key-manager/src/stores/database"
	entities "github.com/consensys/quorum-key-manager/src/stores/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDatabase is a mock of Database interface
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// ETHAccounts mocks base method
func (m *MockDatabase) ETHAccounts(storeID string) database.ETHAccounts {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ETHAccounts", storeID)
	ret0, _ := ret[0].(database.ETHAccounts)
	return ret0
}

// ETHAccounts indicates an expected call of ETHAccounts
func (mr *MockDatabaseMockRecorder) ETHAccounts(storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ETHAccounts", reflect.TypeOf((*MockDatabase)(nil).ETHAccounts), storeID)
}

// Keys mocks base method
func (m *MockDatabase) Keys(storeID string) database.Keys {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Keys", storeID)
	ret0, _ := ret[0].(database.Keys)
	return ret0
}

// Keys indicates an expected call of Keys
func (mr *MockDatabaseMockRecorder) Keys(storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Keys", reflect.TypeOf((*MockDatabase)(nil).Keys), storeID)
}

// Secrets mocks base method
func (m *MockDatabase) Secrets(storeID string) database.Secrets {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Secrets", storeID)
	ret0, _ := ret[0].(database.Secrets)
	return ret0
}

// Secrets indicates an expected call of Secrets
func (mr *MockDatabaseMockRecorder) Secrets(storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Secrets", reflect.TypeOf((*MockDatabase)(nil).Secrets), storeID)
}

// MockETHAccounts is a mock of ETHAccounts interface
type MockETHAccounts struct {
	ctrl     *gomock.Controller
	recorder *MockETHAccountsMockRecorder
}

// MockETHAccountsMockRecorder is the mock recorder for MockETHAccounts
type MockETHAccountsMockRecorder struct {
	mock *MockETHAccounts
}

// NewMockETHAccounts creates a new mock instance
func NewMockETHAccounts(ctrl *gomock.Controller) *MockETHAccounts {
	mock := &MockETHAccounts{ctrl: ctrl}
	mock.recorder = &MockETHAccountsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockETHAccounts) EXPECT() *MockETHAccountsMockRecorder {
	return m.recorder
}

// RunInTransaction mocks base method
func (m *MockETHAccounts) RunInTransaction(ctx context.Context, persistFunc func(database.ETHAccounts) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTransaction", ctx, persistFunc)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTransaction indicates an expected call of RunInTransaction
func (mr *MockETHAccountsMockRecorder) RunInTransaction(ctx, persistFunc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTransaction", reflect.TypeOf((*MockETHAccounts)(nil).RunInTransaction), ctx, persistFunc)
}

// Get mocks base method
func (m *MockETHAccounts) Get(ctx context.Context, addr string) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, addr)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockETHAccountsMockRecorder) Get(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockETHAccounts)(nil).Get), ctx, addr)
}

// GetDeleted mocks base method
func (m *MockETHAccounts) GetDeleted(ctx context.Context, addr string) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeleted", ctx, addr)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeleted indicates an expected call of GetDeleted
func (mr *MockETHAccountsMockRecorder) GetDeleted(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeleted", reflect.TypeOf((*MockETHAccounts)(nil).GetDeleted), ctx, addr)
}

// GetAll mocks base method
func (m *MockETHAccounts) GetAll(ctx context.Context) ([]*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockETHAccountsMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockETHAccounts)(nil).GetAll), ctx)
}

// GetAllDeleted mocks base method
func (m *MockETHAccounts) GetAllDeleted(ctx context.Context) ([]*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDeleted", ctx)
	ret0, _ := ret[0].([]*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDeleted indicates an expected call of GetAllDeleted
func (mr *MockETHAccountsMockRecorder) GetAllDeleted(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDeleted", reflect.TypeOf((*MockETHAccounts)(nil).GetAllDeleted), ctx)
}

// Add mocks base method
func (m *MockETHAccounts) Add(ctx context.Context, account *entities.ETHAccount) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, account)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add
func (mr *MockETHAccountsMockRecorder) Add(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockETHAccounts)(nil).Add), ctx, account)
}

// Update mocks base method
func (m *MockETHAccounts) Update(ctx context.Context, account *entities.ETHAccount) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, account)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockETHAccountsMockRecorder) Update(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockETHAccounts)(nil).Update), ctx, account)
}

// Delete mocks base method
func (m *MockETHAccounts) Delete(ctx context.Context, addr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockETHAccountsMockRecorder) Delete(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockETHAccounts)(nil).Delete), ctx, addr)
}

// Restore mocks base method
func (m *MockETHAccounts) Restore(ctx context.Context, addr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore
func (mr *MockETHAccountsMockRecorder) Restore(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockETHAccounts)(nil).Restore), ctx, addr)
}

// Purge mocks base method
func (m *MockETHAccounts) Purge(ctx context.Context, addr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Purge", ctx, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Purge indicates an expected call of Purge
func (mr *MockETHAccountsMockRecorder) Purge(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Purge", reflect.TypeOf((*MockETHAccounts)(nil).Purge), ctx, addr)
}

// MockKeys is a mock of Keys interface
type MockKeys struct {
	ctrl     *gomock.Controller
	recorder *MockKeysMockRecorder
}

// MockKeysMockRecorder is the mock recorder for MockKeys
type MockKeysMockRecorder struct {
	mock *MockKeys
}

// NewMockKeys creates a new mock instance
func NewMockKeys(ctrl *gomock.Controller) *MockKeys {
	mock := &MockKeys{ctrl: ctrl}
	mock.recorder = &MockKeysMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeys) EXPECT() *MockKeysMockRecorder {
	return m.recorder
}

// RunInTransaction mocks base method
func (m *MockKeys) RunInTransaction(ctx context.Context, persistFunc func(database.Keys) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTransaction", ctx, persistFunc)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTransaction indicates an expected call of RunInTransaction
func (mr *MockKeysMockRecorder) RunInTransaction(ctx, persistFunc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTransaction", reflect.TypeOf((*MockKeys)(nil).RunInTransaction), ctx, persistFunc)
}

// Get mocks base method
func (m *MockKeys) Get(ctx context.Context, id string) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockKeysMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKeys)(nil).Get), ctx, id)
}

// GetDeleted mocks base method
func (m *MockKeys) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeleted", ctx, id)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeleted indicates an expected call of GetDeleted
func (mr *MockKeysMockRecorder) GetDeleted(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeleted", reflect.TypeOf((*MockKeys)(nil).GetDeleted), ctx, id)
}

// GetAll mocks base method
func (m *MockKeys) GetAll(ctx context.Context) ([]*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockKeysMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockKeys)(nil).GetAll), ctx)
}

// GetAllDeleted mocks base method
func (m *MockKeys) GetAllDeleted(ctx context.Context) ([]*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDeleted", ctx)
	ret0, _ := ret[0].([]*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDeleted indicates an expected call of GetAllDeleted
func (mr *MockKeysMockRecorder) GetAllDeleted(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDeleted", reflect.TypeOf((*MockKeys)(nil).GetAllDeleted), ctx)
}

// Add mocks base method
func (m *MockKeys) Add(ctx context.Context, key *entities.Key) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, key)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add
func (mr *MockKeysMockRecorder) Add(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockKeys)(nil).Add), ctx, key)
}

// Update mocks base method
func (m *MockKeys) Update(ctx context.Context, key *entities.Key) (*entities.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, key)
	ret0, _ := ret[0].(*entities.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockKeysMockRecorder) Update(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockKeys)(nil).Update), ctx, key)
}

// Delete mocks base method
func (m *MockKeys) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockKeysMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockKeys)(nil).Delete), ctx, id)
}

// Restore mocks base method
func (m *MockKeys) Restore(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore
func (mr *MockKeysMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockKeys)(nil).Restore), ctx, id)
}

// Purge mocks base method
func (m *MockKeys) Purge(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Purge", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Purge indicates an expected call of Purge
func (mr *MockKeysMockRecorder) Purge(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Purge", reflect.TypeOf((*MockKeys)(nil).Purge), ctx, id)
}

// MockSecrets is a mock of Secrets interface
type MockSecrets struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsMockRecorder
}

// MockSecretsMockRecorder is the mock recorder for MockSecrets
type MockSecretsMockRecorder struct {
	mock *MockSecrets
}

// NewMockSecrets creates a new mock instance
func NewMockSecrets(ctrl *gomock.Controller) *MockSecrets {
	mock := &MockSecrets{ctrl: ctrl}
	mock.recorder = &MockSecretsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSecrets) EXPECT() *MockSecretsMockRecorder {
	return m.recorder
}

// RunInTransaction mocks base method
func (m *MockSecrets) RunInTransaction(ctx context.Context, persistFunc func(database.Secrets) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTransaction", ctx, persistFunc)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTransaction indicates an expected call of RunInTransaction
func (mr *MockSecretsMockRecorder) RunInTransaction(ctx, persistFunc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTransaction", reflect.TypeOf((*MockSecrets)(nil).RunInTransaction), ctx, persistFunc)
}

// Get mocks base method
func (m *MockSecrets) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id, version)
	ret0, _ := ret[0].(*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockSecretsMockRecorder) Get(ctx, id, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSecrets)(nil).Get), ctx, id, version)
}

// GetLatestVersion mocks base method
func (m *MockSecrets) GetLatestVersion(ctx context.Context, id string, isDeleted bool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestVersion", ctx, id, isDeleted)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestVersion indicates an expected call of GetLatestVersion
func (mr *MockSecretsMockRecorder) GetLatestVersion(ctx, id, isDeleted interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestVersion", reflect.TypeOf((*MockSecrets)(nil).GetLatestVersion), ctx, id, isDeleted)
}

// ListVersions mocks base method
func (m *MockSecrets) ListVersions(ctx context.Context, id string, isDeleted bool) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListVersions", ctx, id, isDeleted)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListVersions indicates an expected call of ListVersions
func (mr *MockSecretsMockRecorder) ListVersions(ctx, id, isDeleted interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVersions", reflect.TypeOf((*MockSecrets)(nil).ListVersions), ctx, id, isDeleted)
}

// GetDeleted mocks base method
func (m *MockSecrets) GetDeleted(ctx context.Context, id string) (*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeleted", ctx, id)
	ret0, _ := ret[0].(*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeleted indicates an expected call of GetDeleted
func (mr *MockSecretsMockRecorder) GetDeleted(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeleted", reflect.TypeOf((*MockSecrets)(nil).GetDeleted), ctx, id)
}

// GetAll mocks base method
func (m *MockSecrets) GetAll(ctx context.Context) ([]*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockSecretsMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockSecrets)(nil).GetAll), ctx)
}

// GetAllDeleted mocks base method
func (m *MockSecrets) GetAllDeleted(ctx context.Context) ([]*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDeleted", ctx)
	ret0, _ := ret[0].([]*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDeleted indicates an expected call of GetAllDeleted
func (mr *MockSecretsMockRecorder) GetAllDeleted(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDeleted", reflect.TypeOf((*MockSecrets)(nil).GetAllDeleted), ctx)
}

// Add mocks base method
func (m *MockSecrets) Add(ctx context.Context, secret *entities.Secret) (*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, secret)
	ret0, _ := ret[0].(*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add
func (mr *MockSecretsMockRecorder) Add(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockSecrets)(nil).Add), ctx, secret)
}

// Update mocks base method
func (m *MockSecrets) Update(ctx context.Context, secret *entities.Secret) (*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, secret)
	ret0, _ := ret[0].(*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockSecretsMockRecorder) Update(ctx, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSecrets)(nil).Update), ctx, secret)
}

// Delete mocks base method
func (m *MockSecrets) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockSecretsMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSecrets)(nil).Delete), ctx, id)
}

// Restore mocks base method
func (m *MockSecrets) Restore(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore
func (mr *MockSecretsMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockSecrets)(nil).Restore), ctx, id)
}

// Purge mocks base method
func (m *MockSecrets) Purge(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Purge", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Purge indicates an expected call of Purge
func (mr *MockSecretsMockRecorder) Purge(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Purge", reflect.TypeOf((*MockSecrets)(nil).Purge), ctx, id)
}
