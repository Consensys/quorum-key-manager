package jsonrpc

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeHandler(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the response
		f interface{}

		expectedErrMsg string

		req          *Request
		expectedBody []byte
	}{
		{
			desc: "Valid - in: (context,int) // out: (int,error) // return result",
			f: func(ctx context.Context, i int) (int, error) {
				return i, nil
			},
			req:          NewRequest(&http.Request{}).WithVersion("2.0").WithMethod("testMethod").WithParams([]int{5}),
			expectedBody: []byte(`{"jsonrpc":"","result":5,"error":null,"id":null}`),
		},
		{
			desc: "Valid - in: (context,int) // out: (int,error) // return error",
			f: func(ctx context.Context, i int) (int, error) {
				return 0, fmt.Errorf("test-error")
			},
			req:          NewRequest(&http.Request{}).WithVersion("2.0").WithMethod("testMethod").WithParams([]int{5}),
			expectedBody: []byte(`{"jsonrpc":"","result":null,"error":{"code":-32000,"message":"test-error","data":null},"id":null}`),
		},
		{
			desc: "Valid - in: (string) // out: string // return result",
			f: func(s string) string {
				return s
			},
			req:          NewRequest(&http.Request{}).WithVersion("2.0").WithMethod("testMethod").WithParams([]string{"hello message"}),
			expectedBody: []byte(`{"jsonrpc":"","result":"hello message","error":null,"id":null}`),
		},
		{
			desc: "Valid - in: () // out: string // return result",
			f: func() string {
				return "hello message"
			},
			req:          NewRequest(&http.Request{}).WithVersion("2.0").WithMethod("testMethod").WithParams([]string{"hello message"}),
			expectedBody: []byte(`{"jsonrpc":"","result":"hello message","error":null,"id":null}`),
		},
		{
			desc:           "Invalid - nil func",
			f:              nil,
			expectedErrMsg: "can not generate handler from zero value",
		},
		{
			desc:           "Invalid - int input",
			f:              int(0),
			expectedErrMsg: "expect function but got int",
		},
		{
			desc:           "Invalid - too many outputs",
			f:              func() (int, string, error) { return 0, "", nil },
			expectedErrMsg: "function must return at most two outputs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			handler, err := MakeHandler(tt.f)
			if tt.expectedErrMsg != "" {
				require.Error(t, err, "MakeHandler must error")
				assert.Equal(t, tt.expectedErrMsg, err.Error(), "MakeHandler error message should be correct")
			} else {
				require.NoError(t, err, "MakeHandler should not fail")
				rec := httptest.NewRecorder()
				rw := NewResponseWriter(rec)

				_ = tt.req.WriteBody()

				handler.ServeRPC(rw, tt.req)
				assert.Equal(t, http.StatusOK, rec.Code, "Code should be correct")
				assert.Equal(t, tt.expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Correct body should have been written")
			}
		})
	}
}
