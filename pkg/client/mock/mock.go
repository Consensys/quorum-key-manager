// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	jsonrpc "github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	types "github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/api/types"
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

// MockKeysClient is a mock of KeysClient interface
type MockKeysClient struct {
	ctrl     *gomock.Controller
	recorder *MockKeysClientMockRecorder
}

// MockKeysClientMockRecorder is the mock recorder for MockKeysClient
type MockKeysClientMockRecorder struct {
	mock *MockKeysClient
}

// NewMockKeysClient creates a new mock instance
func NewMockKeysClient(ctrl *gomock.Controller) *MockKeysClient {
	mock := &MockKeysClient{ctrl: ctrl}
	mock.recorder = &MockKeysClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeysClient) EXPECT() *MockKeysClientMockRecorder {
	return m.recorder
}

// CreateKey mocks base method
func (m *MockKeysClient) CreateKey(ctx context.Context, storeName string, request *types.CreateKeyRequest) (*types.KeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateKey", ctx, storeName, request)
	ret0, _ := ret[0].(*types.KeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateKey indicates an expected call of CreateKey
func (mr *MockKeysClientMockRecorder) CreateKey(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateKey", reflect.TypeOf((*MockKeysClient)(nil).CreateKey), ctx, storeName, request)
}

// ImportKey mocks base method
func (m *MockKeysClient) ImportKey(ctx context.Context, storeName string, request *types.ImportKeyRequest) (*types.KeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportKey", ctx, storeName, request)
	ret0, _ := ret[0].(*types.KeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportKey indicates an expected call of ImportKey
func (mr *MockKeysClientMockRecorder) ImportKey(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportKey", reflect.TypeOf((*MockKeysClient)(nil).ImportKey), ctx, storeName, request)
}

// SignKey mocks base method
func (m *MockKeysClient) SignKey(ctx context.Context, storeName, id string, request *types.SignBase64PayloadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignKey", ctx, storeName, id, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignKey indicates an expected call of SignKey
func (mr *MockKeysClientMockRecorder) SignKey(ctx, storeName, id, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignKey", reflect.TypeOf((*MockKeysClient)(nil).SignKey), ctx, storeName, id, request)
}

// GetKey mocks base method
func (m *MockKeysClient) GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKey", ctx, storeName, id)
	ret0, _ := ret[0].(*types.KeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKey indicates an expected call of GetKey
func (mr *MockKeysClientMockRecorder) GetKey(ctx, storeName, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockKeysClient)(nil).GetKey), ctx, storeName, id)
}

// ListKeys mocks base method
func (m *MockKeysClient) ListKeys(ctx context.Context, storeName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListKeys", ctx, storeName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListKeys indicates an expected call of ListKeys
func (mr *MockKeysClientMockRecorder) ListKeys(ctx, storeName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListKeys", reflect.TypeOf((*MockKeysClient)(nil).ListKeys), ctx, storeName)
}

// DestroyKey mocks base method
func (m *MockKeysClient) DestroyKey(ctx context.Context, storeName, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyKey", ctx, storeName, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyKey indicates an expected call of DestroyKey
func (mr *MockKeysClientMockRecorder) DestroyKey(ctx, storeName, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyKey", reflect.TypeOf((*MockKeysClient)(nil).DestroyKey), ctx, storeName, id)
}

// MockEth1Client is a mock of Eth1Client interface
type MockEth1Client struct {
	ctrl     *gomock.Controller
	recorder *MockEth1ClientMockRecorder
}

// MockEth1ClientMockRecorder is the mock recorder for MockEth1Client
type MockEth1ClientMockRecorder struct {
	mock *MockEth1Client
}

// NewMockEth1Client creates a new mock instance
func NewMockEth1Client(ctrl *gomock.Controller) *MockEth1Client {
	mock := &MockEth1Client{ctrl: ctrl}
	mock.recorder = &MockEth1ClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEth1Client) EXPECT() *MockEth1ClientMockRecorder {
	return m.recorder
}

// CreateEth1Account mocks base method
func (m *MockEth1Client) CreateEth1Account(ctx context.Context, storeName string, request *types.CreateEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEth1Account", ctx, storeName, request)
	ret0, _ := ret[0].(*types.Eth1AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEth1Account indicates an expected call of CreateEth1Account
func (mr *MockEth1ClientMockRecorder) CreateEth1Account(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEth1Account", reflect.TypeOf((*MockEth1Client)(nil).CreateEth1Account), ctx, storeName, request)
}

// ImportEth1Account mocks base method
func (m *MockEth1Client) ImportEth1Account(ctx context.Context, storeName string, request *types.ImportEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportEth1Account", ctx, storeName, request)
	ret0, _ := ret[0].(*types.Eth1AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportEth1Account indicates an expected call of ImportEth1Account
func (mr *MockEth1ClientMockRecorder) ImportEth1Account(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportEth1Account", reflect.TypeOf((*MockEth1Client)(nil).ImportEth1Account), ctx, storeName, request)
}

// UpdateEth1Account mocks base method
func (m *MockEth1Client) UpdateEth1Account(ctx context.Context, storeName string, request *types.UpdateEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEth1Account", ctx, storeName, request)
	ret0, _ := ret[0].(*types.Eth1AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEth1Account indicates an expected call of UpdateEth1Account
func (mr *MockEth1ClientMockRecorder) UpdateEth1Account(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEth1Account", reflect.TypeOf((*MockEth1Client)(nil).UpdateEth1Account), ctx, storeName, request)
}

// SignEth1 mocks base method
func (m *MockEth1Client) SignEth1(ctx context.Context, storeName string, request *types.SignHexPayloadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignEth1", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignEth1 indicates an expected call of SignEth1
func (mr *MockEth1ClientMockRecorder) SignEth1(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignEth1", reflect.TypeOf((*MockEth1Client)(nil).SignEth1), ctx, storeName, request)
}

// SignTypedData mocks base method
func (m *MockEth1Client) SignTypedData(ctx context.Context, storeName string, request *types.SignTypedDataRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTypedData", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTypedData indicates an expected call of SignTypedData
func (mr *MockEth1ClientMockRecorder) SignTypedData(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTypedData", reflect.TypeOf((*MockEth1Client)(nil).SignTypedData), ctx, storeName, request)
}

// SignTransaction mocks base method
func (m *MockEth1Client) SignTransaction(ctx context.Context, storeName string, request *types.SignETHTransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTransaction", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTransaction indicates an expected call of SignTransaction
func (mr *MockEth1ClientMockRecorder) SignTransaction(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTransaction", reflect.TypeOf((*MockEth1Client)(nil).SignTransaction), ctx, storeName, request)
}

// SignQuorumPrivateTransaction mocks base method
func (m *MockEth1Client) SignQuorumPrivateTransaction(ctx context.Context, storeName string, request *types.SignQuorumPrivateTransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignQuorumPrivateTransaction", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignQuorumPrivateTransaction indicates an expected call of SignQuorumPrivateTransaction
func (mr *MockEth1ClientMockRecorder) SignQuorumPrivateTransaction(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignQuorumPrivateTransaction", reflect.TypeOf((*MockEth1Client)(nil).SignQuorumPrivateTransaction), ctx, storeName, request)
}

// SignEEATransaction mocks base method
func (m *MockEth1Client) SignEEATransaction(ctx context.Context, storeName string, request *types.SignEEATransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignEEATransaction", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignEEATransaction indicates an expected call of SignEEATransaction
func (mr *MockEth1ClientMockRecorder) SignEEATransaction(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignEEATransaction", reflect.TypeOf((*MockEth1Client)(nil).SignEEATransaction), ctx, storeName, request)
}

// GetEth1Account mocks base method
func (m *MockEth1Client) GetEth1Account(ctx context.Context, storeName, account string) (*types.Eth1AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEth1Account", ctx, storeName, account)
	ret0, _ := ret[0].(*types.Eth1AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEth1Account indicates an expected call of GetEth1Account
func (mr *MockEth1ClientMockRecorder) GetEth1Account(ctx, storeName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEth1Account", reflect.TypeOf((*MockEth1Client)(nil).GetEth1Account), ctx, storeName, account)
}

// ListEth1Account mocks base method
func (m *MockEth1Client) ListEth1Account(ctx context.Context, storeName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEth1Account", ctx, storeName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEth1Account indicates an expected call of ListEth1Account
func (mr *MockEth1ClientMockRecorder) ListEth1Account(ctx, storeName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEth1Account", reflect.TypeOf((*MockEth1Client)(nil).ListEth1Account), ctx, storeName)
}

// DeleteEth1Account mocks base method
func (m *MockEth1Client) DeleteEth1Account(ctx context.Context, storeName, account string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEth1Account", ctx, storeName, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEth1Account indicates an expected call of DeleteEth1Account
func (mr *MockEth1ClientMockRecorder) DeleteEth1Account(ctx, storeName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEth1Account", reflect.TypeOf((*MockEth1Client)(nil).DeleteEth1Account), ctx, storeName, account)
}

// DestroyEth1Account mocks base method
func (m *MockEth1Client) DestroyEth1Account(ctx context.Context, storeName, account string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyEth1Account", ctx, storeName, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyEth1Account indicates an expected call of DestroyEth1Account
func (mr *MockEth1ClientMockRecorder) DestroyEth1Account(ctx, storeName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyEth1Account", reflect.TypeOf((*MockEth1Client)(nil).DestroyEth1Account), ctx, storeName, account)
}

// RestoreEth1Account mocks base method
func (m *MockEth1Client) RestoreEth1Account(ctx context.Context, storeName, account string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreEth1Account", ctx, storeName, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreEth1Account indicates an expected call of RestoreEth1Account
func (mr *MockEth1ClientMockRecorder) RestoreEth1Account(ctx, storeName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreEth1Account", reflect.TypeOf((*MockEth1Client)(nil).RestoreEth1Account), ctx, storeName, account)
}

// ECRecover mocks base method
func (m *MockEth1Client) ECRecover(ctx context.Context, storeName string, request *types.ECRecoverRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ECRecover", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ECRecover indicates an expected call of ECRecover
func (mr *MockEth1ClientMockRecorder) ECRecover(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ECRecover", reflect.TypeOf((*MockEth1Client)(nil).ECRecover), ctx, storeName, request)
}

// VerifySignature mocks base method
func (m *MockEth1Client) VerifySignature(ctx context.Context, storeName string, request *types.VerifyEth1SignatureRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySignature", ctx, storeName, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifySignature indicates an expected call of VerifySignature
func (mr *MockEth1ClientMockRecorder) VerifySignature(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySignature", reflect.TypeOf((*MockEth1Client)(nil).VerifySignature), ctx, storeName, request)
}

// VerifyTypedDataSignature mocks base method
func (m *MockEth1Client) VerifyTypedDataSignature(ctx context.Context, storeName string, request *types.VerifyTypedDataRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyTypedDataSignature", ctx, storeName, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyTypedDataSignature indicates an expected call of VerifyTypedDataSignature
func (mr *MockEth1ClientMockRecorder) VerifyTypedDataSignature(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyTypedDataSignature", reflect.TypeOf((*MockEth1Client)(nil).VerifyTypedDataSignature), ctx, storeName, request)
}

// MockJSONRPC is a mock of JSONRPC interface
type MockJSONRPC struct {
	ctrl     *gomock.Controller
	recorder *MockJSONRPCMockRecorder
}

// MockJSONRPCMockRecorder is the mock recorder for MockJSONRPC
type MockJSONRPCMockRecorder struct {
	mock *MockJSONRPC
}

// NewMockJSONRPC creates a new mock instance
func NewMockJSONRPC(ctrl *gomock.Controller) *MockJSONRPC {
	mock := &MockJSONRPC{ctrl: ctrl}
	mock.recorder = &MockJSONRPCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJSONRPC) EXPECT() *MockJSONRPCMockRecorder {
	return m.recorder
}

// Call mocks base method
func (m *MockJSONRPC) Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, nodeID, method}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Call", varargs...)
	ret0, _ := ret[0].(*jsonrpc.ResponseMsg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Call indicates an expected call of Call
func (mr *MockJSONRPCMockRecorder) Call(ctx, nodeID, method interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, nodeID, method}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockJSONRPC)(nil).Call), varargs...)
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

// CreateKey mocks base method
func (m *MockKeyManagerClient) CreateKey(ctx context.Context, storeName string, request *types.CreateKeyRequest) (*types.KeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateKey", ctx, storeName, request)
	ret0, _ := ret[0].(*types.KeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateKey indicates an expected call of CreateKey
func (mr *MockKeyManagerClientMockRecorder) CreateKey(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateKey", reflect.TypeOf((*MockKeyManagerClient)(nil).CreateKey), ctx, storeName, request)
}

// ImportKey mocks base method
func (m *MockKeyManagerClient) ImportKey(ctx context.Context, storeName string, request *types.ImportKeyRequest) (*types.KeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportKey", ctx, storeName, request)
	ret0, _ := ret[0].(*types.KeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportKey indicates an expected call of ImportKey
func (mr *MockKeyManagerClientMockRecorder) ImportKey(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportKey", reflect.TypeOf((*MockKeyManagerClient)(nil).ImportKey), ctx, storeName, request)
}

// SignKey mocks base method
func (m *MockKeyManagerClient) SignKey(ctx context.Context, storeName, id string, request *types.SignBase64PayloadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignKey", ctx, storeName, id, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignKey indicates an expected call of SignKey
func (mr *MockKeyManagerClientMockRecorder) SignKey(ctx, storeName, id, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignKey", reflect.TypeOf((*MockKeyManagerClient)(nil).SignKey), ctx, storeName, id, request)
}

// GetKey mocks base method
func (m *MockKeyManagerClient) GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKey", ctx, storeName, id)
	ret0, _ := ret[0].(*types.KeyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKey indicates an expected call of GetKey
func (mr *MockKeyManagerClientMockRecorder) GetKey(ctx, storeName, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockKeyManagerClient)(nil).GetKey), ctx, storeName, id)
}

// ListKeys mocks base method
func (m *MockKeyManagerClient) ListKeys(ctx context.Context, storeName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListKeys", ctx, storeName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListKeys indicates an expected call of ListKeys
func (mr *MockKeyManagerClientMockRecorder) ListKeys(ctx, storeName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListKeys", reflect.TypeOf((*MockKeyManagerClient)(nil).ListKeys), ctx, storeName)
}

// DestroyKey mocks base method
func (m *MockKeyManagerClient) DestroyKey(ctx context.Context, storeName, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyKey", ctx, storeName, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyKey indicates an expected call of DestroyKey
func (mr *MockKeyManagerClientMockRecorder) DestroyKey(ctx, storeName, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyKey", reflect.TypeOf((*MockKeyManagerClient)(nil).DestroyKey), ctx, storeName, id)
}

// CreateEth1Account mocks base method
func (m *MockKeyManagerClient) CreateEth1Account(ctx context.Context, storeName string, request *types.CreateEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEth1Account", ctx, storeName, request)
	ret0, _ := ret[0].(*types.Eth1AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEth1Account indicates an expected call of CreateEth1Account
func (mr *MockKeyManagerClientMockRecorder) CreateEth1Account(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEth1Account", reflect.TypeOf((*MockKeyManagerClient)(nil).CreateEth1Account), ctx, storeName, request)
}

// ImportEth1Account mocks base method
func (m *MockKeyManagerClient) ImportEth1Account(ctx context.Context, storeName string, request *types.ImportEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportEth1Account", ctx, storeName, request)
	ret0, _ := ret[0].(*types.Eth1AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImportEth1Account indicates an expected call of ImportEth1Account
func (mr *MockKeyManagerClientMockRecorder) ImportEth1Account(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportEth1Account", reflect.TypeOf((*MockKeyManagerClient)(nil).ImportEth1Account), ctx, storeName, request)
}

// UpdateEth1Account mocks base method
func (m *MockKeyManagerClient) UpdateEth1Account(ctx context.Context, storeName string, request *types.UpdateEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEth1Account", ctx, storeName, request)
	ret0, _ := ret[0].(*types.Eth1AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEth1Account indicates an expected call of UpdateEth1Account
func (mr *MockKeyManagerClientMockRecorder) UpdateEth1Account(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEth1Account", reflect.TypeOf((*MockKeyManagerClient)(nil).UpdateEth1Account), ctx, storeName, request)
}

// SignEth1 mocks base method
func (m *MockKeyManagerClient) SignEth1(ctx context.Context, storeName string, request *types.SignHexPayloadRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignEth1", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignEth1 indicates an expected call of SignEth1
func (mr *MockKeyManagerClientMockRecorder) SignEth1(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignEth1", reflect.TypeOf((*MockKeyManagerClient)(nil).SignEth1), ctx, storeName, request)
}

// SignTypedData mocks base method
func (m *MockKeyManagerClient) SignTypedData(ctx context.Context, storeName string, request *types.SignTypedDataRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTypedData", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTypedData indicates an expected call of SignTypedData
func (mr *MockKeyManagerClientMockRecorder) SignTypedData(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTypedData", reflect.TypeOf((*MockKeyManagerClient)(nil).SignTypedData), ctx, storeName, request)
}

// SignTransaction mocks base method
func (m *MockKeyManagerClient) SignTransaction(ctx context.Context, storeName string, request *types.SignETHTransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTransaction", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTransaction indicates an expected call of SignTransaction
func (mr *MockKeyManagerClientMockRecorder) SignTransaction(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTransaction", reflect.TypeOf((*MockKeyManagerClient)(nil).SignTransaction), ctx, storeName, request)
}

// SignQuorumPrivateTransaction mocks base method
func (m *MockKeyManagerClient) SignQuorumPrivateTransaction(ctx context.Context, storeName string, request *types.SignQuorumPrivateTransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignQuorumPrivateTransaction", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignQuorumPrivateTransaction indicates an expected call of SignQuorumPrivateTransaction
func (mr *MockKeyManagerClientMockRecorder) SignQuorumPrivateTransaction(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignQuorumPrivateTransaction", reflect.TypeOf((*MockKeyManagerClient)(nil).SignQuorumPrivateTransaction), ctx, storeName, request)
}

// SignEEATransaction mocks base method
func (m *MockKeyManagerClient) SignEEATransaction(ctx context.Context, storeName string, request *types.SignEEATransactionRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignEEATransaction", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignEEATransaction indicates an expected call of SignEEATransaction
func (mr *MockKeyManagerClientMockRecorder) SignEEATransaction(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignEEATransaction", reflect.TypeOf((*MockKeyManagerClient)(nil).SignEEATransaction), ctx, storeName, request)
}

// GetEth1Account mocks base method
func (m *MockKeyManagerClient) GetEth1Account(ctx context.Context, storeName, account string) (*types.Eth1AccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEth1Account", ctx, storeName, account)
	ret0, _ := ret[0].(*types.Eth1AccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEth1Account indicates an expected call of GetEth1Account
func (mr *MockKeyManagerClientMockRecorder) GetEth1Account(ctx, storeName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEth1Account", reflect.TypeOf((*MockKeyManagerClient)(nil).GetEth1Account), ctx, storeName, account)
}

// ListEth1Account mocks base method
func (m *MockKeyManagerClient) ListEth1Account(ctx context.Context, storeName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEth1Account", ctx, storeName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEth1Account indicates an expected call of ListEth1Account
func (mr *MockKeyManagerClientMockRecorder) ListEth1Account(ctx, storeName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEth1Account", reflect.TypeOf((*MockKeyManagerClient)(nil).ListEth1Account), ctx, storeName)
}

// DeleteEth1Account mocks base method
func (m *MockKeyManagerClient) DeleteEth1Account(ctx context.Context, storeName, account string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEth1Account", ctx, storeName, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEth1Account indicates an expected call of DeleteEth1Account
func (mr *MockKeyManagerClientMockRecorder) DeleteEth1Account(ctx, storeName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEth1Account", reflect.TypeOf((*MockKeyManagerClient)(nil).DeleteEth1Account), ctx, storeName, account)
}

// DestroyEth1Account mocks base method
func (m *MockKeyManagerClient) DestroyEth1Account(ctx context.Context, storeName, account string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DestroyEth1Account", ctx, storeName, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// DestroyEth1Account indicates an expected call of DestroyEth1Account
func (mr *MockKeyManagerClientMockRecorder) DestroyEth1Account(ctx, storeName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DestroyEth1Account", reflect.TypeOf((*MockKeyManagerClient)(nil).DestroyEth1Account), ctx, storeName, account)
}

// RestoreEth1Account mocks base method
func (m *MockKeyManagerClient) RestoreEth1Account(ctx context.Context, storeName, account string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreEth1Account", ctx, storeName, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreEth1Account indicates an expected call of RestoreEth1Account
func (mr *MockKeyManagerClientMockRecorder) RestoreEth1Account(ctx, storeName, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreEth1Account", reflect.TypeOf((*MockKeyManagerClient)(nil).RestoreEth1Account), ctx, storeName, account)
}

// ECRecover mocks base method
func (m *MockKeyManagerClient) ECRecover(ctx context.Context, storeName string, request *types.ECRecoverRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ECRecover", ctx, storeName, request)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ECRecover indicates an expected call of ECRecover
func (mr *MockKeyManagerClientMockRecorder) ECRecover(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ECRecover", reflect.TypeOf((*MockKeyManagerClient)(nil).ECRecover), ctx, storeName, request)
}

// VerifySignature mocks base method
func (m *MockKeyManagerClient) VerifySignature(ctx context.Context, storeName string, request *types.VerifyEth1SignatureRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySignature", ctx, storeName, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifySignature indicates an expected call of VerifySignature
func (mr *MockKeyManagerClientMockRecorder) VerifySignature(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySignature", reflect.TypeOf((*MockKeyManagerClient)(nil).VerifySignature), ctx, storeName, request)
}

// VerifyTypedDataSignature mocks base method
func (m *MockKeyManagerClient) VerifyTypedDataSignature(ctx context.Context, storeName string, request *types.VerifyTypedDataRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyTypedDataSignature", ctx, storeName, request)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyTypedDataSignature indicates an expected call of VerifyTypedDataSignature
func (mr *MockKeyManagerClientMockRecorder) VerifyTypedDataSignature(ctx, storeName, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyTypedDataSignature", reflect.TypeOf((*MockKeyManagerClient)(nil).VerifyTypedDataSignature), ctx, storeName, request)
}

// Call mocks base method
func (m *MockKeyManagerClient) Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, nodeID, method}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Call", varargs...)
	ret0, _ := ret[0].(*jsonrpc.ResponseMsg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Call indicates an expected call of Call
func (mr *MockKeyManagerClientMockRecorder) Call(ctx, nodeID, method interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, nodeID, method}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockKeyManagerClient)(nil).Call), varargs...)
}
