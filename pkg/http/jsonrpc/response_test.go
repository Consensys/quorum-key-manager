package jsonrpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	type TestResult struct {
		Field string `json:"test-field"`
	}

	tests := []struct {
		desc string

		body   []byte
		status int

		expectedReadBodyErrMsg string

		expectedVersion  string
		expectedID       string
		expectedResult   TestResult
		expectedErrorMsg *ErrorMsg
	}{
		{
			desc:                   "200 with valid body",
			body:                   []byte(`{"jsonrpc": "1.0", "id": "25", "result": {"test-field": "test-value", "test-field2":[1,2,3]}, "error": {"code": -32600, "message":"test error"}}`),
			status:                 http.StatusOK,
			expectedReadBodyErrMsg: "",
			expectedVersion:        "1.0",
			expectedID:             "25",
			expectedResult:         TestResult{Field: "test-value"},
			expectedErrorMsg:       &ErrorMsg{Code: -32600, Message: "test error"},
		},
		{
			desc:                   "500 with valid body",
			body:                   []byte(`{"jsonrpc": "1.0", "id": "25", "result": {"test-field": "test-value", "test-field2":[1,2,3]}, "error": {"code": -32600, "message":"test error"}}`),
			status:                 http.StatusInternalServerError,
			expectedReadBodyErrMsg: "invalid http response: Internal Server Error (code=500)",
			expectedVersion:        "",
			expectedID:             "",
			expectedResult:         TestResult{},
			expectedErrorMsg:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Create request
			buf := bytes.NewReader(tt.body)
			resp := NewResponse(&http.Response{
				StatusCode: tt.status,
				Body:       ioutil.NopCloser(buf),
			})

			err := resp.ReadBody()
			if tt.expectedReadBodyErrMsg != "" {
				require.Error(t, err, "ReadBody should error")
				require.Equal(t, tt.expectedReadBodyErrMsg, err.Error(), "ReadBody error should be correct")
			} else {
				require.NoError(t, err, "ReadBody should not error")

				id := new(string)
				err = resp.UnmarshalID(id)
				require.NoError(t, err, "UnmarshalID should not error")
				assert.Equal(t, tt.expectedID, *id, "ID should be correct")
				assert.Equal(t, id, resp.ID(), "ID should have been set correctly")

				// Unmarshal Params
				result := new(TestResult)
				err = resp.UnmarshalResult(result)
				require.NoError(t, err, "UnmarshalResult should not error")
				assert.Equal(t, tt.expectedResult, *result, "Result should be correct")
				assert.Equal(t, result, resp.Result(), "Result should have been set correctly")
			}

			if tt.expectedErrorMsg != nil {
				assert.Equal(t, tt.expectedErrorMsg.Code, resp.Error().(*ErrorMsg).Code, "Error code should be correct")
				assert.Equal(t, tt.expectedErrorMsg.Message, resp.Error().(*ErrorMsg).Message, "Error message be correct")
			} else {
				_, ok := resp.Error().(*ErrorMsg)
				require.False(t, ok, "Error should not be a JSON-RPC")
			}

			assert.Equal(t, tt.expectedVersion, resp.Version(), "Version should be correct")
		})
	}
}
