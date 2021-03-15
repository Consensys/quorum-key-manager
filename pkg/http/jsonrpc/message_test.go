package jsonrpc

import (
	"encoding/json"
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
				assert.Equal(t, paramsV.Interface(), msg.Params, "Params should have been set correctly")
			} else {
				assert.Nil(t, msg.Params, "Params should be nil")
				assert.Nil(t, msg.raw.Params, "Raw Params should be zero")
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
			msg:          &RequestMsg{Version: "2.0", Method: "testMethod", ID: int(0), Params: &TestParams{Field: "test-value"}},
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

		expectedErrMsg string
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
			desc:           "invalid request with no method",
			msg:            &RequestMsg{Version: "2.0", ID: 0},
			expectedErrMsg: "missing method",
		},
		{
			desc:           "invalid request with no version",
			msg:            &RequestMsg{Method: "testMethod", ID: 0},
			expectedErrMsg: "missing version",
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
			desc:           "invalid request with array id",
			msg:            &RequestMsg{Version: "2.0", Method: "testMethod", ID: []int{25}},
			expectedErrMsg: "invalid id (should be int or string but got []int)",
		},
		{
			desc:           "valid request with object id",
			msg:            &RequestMsg{Version: "2.0", Method: "testMethod", ID: &TestID{Field: "test-value"}},
			expectedErrMsg: "invalid id (should be int or string but got jsonrpc.TestID)",
		},
		{
			desc:           "invalid request with invalid json.RawMessage array id",
			msg:            &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`[1,2,3]`)},
			expectedErrMsg: "invalid id [1,2,3]",
		},
		{
			desc:           "invalid request with invalid json.RawMessage object id",
			msg:            &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`{"test-field":"test-value"}`)},
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
			expectedBody: []byte(`{"jsonrpc":"2.0","id":0,"result":{"test-field":"test-value"},"error":{"code":-32600,"message":"test message","data":{"test-field":"test-value"}}}`),
		},
		{
			desc:         "response with no fields",
			msg:          &ResponseMsg{},
			expectedBody: []byte(`{"jsonrpc":"","id":null,"result":null,"error":null}`),
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

func TestUnmarshalErrorMsg(t *testing.T) {
	type TestData struct {
		Field1 string `json:"test-field1"`
		Field2 []int  `json:"test-field2"`
	}

	tests := []struct {
		desc string

		// JSON body of the response
		body []byte

		expectedCode    int
		expectedMessage string
		expectedData    interface{}
	}{
		{
			desc:            "error with all fields",
			body:            []byte(`{"code": -32600, "message":"test error", "data": {"test-field1": "test-value"}}`),
			expectedCode:    -32600,
			expectedMessage: "test error",
			expectedData:    TestData{Field1: "test-value"},
		},
		{
			desc:            "error with null data",
			body:            []byte(`{"code": -32600, "message":"test error", "data": null}`),
			expectedCode:    -32600,
			expectedMessage: "test error",
			expectedData:    nil,
		},
		{
			desc:            "error without fields",
			body:            []byte(`{}`),
			expectedCode:    0,
			expectedMessage: "",
			expectedData:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(ErrorMsg)
			err := json.Unmarshal(tt.body, msg)
			require.NoError(t, err, "Unmarshalling should not fail")

			// Assert unmarshal values
			assert.Equal(t, tt.expectedCode, msg.Code, "Code should be correct")
			assert.Equal(t, tt.expectedMessage, msg.Message, "Message should be correct")

			if tt.expectedData != nil {
				dataV := reflect.New(reflect.TypeOf(tt.expectedData))
				err = msg.UnmarshalData(dataV.Interface())
				require.NoError(t, err, "UnmarshalData should not error")
				assert.Equal(t, tt.expectedData, dataV.Elem().Interface(), "Data should be correct")
				assert.Equal(t, dataV.Interface(), msg.Data, "Data should have been set correctly")
			} else {
				assert.Nil(t, msg.Data, "Data should be nil")
				assert.Nil(t, msg.raw.Data, "Raw Data should be zero")
			}
		})
	}
}

func TestMarshalErrorMsg(t *testing.T) {
	type TestData struct {
		Field1 string `json:"test-field1,omitempty"`
		Field2 []int  `json:"test-field2,omitempty"`
	}

	tests := []struct {
		desc string

		// JSON body of the request
		msg *ErrorMsg

		expectedBody   []byte
		expectedErrMsg string
	}{
		{
			desc: "error with all fields",
			msg: &ErrorMsg{
				Code:    -32600,
				Message: "test message",
				Data:    &TestData{Field1: "test-value"},
			},
			expectedBody: []byte(`{"code":-32600,"message":"test message","data":{"test-field1":"test-value"}}`),
		},
		{
			desc:         "error with no fields",
			msg:          &ErrorMsg{},
			expectedBody: []byte(`{"code":0,"message":"","data":null}`),
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
			desc:         "all fiels request",
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
			expectedBody: []byte(`{"jsonrpc":"","id":null,"result":null,"error":null}`),
		},
		{
			desc:         "all fields response",
			body:         []byte(`{"jsonrpc":"2.0","result":true,"error":{"code":-32600,"message":"test message","data":{"test-field":"test-value"}},"id":"abcd"}`),
			expectedBody: []byte(`{"jsonrpc":"2.0","id":"abcd","result":true,"error":{"code":-32600,"message":"test message","data":{"test-field":"test-value"}}}`),
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

func TestErrorMsgMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		body         []byte
		expectedBody []byte
	}{
		{
			desc:         "empty error",
			body:         []byte(`{}`),
			expectedBody: []byte(`{"code":0,"message":"","data":null}`),
		},
		{
			desc:         "all fields error",
			body:         []byte(`{"code":-32600,"message":"test message","data":true}`),
			expectedBody: []byte(`{"code":-32600,"message":"test message","data":true}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(ErrorMsg)
			err := json.Unmarshal(tt.body, msg)
			require.NoError(t, err, "Unmarshalling should not fail")

			b, err := json.Marshal(msg)
			require.NoError(t, err, "Marshal should not fail")
			assert.Equal(t, tt.expectedBody, b, "Body should match")
		})
	}
}
