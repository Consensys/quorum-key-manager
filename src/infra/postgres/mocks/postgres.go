// Code generated by MockGen. DO NOT EDIT.
// Source: postgres.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	postgres "github.com/consensys/quorum-key-manager/src/infra/postgres"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreateTable mocks base method.
func (m *MockClient) CreateTable(ctx context.Context, model interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTable", ctx, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTable indicates an expected call of CreateTable.
func (mr *MockClientMockRecorder) CreateTable(ctx, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTable", reflect.TypeOf((*MockClient)(nil).CreateTable), ctx, model)
}

// DeletePK mocks base method.
func (m *MockClient) DeletePK(ctx context.Context, model ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range model {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeletePK", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePK indicates an expected call of DeletePK.
func (mr *MockClientMockRecorder) DeletePK(ctx interface{}, model ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, model...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePK", reflect.TypeOf((*MockClient)(nil).DeletePK), varargs...)
}

// DropTable mocks base method.
func (m *MockClient) DropTable(ctx context.Context, model interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DropTable", ctx, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// DropTable indicates an expected call of DropTable.
func (mr *MockClientMockRecorder) DropTable(ctx, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropTable", reflect.TypeOf((*MockClient)(nil).DropTable), ctx, model)
}

// ForceDeletePK mocks base method.
func (m *MockClient) ForceDeletePK(ctx context.Context, model ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range model {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ForceDeletePK", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForceDeletePK indicates an expected call of ForceDeletePK.
func (mr *MockClientMockRecorder) ForceDeletePK(ctx interface{}, model ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, model...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForceDeletePK", reflect.TypeOf((*MockClient)(nil).ForceDeletePK), varargs...)
}

// Insert mocks base method.
func (m *MockClient) Insert(ctx context.Context, model ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range model {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Insert", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockClientMockRecorder) Insert(ctx interface{}, model ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, model...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockClient)(nil).Insert), varargs...)
}

// Ping mocks base method.
func (m *MockClient) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockClientMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockClient)(nil).Ping), ctx)
}

// RunInTransaction mocks base method.
func (m *MockClient) RunInTransaction(ctx context.Context, persist func(postgres.Client) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTransaction", ctx, persist)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTransaction indicates an expected call of RunInTransaction.
func (mr *MockClientMockRecorder) RunInTransaction(ctx, persist interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTransaction", reflect.TypeOf((*MockClient)(nil).RunInTransaction), ctx, persist)
}

// Select mocks base method.
func (m *MockClient) Select(ctx context.Context, model ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range model {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Select", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Select indicates an expected call of Select.
func (mr *MockClientMockRecorder) Select(ctx interface{}, model ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, model...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockClient)(nil).Select), varargs...)
}

// SelectDeleted mocks base method.
func (m *MockClient) SelectDeleted(ctx context.Context, model ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range model {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectDeleted", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectDeleted indicates an expected call of SelectDeleted.
func (mr *MockClientMockRecorder) SelectDeleted(ctx interface{}, model ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, model...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectDeleted", reflect.TypeOf((*MockClient)(nil).SelectDeleted), varargs...)
}

// SelectDeletedPK mocks base method.
func (m *MockClient) SelectDeletedPK(ctx context.Context, model ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range model {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectDeletedPK", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectDeletedPK indicates an expected call of SelectDeletedPK.
func (mr *MockClientMockRecorder) SelectDeletedPK(ctx interface{}, model ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, model...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectDeletedPK", reflect.TypeOf((*MockClient)(nil).SelectDeletedPK), varargs...)
}

// SelectMany mocks base method.
func (m *MockClient) SelectMany(ctx context.Context, model, dst interface{}, condition string, params ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, model, dst, condition}
	for _, a := range params {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectMany", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectMany indicates an expected call of SelectMany.
func (mr *MockClientMockRecorder) SelectMany(ctx, model, dst, condition interface{}, params ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, model, dst, condition}, params...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectMany", reflect.TypeOf((*MockClient)(nil).SelectMany), varargs...)
}

// SelectPK mocks base method.
func (m *MockClient) SelectPK(ctx context.Context, model ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range model {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectPK", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectPK indicates an expected call of SelectPK.
func (mr *MockClientMockRecorder) SelectPK(ctx interface{}, model ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, model...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectPK", reflect.TypeOf((*MockClient)(nil).SelectPK), varargs...)
}

// UpdatePK mocks base method.
func (m *MockClient) UpdatePK(ctx context.Context, model ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range model {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdatePK", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePK indicates an expected call of UpdatePK.
func (mr *MockClientMockRecorder) UpdatePK(ctx interface{}, model ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, model...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePK", reflect.TypeOf((*MockClient)(nil).UpdatePK), varargs...)
}
