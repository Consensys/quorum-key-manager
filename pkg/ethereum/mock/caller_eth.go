// Code generated by MockGen. DO NOT EDIT.
// Source: caller_eth.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	big "math/big"
	reflect "reflect"

	ethereum "github.com/consensysquorum/quorum-key-manager/pkg/ethereum"
	common "github.com/ethereum/go-ethereum/common"
	gomock "github.com/golang/mock/gomock"
)

// MockEthCaller is a mock of EthCaller interface.
type MockEthCaller struct {
	ctrl     *gomock.Controller
	recorder *MockEthCallerMockRecorder
}

// MockEthCallerMockRecorder is the mock recorder for MockEthCaller.
type MockEthCallerMockRecorder struct {
	mock *MockEthCaller
}

// NewMockEthCaller creates a new mock instance.
func NewMockEthCaller(ctrl *gomock.Controller) *MockEthCaller {
	mock := &MockEthCaller{ctrl: ctrl}
	mock.recorder = &MockEthCallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEthCaller) EXPECT() *MockEthCallerMockRecorder {
	return m.recorder
}

// ChainID mocks base method.
func (m *MockEthCaller) ChainID(arg0 context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChainID", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChainID indicates an expected call of ChainID.
func (mr *MockEthCallerMockRecorder) ChainID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChainID", reflect.TypeOf((*MockEthCaller)(nil).ChainID), arg0)
}

// EstimateGas mocks base method.
func (m *MockEthCaller) EstimateGas(arg0 context.Context, arg1 *ethereum.CallMsg) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimateGas", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EstimateGas indicates an expected call of EstimateGas.
func (mr *MockEthCallerMockRecorder) EstimateGas(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimateGas", reflect.TypeOf((*MockEthCaller)(nil).EstimateGas), arg0, arg1)
}

// GasPrice mocks base method.
func (m *MockEthCaller) GasPrice(arg0 context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GasPrice", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GasPrice indicates an expected call of GasPrice.
func (mr *MockEthCallerMockRecorder) GasPrice(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GasPrice", reflect.TypeOf((*MockEthCaller)(nil).GasPrice), arg0)
}

// GetTransactionCount mocks base method.
func (m *MockEthCaller) GetTransactionCount(arg0 context.Context, arg1 common.Address, arg2 ethereum.BlockNumber) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionCount indicates an expected call of GetTransactionCount.
func (mr *MockEthCallerMockRecorder) GetTransactionCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionCount", reflect.TypeOf((*MockEthCaller)(nil).GetTransactionCount), arg0, arg1, arg2)
}

// SendRawPrivateTransaction mocks base method.
func (m *MockEthCaller) SendRawPrivateTransaction(arg0 context.Context, arg1 []byte, arg2 *ethereum.PrivateArgs) (common.Hash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRawPrivateTransaction", arg0, arg1, arg2)
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendRawPrivateTransaction indicates an expected call of SendRawPrivateTransaction.
func (mr *MockEthCallerMockRecorder) SendRawPrivateTransaction(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRawPrivateTransaction", reflect.TypeOf((*MockEthCaller)(nil).SendRawPrivateTransaction), arg0, arg1, arg2)
}

// SendRawTransaction mocks base method.
func (m *MockEthCaller) SendRawTransaction(arg0 context.Context, arg1 []byte) (common.Hash, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRawTransaction", arg0, arg1)
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendRawTransaction indicates an expected call of SendRawTransaction.
func (mr *MockEthCallerMockRecorder) SendRawTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRawTransaction", reflect.TypeOf((*MockEthCaller)(nil).SendRawTransaction), arg0, arg1)
}
