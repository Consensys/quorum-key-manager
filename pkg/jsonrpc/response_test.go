package jsonrpc

import (
	"bytes"
	"fmt"
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
		expectedErrorMsg error
	}{
		{
			desc:                   "code 200 with valid succes body",
			body:                   []byte(`{"jsonrpc": "1.0", "id": "25", "result": {"test-field": "test-value", "test-field2":[1,2,3]}}`),
			status:                 http.StatusOK,
			expectedReadBodyErrMsg: "",
			expectedVersion:        "1.0",
			expectedID:             "25",
			expectedResult:         TestResult{Field: "test-value"},
			expectedErrorMsg:       nil,
		},
		{
			desc:                   "code 200 with valid failure body",
			body:                   []byte(`{"jsonrpc": "1.0", "id": "25",  "error": {"code": -32600, "message":"test error"}}`),
			status:                 http.StatusOK,
			expectedReadBodyErrMsg: "",
			expectedVersion:        "1.0",
			expectedID:             "25",
			expectedResult:         TestResult{},
			expectedErrorMsg:       &ErrorMsg{Code: -32600, Message: "test error"},
		},
		{
			desc:                   "code 500",
			body:                   nil,
			status:                 http.StatusInternalServerError,
			expectedReadBodyErrMsg: "invalid http response: Internal Server Error (code=500)",
			expectedVersion:        "",
			expectedID:             "",
			expectedResult:         TestResult{},
			expectedErrorMsg:       fmt.Errorf("invalid http response: Internal Server Error (code=500)"),
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
				assert.Equal(t, *id, resp.ID(), "ID should have been set correctly")

				// Unmarshal Params
				result := new(TestResult)
				err = resp.UnmarshalResult(result)
				require.NoError(t, err, "UnmarshalResult should not error")
				assert.Equal(t, tt.expectedResult, *result, "Result should be correct")
				assert.Equal(t, *result, resp.Result(), "Result should have been set correctly")
			}

			if tt.expectedErrorMsg != nil {
				require.NotNil(t, resp.Error(), "Resp should have an error")
				expectedMsg, ok := tt.expectedErrorMsg.(*ErrorMsg)
				if ok {
					assert.Equal(t, expectedMsg.Code, resp.Error().(*ErrorMsg).Code, "Error code should be correct")
					assert.Equal(t, expectedMsg.Message, resp.Error().(*ErrorMsg).Message, "Error message be correct")
				} else {
					errMsg, ok := resp.Error().(*ErrorMsg)
					assert.False(t, ok, "Error should not cast into %T", errMsg)
					assert.Equal(t, tt.expectedErrorMsg.Error(), resp.Error().Error(), "Error message shoud match")
				}
			} else {
				assert.Nil(t, resp.Error(), "Resp should not have an error")
			}

			assert.Equal(t, tt.expectedVersion, resp.Version(), "Version should be correct")
		})
	}
}
