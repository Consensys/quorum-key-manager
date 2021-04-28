// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	types "github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSecretsClient is a mock of SecretsClient interface
type MockSecretsClient struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsClientMockRecorder
}

// MockSecretsClientMockRecorder is the mock recorder for MockSecretsClient
type MockSecretsClientMockRecorder struct {
	mock *MockSecretsClient
}

// NewMockSecretsClient creates a new mock instance
func NewMockSecretsClient(ctrl *gomock.Controller) *MockSecretsClient {
	mock := &MockSecretsClient{ctrl: ctrl}
	mock.recorder = &MockSecretsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSecretsClient) EXPECT() *MockSecretsClientMockRecorder {
	return m.recorder
}

// SetSecret mocks base method
func (m *MockSecretsClient) SetSecret(ctx context.Context, storeName string, request *types.SetSecretRequest) (*types.SecretResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSecret", ctx, storeName, request)
	ret0, _ := ret[0].(*types.SecretResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetSecret indicates an expected call of SetSecret
func (mr *MockSecretsClientMockRecorder) SetSecret(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSecret", reflect.TypeOf((*MockSecretsClient)(nil).SetSecret), ctx, storeName, request)
}

// GetSecret mocks base method
func (m *MockSecretsClient) GetSecret(ctx context.Context, storeName, id, version string) (*types.SecretResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", ctx, storeName, id, version)
	ret0, _ := ret[0].(*types.SecretResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret
func (mr *MockSecretsClientMockRecorder) GetSecret(ctx, storeName, id, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecretsClient)(nil).GetSecret), ctx, storeName, id, version)
}

// ListSecrets mocks base method
func (m *MockSecretsClient) ListSecrets(ctx context.Context, storeName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecrets", ctx, storeName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecrets indicates an expected call of ListSecrets
func (mr *MockSecretsClientMockRecorder) ListSecrets(ctx, storeName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecrets", reflect.TypeOf((*MockSecretsClient)(nil).ListSecrets), ctx, storeName)
}

// MockKeyManagerClient is a mock of KeyManagerClient interface
type MockKeyManagerClient struct {
	ctrl     *gomock.Controller
	recorder *MockKeyManagerClientMockRecorder
}

// MockKeyManagerClientMockRecorder is the mock recorder for MockKeyManagerClient
type MockKeyManagerClientMockRecorder struct {
	mock *MockKeyManagerClient
}

// NewMockKeyManagerClient creates a new mock instance
func NewMockKeyManagerClient(ctrl *gomock.Controller) *MockKeyManagerClient {
	mock := &MockKeyManagerClient{ctrl: ctrl}
	mock.recorder = &MockKeyManagerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeyManagerClient) EXPECT() *MockKeyManagerClientMockRecorder {
	return m.recorder
}

// SetSecret mocks base method
func (m *MockKeyManagerClient) SetSecret(ctx context.Context, storeName string, request *types.SetSecretRequest) (*types.SecretResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetSecret", ctx, storeName, request)
	ret0, _ := ret[0].(*types.SecretResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetSecret indicates an expected call of SetSecret
func (mr *MockKeyManagerClientMockRecorder) SetSecret(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSecret", reflect.TypeOf((*MockKeyManagerClient)(nil).SetSecret), ctx, storeName, request)
}

// GetSecret mocks base method
func (m *MockKeyManagerClient) GetSecret(ctx context.Context, storeName, id, version string) (*types.SecretResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", ctx, storeName, id, version)
	ret0, _ := ret[0].(*types.SecretResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret
func (mr *MockKeyManagerClientMockRecorder) GetSecret(ctx, storeName, id, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockKeyManagerClient)(nil).GetSecret), ctx, storeName, id, version)
}

// ListSecrets mocks base method
func (m *MockKeyManagerClient) ListSecrets(ctx context.Context, storeName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecrets", ctx, storeName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecrets indicates an expected call of ListSecrets
func (mr *MockKeyManagerClientMockRecorder) ListSecrets(ctx, storeName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecrets", reflect.TypeOf((*MockKeyManagerClient)(nil).ListSecrets), ctx, storeName)
}
