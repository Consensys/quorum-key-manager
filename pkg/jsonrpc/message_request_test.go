package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalRequestMsg(t *testing.T) {
	type TestParams struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		// JSON body of the request
		body []byte

		expectedVersion string
		expectedMethod  string
		expectedID      interface{}
		expectedParams  interface{}
	}{
		{
			desc:            "request with all fields",
			body:            []byte(`{"jsonrpc": "2.0", "id": "25", "method": "testMethod", "params": {"test-field": "test-value"}}`),
			expectedVersion: "2.0",
			expectedID:      "25",
			expectedMethod:  "testMethod",
			expectedParams:  TestParams{Field: "test-value"},
		},
		{
			desc:            "request without null params and null id",
			body:            []byte(`{"jsonrpc": "2.0", "id": null, "method": "testMethod", "params": null}`),
			expectedVersion: "2.0",
			expectedID:      nil,
			expectedMethod:  "testMethod",
			expectedParams:  nil,
		},
		{
			desc:            "request without arguments",
			body:            []byte(`{}`),
			expectedVersion: "",
			expectedID:      nil,
			expectedMethod:  "",
			expectedParams:  nil,
		},
		{
			desc:            "request with version and method",
			body:            []byte(`{"jsonrpc": "2.0", "method": "testMethod"}`),
			expectedVersion: "2.0",
			expectedID:      nil,
			expectedMethod:  "testMethod",
			expectedParams:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(RequestMsg)
			err := json.Unmarshal(tt.body, msg)
			require.NoError(t, err, "Unmarshalling should not fail")

			// Assert unmarshal values
			assert.Equal(t, tt.expectedVersion, msg.Version, "Version should be correct")
			assert.Equal(t, tt.expectedMethod, msg.Method, "Method should be correct")

			if tt.expectedParams != nil {
				paramsV := reflect.New(reflect.TypeOf(tt.expectedParams))
				err = msg.UnmarshalParams(paramsV.Interface())
				require.NoError(t, err, "UnmarshalParams should not error")
				assert.Equal(t, tt.expectedParams, paramsV.Elem().Interface(), "Params should be correct")
				assert.Equal(t, paramsV.Elem().Interface(), msg.Params, "Params should have been set correctly")
			} else {
				assert.Nil(t, msg.Params, "Params should be nil")
				assert.Nil(t, msg.raw.Params, "Raw Params should be zero")
			}

			if tt.expectedID != nil {
				paramsID := reflect.New(reflect.TypeOf(tt.expectedID))
				err = msg.UnmarshalID(paramsID.Interface())
				require.NoError(t, err, "UnmarshalID should not error")
				assert.Equal(t, tt.expectedID, paramsID.Elem().Interface(), "ID should be correct")
				assert.Equal(t, paramsID.Elem().Interface(), msg.ID, "ID should have been set correctly")
			} else {
				assert.Nil(t, msg.ID, "ID should be nil")
				assert.Nil(t, msg.raw.ID, "Raw Params should be zero")
			}
		})
	}
}

func TestMarshalRequestMsg(t *testing.T) {
	type TestParams struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		// JSON body of the request
		msg *RequestMsg

		expectedBody   []byte
		expectedErrMsg string
	}{
		{
			desc:         "request with all fields",
			msg:          &RequestMsg{Version: "2.0", Method: "testMethod", ID: 0, Params: &TestParams{Field: "test-value"}},
			expectedBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":{"test-field":"test-value"},"id":0}`),
		},
		{
			desc:         "request with no id and no params",
			msg:          &RequestMsg{Version: "2.0", Method: "testMethod"},
			expectedBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":null,"id":null}`),
		},
		{
			desc:         "request with no field",
			msg:          &RequestMsg{},
			expectedBody: []byte(`{"jsonrpc":"","method":"","params":null,"id":null}`),
		},
		{
			desc:         "request with slice id",
			msg:          &RequestMsg{Version: "2.0", Method: "testMethod", ID: []int{27}},
			expectedBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":null,"id":[27]}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			b, err := json.Marshal(tt.msg)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "Marshal should not fail")
				assert.Equal(t, tt.expectedBody, b, "Body should match")
			} else {
				require.Error(t, err, "Marshal should fail")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "Error message should match")
			}
		})
	}
}

func TestRequestMsgValidate(t *testing.T) {
	type TestID struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		// JSON body of the request
		msg *RequestMsg

		expectedErr error
	}{
		{
			desc: "valid request with no params",
			msg:  &RequestMsg{Version: "2.0", Method: "testMethod", ID: 0},
		},
		{
			desc: "valid request with params",
			msg:  &RequestMsg{Version: "2.0", Method: "testMethod", ID: 0, Params: json.RawMessage(`{"test-field": "test-value"}`)},
		},
		{
			desc:        "invalid request with no method",
			msg:         &RequestMsg{Version: "2.0", ID: 0},
			expectedErr: InvalidRequest(fmt.Errorf("missing method")),
		},
		{
			desc:        "invalid request with no version",
			msg:         &RequestMsg{Method: "testMethod", ID: 0},
			expectedErr: InvalidRequest(fmt.Errorf("missing version")),
		},
		{
			desc: "valid request with no id",
			msg:  &RequestMsg{Version: "2.0", Method: "testMethod"},
		},
		{
			desc: "valid request with valid json.RawMessage id",
			msg:  &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`"abcd"`)},
		},
		{
			desc:        "invalid request with array id",
			msg:         &RequestMsg{Version: "2.0", Method: "testMethod", ID: []int{25}},
			expectedErr: InvalidRequest(fmt.Errorf("invalid id (should be int or string but got []int)")),
		},
		{
			desc:        "valid request with object id",
			msg:         &RequestMsg{Version: "2.0", Method: "testMethod", ID: &TestID{Field: "test-value"}},
			expectedErr: InvalidRequest(fmt.Errorf("invalid id (should be int or string but got jsonrpc.TestID)")),
		},
		{
			desc:        "invalid request with invalid json.RawMessage array id",
			msg:         &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`[1,2,3]`)},
			expectedErr: InvalidRequest(fmt.Errorf("invalid id [1,2,3]")),
		},
		{
			desc:        "invalid request with invalid json.RawMessage object id",
			msg:         &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`{"test-field":"test-value"}`)},
			expectedErr: InvalidRequest(fmt.Errorf("invalid id {\"test-field\":\"test-value\"}")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.expectedErr == nil {
				require.NoError(t, err, "Validate should not fail")
			} else {
				require.Error(t, err, "Validate should fail")
				assert.Equal(t, tt.expectedErr, err, "Error message should match")
			}
		})
	}
}

func TestRequestMsgMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		body         []byte
		expectedBody []byte
	}{
		{
			desc:         "empty request",
			body:         []byte(`{}`),
			expectedBody: []byte(`{"jsonrpc":"","method":"","params":null,"id":null}`),
		},
		{
			desc:         "all fields request",
			body:         []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"abcd"}`),
			expectedBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"abcd"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(RequestMsg)
			err := json.Unmarshal(tt.body, msg)
			require.NoError(t, err, "Unmarshalling should not fail")

			b, err := json.Marshal(msg)
			require.NoError(t, err, "Marshal should not fail")
			assert.Equal(t, tt.expectedBody, b, "Body should match")
		})
	}
}

func TestCopyRequestMsg(t *testing.T) {
	b := []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"abcd"}`)

	msg := new(RequestMsg)
	err := json.Unmarshal(b, msg)
	require.NoError(t, err, "Unmarshal must not error")

	cpy := msg.Copy()
	assert.Equal(t, msg.Version, cpy.Version, "Version should equal")
	assert.Equal(t, msg.Method, cpy.Method, "Method should equal")
	assert.Equal(t, msg.Params, cpy.Params, "Method should equal")
	assert.Equal(t, msg.ID, cpy.ID, "ID should equal")
	assert.Equal(t, msg.raw.Version, cpy.raw.Version, "raw.Version should equal")
	assert.Equal(t, msg.raw.Method, cpy.raw.Method, "raw.Method should equal")
	assert.Equal(t, msg.raw.Params, cpy.raw.Params, "raw.Params should equal")
	assert.Equal(t, msg.raw.ID, cpy.raw.ID, "raw.ID should equal")
}
