// Code generated by MockGen. DO NOT EDIT.
// Source: alias_parser.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	aliases "github.com/consensys/quorum-key-manager/src/aliases"
	gomock "github.com/golang/mock/gomock"
)

// MockAliasParser is a mock of AliasParser interface.
type MockAliasParser struct {
	ctrl     *gomock.Controller
	recorder *MockAliasParserMockRecorder
}

// MockAliasParserMockRecorder is the mock recorder for MockAliasParser.
type MockAliasParserMockRecorder struct {
	mock *MockAliasParser
}

// NewMockAliasParser creates a new mock instance.
func NewMockAliasParser(ctrl *gomock.Controller) *MockAliasParser {
	mock := &MockAliasParser{ctrl: ctrl}
	mock.recorder = &MockAliasParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAliasParser) EXPECT() *MockAliasParserMockRecorder {
	return m.recorder
}

// ParseAlias mocks base method.
func (m *MockAliasParser) ParseAlias(alias string) (string, string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseAlias", alias)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// ParseAlias indicates an expected call of ParseAlias.
func (mr *MockAliasParserMockRecorder) ParseAlias(alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseAlias", reflect.TypeOf((*MockAliasParser)(nil).ParseAlias), alias)
}

// ReplaceAliases mocks base method.
func (m *MockAliasParser) ReplaceAliases(ctx context.Context, aliasBackend aliases.AliasBackend, addrs []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceAliases", ctx, aliasBackend, addrs)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceAliases indicates an expected call of ReplaceAliases.
func (mr *MockAliasParserMockRecorder) ReplaceAliases(ctx, aliasBackend, addrs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceAliases", reflect.TypeOf((*MockAliasParser)(nil).ReplaceAliases), ctx, aliasBackend, addrs)
}
