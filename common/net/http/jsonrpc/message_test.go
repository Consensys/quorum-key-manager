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
	tests := []struct {
		desc string

		// JSON body of the request
		body []byte

		expectedVersion string
		expectedID      json.RawMessage
		expectedMethod  string
		expectedParams  json.RawMessage
	}{
		{
			desc:            "valid request with params",
			body:            []byte(`{"jsonrpc": "2.0", "id": "25", "method": "testMethod", "params": {"test-field": "test-value"}}`),
			expectedVersion: "2.0",
			expectedID:      json.RawMessage([]byte(`"25"`)),
			expectedMethod:  "testMethod",
			expectedParams:  json.RawMessage([]byte(`{"test-field": "test-value"}`)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(RequestMsg)
			err := json.Unmarshal(tt.body, msg)
			require.NoError(t, err, "Unmarshalling should not fail")

			// Assert unmarshal values
			assert.Equal(t, tt.expectedVersion, msg.Version, "Version should be correct")
			assert.Equal(t, tt.expectedID, msg.ID, "ID should be correct")
			assert.Equal(t, tt.expectedMethod, msg.Method, "Method should be correct")
			assert.Equal(t, tt.expectedParams, msg.Params, "Params should be correct")
		})
	}
}

func TestMarshalRequestMsg(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the request
		msg *RequestMsg

		expectedBody   []byte
		expectedErrMsg string
	}{
		{
			desc:         "valid request with no method",
			msg:          &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`0`), Params: json.RawMessage([]byte(`{"test-field": "test-value"}`))},
			expectedBody: []byte(`{"jsonrpc":"2.0","id":0,"method":"testMethod","params":{"test-field":"test-value"}}`),
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
	tests := []struct {
		desc string

		// JSON body of the request
		msg *RequestMsg

		expectedErrMsg string
	}{
		{
			desc: "valid request with no params",
			msg:  &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`0`)},
		},
		{
			desc: "valid request with params",
			msg:  &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`0`), Params: json.RawMessage(`{"test-field": "test-value"}`)},
		},
		{
			desc:           "invalid request with no method",
			msg:            &RequestMsg{Version: "2.0", ID: json.RawMessage(`0`)},
			expectedErrMsg: "missing method",
		},
		{
			desc:           "invalid request with no version",
			msg:            &RequestMsg{Method: "testMethod", ID: json.RawMessage(`0`)},
			expectedErrMsg: "missing version",
		},
		{
			desc:           "invalid request with no id",
			msg:            &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(nil)},
			expectedErrMsg: "missing id",
		},
		{
			desc:           "invalid request with array id",
			msg:            &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage(`["25"]`)},
			expectedErrMsg: "invalid id [\"25\"]",
		},
		{
			desc:           "valid request with object id",
			msg:            &RequestMsg{Version: "2.0", Method: "testMethod", ID: json.RawMessage([]byte(`{"id": "25"}`))},
			expectedErrMsg: "invalid id {\"id\": \"25\"}",
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

func TestRequestMsgUnmarshalParams(t *testing.T) {
	type TestParams struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		// JSON body of the request
		msg *RequestMsg

		expectedParams interface{}
		expectedErrMsg string
	}{
		{
			msg:            &RequestMsg{Params: json.RawMessage(`{"test-field": "test-value"}`)},
			expectedParams: TestParams{Field: "test-value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			val := reflect.New(reflect.TypeOf(tt.expectedParams))
			err := tt.msg.UnmarshalParams(val.Interface())
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "UnmarshalParams should not fail")
				assert.Equal(t, tt.expectedParams, val.Elem().Interface(), "Params should match")
			} else {
				require.Error(t, err, "UnmarshalParams should fail")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "Error message should match")
			}
		})
	}
}

func TestRequestMsgWithID(t *testing.T) {
	type TestID struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		id interface{}

		expectedID     json.RawMessage
		expectedErrMsg string
	}{
		{
			desc:       "int",
			id:         int(27),
			expectedID: json.RawMessage(`27`),
		},
		{
			desc:       "string",
			id:         "abcd",
			expectedID: json.RawMessage(`"abcd"`),
		},
		{
			desc:       "nil",
			id:         nil,
			expectedID: json.RawMessage(`null`),
		},
		{
			desc:           "slice",
			id:             []int{1, 2, 3},
			expectedID:     json.RawMessage(nil),
			expectedErrMsg: "invalid id [1,2,3]",
		},
		{
			desc:           "struct",
			id:             TestID{Field: "test-value"},
			expectedID:     json.RawMessage(nil),
			expectedErrMsg: "invalid id {\"test-field\":\"test-value\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(RequestMsg)

			err := msg.WithID(tt.id)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "WithID should not fail")
			} else {
				require.Error(t, err, "WithID should fail")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "Error message should match")
			}

			assert.Equal(t, tt.expectedID, msg.ID, "ID should match")
		})
	}
}

func TestRequestMsgWithParams(t *testing.T) {
	type TestParams struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		params interface{}

		expectedParams json.RawMessage
		expectedErrMsg string
	}{
		{
			desc:           "int",
			params:         int(27),
			expectedParams: json.RawMessage(`27`),
		},
		{
			desc:           "string",
			params:         "abcd",
			expectedParams: json.RawMessage(`"abcd"`),
		},
		{
			desc:           "nil",
			params:         nil,
			expectedParams: json.RawMessage(`null`),
		},
		{
			desc:           "slice",
			params:         []int{1, 2, 3},
			expectedParams: json.RawMessage(`[1,2,3]`),
		},
		{
			desc:           "struct",
			params:         TestParams{Field: "test-value"},
			expectedParams: json.RawMessage(`{"test-field":"test-value"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(RequestMsg)

			err := msg.WithParams(tt.params)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "UnmarshalParams should not fail")
			} else {
				require.Error(t, err, "UnmarshalParams should fail")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "Error message should match")
			}

			assert.Equal(t, tt.expectedParams, msg.Params, "ID should match")
		})
	}
}

func TestUnmarshalResponseMsg(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the response
		body []byte

		expectedVersion string
		expectedID      json.RawMessage
		expectedResult  json.RawMessage
		expectedError   *ErrorMsg
	}{
		{
			desc:            "valid response with result",
			body:            []byte(`{"jsonrpc": "2.0", "id": "25", "result": {"test-field1": "test-value", "test-field2":[1,2,3]}, "error": {"code": -32600, "message":"test error"}}`),
			expectedVersion: "2.0",
			expectedID:      json.RawMessage([]byte(`"25"`)),
			expectedResult:  json.RawMessage([]byte(`{"test-field1": "test-value", "test-field2":[1,2,3]}`)),
			expectedError:   &ErrorMsg{Code: -32600, Message: "test error", Data: json.RawMessage(nil)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(ResponseMsg)
			err := json.Unmarshal(tt.body, msg)
			require.NoError(t, err, "Unmarshalling should not fail")

			// Assert unmarshal values
			assert.Equal(t, tt.expectedVersion, msg.Version, "Version should be correct")
			assert.Equal(t, tt.expectedID, msg.ID, "ID should be correct")
			assert.Equal(t, tt.expectedResult, msg.Result, "Result should be correct")
			if tt.expectedError != nil {
				assert.Equal(t, *tt.expectedError, *msg.Error, "Error should match")
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
			desc: "valid request with no method",
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
			msg:  &ResponseMsg{Version: "2.0", ID: json.RawMessage(`0`), Result: json.RawMessage(`true`)},
		},
		{
			desc: "valid failure response",
			msg:  &ResponseMsg{Version: "2.0", ID: json.RawMessage(`0`), Error: &ErrorMsg{Code: -32600}},
		},
		{
			desc:           "valid success response without result",
			msg:            &ResponseMsg{Version: "2.0", ID: json.RawMessage(`0`)},
			expectedErrMsg: "missing result on success",
		},
		{
			desc:           "invalid failure response with result",
			msg:            &ResponseMsg{Version: "2.0", ID: json.RawMessage(`0`), Error: &ErrorMsg{Code: -32600}, Result: json.RawMessage(`true`)},
			expectedErrMsg: "non empty result on failure",
		},
		{
			desc:           "invalid response with no version",
			msg:            &ResponseMsg{ID: json.RawMessage(`0`), Result: json.RawMessage(`true`)},
			expectedErrMsg: "missing version",
		},
		{
			desc:           "invalid response with no id",
			msg:            &ResponseMsg{Version: "2.0", ID: json.RawMessage(nil), Result: json.RawMessage(`true`)},
			expectedErrMsg: "missing id",
		},
		{
			desc:           "invalid response with array id",
			msg:            &ResponseMsg{Version: "2.0", ID: json.RawMessage(`["25"]`), Result: json.RawMessage(`true`)},
			expectedErrMsg: "invalid id [\"25\"]",
		},
		{
			desc:           "valid response with object id",
			msg:            &ResponseMsg{Version: "2.0", ID: json.RawMessage([]byte(`{"id": "25"}`)), Result: json.RawMessage(`true`)},
			expectedErrMsg: "invalid id {\"id\": \"25\"}",
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

func TestResponseMsgUnmarshalResult(t *testing.T) {
	type TestResult struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		// JSON body of the request
		msg *ResponseMsg

		expectedResult interface{}
		expectedErrMsg string
	}{
		{
			msg:            &ResponseMsg{Result: json.RawMessage(`{"test-field": "test-value"}`)},
			expectedResult: TestResult{Field: "test-value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			val := reflect.New(reflect.TypeOf(tt.expectedResult))
			err := tt.msg.UnmarshalResult(val.Interface())
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "UnmarshalParams should not fail")
				assert.Equal(t, tt.expectedResult, val.Elem().Interface(), "Result should match")
			} else {
				require.Error(t, err, "UnmarshalParams should fail")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "Error message should match")
			}
		})
	}
}

func TestResponseMessageWithID(t *testing.T) {
	type TestID struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		id interface{}

		expectedID     json.RawMessage
		expectedErrMsg string
	}{
		{
			desc:       "int",
			id:         int(27),
			expectedID: json.RawMessage(`27`),
		},
		{
			desc:       "string",
			id:         "abcd",
			expectedID: json.RawMessage(`"abcd"`),
		},
		{
			desc:       "nil",
			id:         nil,
			expectedID: json.RawMessage(`null`),
		},
		{
			desc:           "slice",
			id:             []int{1, 2, 3},
			expectedID:     json.RawMessage(nil),
			expectedErrMsg: "invalid id [1,2,3]",
		},
		{
			desc:           "struct",
			id:             TestID{Field: "test-value"},
			expectedID:     json.RawMessage(nil),
			expectedErrMsg: "invalid id {\"test-field\":\"test-value\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(ResponseMsg)

			err := msg.WithID(tt.id)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "WithID should not fail")
			} else {
				require.Error(t, err, "WithID should fail")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "Error message should match")
			}

			assert.Equal(t, tt.expectedID, msg.ID, "ID should match")
		})
	}
}

func TestResponseMsgWithResult(t *testing.T) {
	type TestParams struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		result interface{}

		expectedResult json.RawMessage
		expectedErrMsg string
	}{
		{
			desc:           "int",
			result:         int(27),
			expectedResult: json.RawMessage(`27`),
		},
		{
			desc:           "string",
			result:         "abcd",
			expectedResult: json.RawMessage(`"abcd"`),
		},
		{
			desc:           "nil",
			result:         nil,
			expectedResult: json.RawMessage(`null`),
		},
		{
			desc:           "slice",
			result:         []int{1, 2, 3},
			expectedResult: json.RawMessage(`[1,2,3]`),
		},
		{
			desc:           "struct",
			result:         TestParams{Field: "test-value"},
			expectedResult: json.RawMessage(`{"test-field":"test-value"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(ResponseMsg)

			err := msg.WithResult(tt.result)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "WithResult should not fail")
			} else {
				require.Error(t, err, "WithResult should fail")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "Error message should match")
			}

			assert.Equal(t, tt.expectedResult, msg.Result, "ID should match")
		})
	}
}

func TestResponseMsgWithError(t *testing.T) {
	tests := []struct {
		desc string

		err error

		expectedError *ErrorMsg
	}{
		{
			desc:          "any error",
			err:           fmt.Errorf("test error"),
			expectedError: &ErrorMsg{Message: "test error"},
		},
		{
			desc:          "error message",
			err:           &ErrorMsg{Code: -32600, Message: "test error"},
			expectedError: &ErrorMsg{Code: -32600, Message: "test error"},
		},
		{
			desc:          "nil",
			err:           nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(ResponseMsg)

			msg.WithError(tt.err)

			assert.Equal(t, tt.expectedError, msg.Error, "Error should match")
		})
	}
}
