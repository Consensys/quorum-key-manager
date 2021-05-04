package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)

	client := NewHTTPClient(&http.Client{Transport: transport})

	tests := []struct {
		desc string

		reqMsg          *RequestMsg
		expectedReqBody []byte

		respStatus int
		respBody   []byte
		respErr    error

		expectedError   error
		expectedRespMsg *ResponseMsg
	}{
		{
			desc: "Base scenario",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respStatus:      http.StatusOK,
			respBody:        []byte(`{"jsonrpc": "2.0","id":1, "result":"abcd"}`),
			expectedRespMsg: &ResponseMsg{
				Version: "2.0",
				Result:  "abcd",
				ID:      1,
			},
		},
		{
			desc: "Invalid request",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedError: &ErrorMsg{
				Code:    -32600,
				Message: "Invalid Request",
				Data: map[string]interface{}{
					"message": "missing method",
				},
			},
			expectedRespMsg: nil,
		},
		{
			desc: "Downstream error",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respStatus:      http.StatusNotFound,
			respErr:         io.EOF,
			expectedError: &ErrorMsg{
				Code:    -32000,
				Message: "Downstream error",
				Data: map[string]interface{}{
					"status":  502,
					"message": "Bad Gateway",
				},
			},
		},
		{
			desc: "Invalid downstream HTTP status",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respStatus:      http.StatusNotFound,
			expectedError: &ErrorMsg{
				Code:    -32001,
				Message: "Invalid downstream HTTP status",
				Data: map[string]interface{}{
					"status":  404,
					"message": "Not Found",
				},
			},
		},
		{
			desc: "Invalid downstream response (parse error)",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respStatus:      http.StatusOK,
			respBody:        []byte(`hello world`),
			expectedError: &ErrorMsg{
				Code:    -32003,
				Message: "Invalid downstream JSON-RPC response",
				Data: map[string]interface{}{
					"message": "invalid character 'h' looking for beginning of value",
				},
			},
		},
		{
			desc: "Invalid downstream response (validate error)",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respStatus:      http.StatusOK,
			respBody:        []byte(`{"id":1, "result":"abcd"}`),
			expectedError: &ErrorMsg{
				Code:    -32003,
				Message: "Invalid downstream JSON-RPC response",
				Data: map[string]interface{}{
					"message": "missing version",
				},
			},
		},
		{
			desc: "Invalid downstream response (invalid id)",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respStatus:      http.StatusOK,
			respBody:        []byte(`{"id":1,"jsonrpc":"2.0","id":[1,2,3], "result":"abcd"}`),
			expectedError: &ErrorMsg{
				Code:    -32003,
				Message: "Invalid downstream JSON-RPC response",
				Data: map[string]interface{}{
					"message": "invalid id [1,2,3]",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if len(tt.expectedReqBody) > 0 {
				m := testutils.RequestMatcher(t, "", tt.expectedReqBody)
				if tt.respErr == nil {
					header := make(http.Header)
					header.Set("Content-Type", "application/json")

					transport.EXPECT().RoundTrip(m).Return(&http.Response{
						StatusCode: tt.respStatus,
						Body:       ioutil.NopCloser(bytes.NewReader(tt.respBody)),
						Header:     header,
					}, nil)
				} else {
					transport.EXPECT().RoundTrip(m).Return(nil, tt.respErr)
				}
			}

			resp, err := client.Do(tt.reqMsg)
			if tt.expectedError != nil {
				require.Error(t, err, "Do must error")
				assert.Equal(t, tt.expectedError, err, "Error must be valid")
			} else {
				require.NoError(t, err, "Do must not error")
			}

			if tt.expectedRespMsg != nil {
				require.NotNil(t, resp, "Do must return non nil message")
				expectedRespBody, _ := json.Marshal(tt.expectedRespMsg)
				bodyResp, _ := json.Marshal(resp)
				assert.Equal(t, expectedRespBody, bodyResp, "Body should match")
			} else {
				require.Nil(t, resp, "Do must return non nil message")
			}
		})
	}
}
