package jsonrpc

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalResponseMsg(t *testing.T) {
	type TestResult struct {
		Field1 string `json:"test-field1"`
		Field2 []int  `json:"test-field2"`
	}

	tests := []struct {
		desc string

		// JSON body of the response
		body []byte

		expectedVersion string
		expectedID      interface{}
		expectedResult  interface{}
		expectedError   *ErrorMsg
	}{
		{
			desc:            "response with all fields",
			body:            []byte(`{"jsonrpc": "2.0", "id": "25", "result": [{"test-field1": "test-value", "test-field2":[1,2,3]}], "error": {"code": -32600, "message":"test error"}}`),
			expectedVersion: "2.0",
			expectedID:      "25",
			expectedResult: []*TestResult{
				&TestResult{
					Field1: "test-value",
					Field2: []int{1, 2, 3},
				},
			},
			expectedError: &ErrorMsg{Code: -32600, Message: "test error", Data: nil},
		},
		{
			desc:            "response with null result, null id and null error",
			body:            []byte(`{"jsonrpc": "2.0", "id": null, "method": "testMethod", "result": null, "error": null}`),
			expectedVersion: "2.0",
			expectedID:      nil,
			expectedResult:  nil,
			expectedError:   nil,
		},
		{
			desc:            "response without fields",
			body:            []byte(`{}`),
			expectedVersion: "",
			expectedID:      nil,
			expectedResult:  nil,
			expectedError:   nil,
		},
		{
			desc:            "response with version",
			body:            []byte(`{"jsonrpc": "2.0"}`),
			expectedVersion: "2.0",
			expectedID:      nil,
			expectedResult:  nil,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(ResponseMsg)
			err := json.Unmarshal(tt.body, msg)
			require.NoError(t, err, "Unmarshalling should not fail")

			// Assert unmarshal values
			assert.Equal(t, tt.expectedVersion, msg.Version, "Version should be correct")

			if tt.expectedResult != nil {
				resultV := reflect.New(reflect.TypeOf(tt.expectedResult))
				err = msg.UnmarshalResult(resultV.Interface())
				require.NoError(t, err, "expectedResult should not error")
				assert.Equal(t, tt.expectedResult, resultV.Elem().Interface(), "Result should be correct")
				assert.Equal(t, resultV.Interface(), msg.Result, "Result should have been set correctly")
			} else {
				assert.Equal(t, nil, msg.Result, "Params should be nil")
				assert.Nil(t, msg.raw.Result, "Raw Params should be zero")
			}

			if tt.expectedID != nil {
				paramsID := reflect.New(reflect.TypeOf(tt.expectedID))
				err = msg.UnmarshalID(paramsID.Interface())
				require.NoError(t, err, "UnmarshalID should not error")
				assert.Equal(t, tt.expectedID, paramsID.Elem().Interface(), "ID should be correct")
				assert.Equal(t, paramsID.Interface(), msg.ID, "ID should have been set correctly")
			} else {
				assert.Nil(t, msg.ID, "ID should be nil")
				assert.Nil(t, msg.raw.ID, "Raw Params should be zero")
			}

			if tt.expectedError != nil {
				require.NotNil(t, msg.raw.Error, "JSON-RPC Error should not be nil")
				errMgs := new(ErrorMsg)
				err = json.Unmarshal(*msg.raw.Error, errMgs)
				require.NoError(t, err, "Unmarshal error should not error")
				assert.Equal(t, tt.expectedError.Code, errMgs.Code, "Error Code should be correct")
				assert.Equal(t, tt.expectedError.Message, errMgs.Message, "Error Message should be correct")
			} else {
				assert.Nil(t, msg.raw.Error, "JSON-RPC Error should not be nil")
				assert.Nil(t, msg.Error, "Error should be nil")
			}
		})
	}
}

func TestMarshalResponseMsg(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		msg *ResponseMsg

		expectedBody   []byte
		expectedErrMsg string
	}{
		{
			desc: "response with all fields",
			msg: &ResponseMsg{
				Version: "2.0",
				ID:      json.RawMessage(`0`),
				Result:  json.RawMessage([]byte(`{"test-field": "test-value"}`)),
				Error: &ErrorMsg{
					Code:    -32600,
					Message: "test message",
					Data:    json.RawMessage(`{"test-field": "test-value"}`),
				}},
			expectedBody: []byte(`{"jsonrpc":"2.0","result":{"test-field":"test-value"},"error":{"code":-32600,"message":"test message","data":{"test-field":"test-value"}},"id":0}`),
		},
		{
			desc:         "response with no fields",
			msg:          &ResponseMsg{},
			expectedBody: []byte(`{"jsonrpc":"","result":null,"error":null,"id":null}`),
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

func TestResponseMsgValidate(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		msg *ResponseMsg

		expectedErrMsg string
	}{
		{
			desc: "valid success response",
			msg:  &ResponseMsg{Version: "2.0", ID: 0, Result: json.RawMessage(`true`)},
		},
		{
			desc: "valid failure response",
			msg:  &ResponseMsg{Version: "2.0", ID: 0, Error: &ErrorMsg{Code: -32600}},
		},
		{
			desc: "valid request with valid json.RawMessage id",
			msg:  &ResponseMsg{Version: "2.0", ID: json.RawMessage(`"abcd"`), Result: true},
		},
		{
			desc:           "invalid success response without result",
			msg:            &ResponseMsg{Version: "2.0", ID: 0},
			expectedErrMsg: "missing result on success",
		},

		{
			desc:           "invalid failure response with result",
			msg:            &ResponseMsg{Version: "2.0", ID: 0, Error: &ErrorMsg{Code: -32600}, Result: json.RawMessage(`true`)},
			expectedErrMsg: "non empty result on failure",
		},
		{
			desc:           "invalid response with no version",
			msg:            &ResponseMsg{ID: 0, Result: true},
			expectedErrMsg: "missing version",
		},
		{
			desc: "valid response with no id",
			msg:  &ResponseMsg{Version: "2.0", Result: true},
		},
		{
			desc:           "invalid response with array id",
			msg:            &ResponseMsg{Version: "2.0", ID: []int{25}, Result: true},
			expectedErrMsg: "invalid id (should be int or string but got []int)",
		},
		{
			desc:           "invalid request with invalid json.RawMessage array id",
			msg:            &ResponseMsg{Version: "2.0", ID: json.RawMessage(`[1,2,3]`), Result: true},
			expectedErrMsg: "invalid id [1,2,3]",
		},
		{
			desc:           "invalid request with invalid json.RawMessage object id",
			msg:            &ResponseMsg{Version: "2.0", ID: json.RawMessage(`{"test-field":"test-value"}`), Result: true},
			expectedErrMsg: "invalid id {\"test-field\":\"test-value\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "Validate should not fail")
			} else {
				require.Error(t, err, "Validate should fail")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "Error message should match")
			}
		})
	}
}

func TestResponseMsgMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		body         []byte
		expectedBody []byte
	}{
		{
			desc:         "empty response",
			body:         []byte(`{}`),
			expectedBody: []byte(`{"jsonrpc":"","result":null,"error":null,"id":null}`),
		},
		{
			desc:         "all fields response",
			body:         []byte(`{"jsonrpc":"2.0","result":true,"error":{"code":-32600,"message":"test message","data":{"test-field":"test-value"}},"id":"abcd"}`),
			expectedBody: []byte(`{"jsonrpc":"2.0","result":true,"error":{"code":-32600,"message":"test message","data":{"test-field":"test-value"}},"id":"abcd"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(ResponseMsg)
			err := json.Unmarshal(tt.body, msg)
			require.NoError(t, err, "Unmarshalling should not fail")

			b, err := json.Marshal(msg)
			require.NoError(t, err, "Marshal should not fail")
			assert.Equal(t, tt.expectedBody, b, "Body should match")
		})
	}
}
