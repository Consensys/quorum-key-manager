package jsonrpc

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				assert.Equal(t, dataV.Elem().Interface(), msg.Data, "Data should have been set correctly")
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
