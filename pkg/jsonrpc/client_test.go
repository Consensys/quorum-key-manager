package jsonrpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/consensysquorum/quorum-key-manager/pkg/http/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithIncrementalID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)

	nilIDClient := WithIncrementalID(nil)(NewHTTPClient(&http.Client{Transport: transport}))
	stringIDClient := WithIncrementalID("abcd")(NewHTTPClient(&http.Client{Transport: transport}))
	intIDClient := WithIncrementalID(47)(NewHTTPClient(&http.Client{Transport: transport}))

	tests := []struct {
		desc string

		client Client

		reqMsg          *RequestMsg
		expectedReqBody []byte
	}{
		{
			desc:   "Base nil - ID specified",
			client: nilIDClient,
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
		},
		{
			desc:   "Base nil- ID unspecified #1",
			client: nilIDClient,
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"1"}`),
		},
		{
			desc:   "Base nil - ID unspecified #2",
			client: nilIDClient,
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"2"}`),
		},
		{
			desc:   "Base string- ID unspecified #1",
			client: stringIDClient,
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"abcd.1"}`),
		},
		{
			desc:   "Base string - ID unspecified #2",
			client: stringIDClient,
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"abcd.2"}`),
		},
		{
			desc:   "Base int- ID unspecified #1",
			client: intIDClient,
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"47.1"}`),
		},
		{
			desc:   "Base int - ID unspecified #2",
			client: intIDClient,
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"47.2"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if len(tt.expectedReqBody) > 0 {
				m := testutils.RequestMatcher(t, "", tt.expectedReqBody)

				header := make(http.Header)
				header.Set("Content-Type", "application/json")

				transport.EXPECT().RoundTrip(m).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"jsonrpc": "2.0","id":1, "result":"abcd"}`))),
					Header:     header,
				}, nil)
			}

			_, err := tt.client.Do(tt.reqMsg)
			require.NoError(t, err, "Do must not error")
		})
	}
}

func TestWithVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)

	defaultVersionClient := WithVersion("")(NewHTTPClient(&http.Client{Transport: transport}))
	versionClient := WithVersion("2.1")(NewHTTPClient(&http.Client{Transport: transport}))

	tests := []struct {
		desc string

		client Client

		reqMsg          *RequestMsg
		expectedReqBody []byte
	}{
		{
			desc:   "Default version - Version specified",
			client: defaultVersionClient,
			reqMsg: &RequestMsg{
				Version: "3.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"3.0","method":"testMethod","params":[1,2,3],"id":1}`),
		},
		{
			desc:   "Default version - Version non specified",
			client: defaultVersionClient,
			reqMsg: &RequestMsg{
				Method: "testMethod",
				Params: []int{1, 2, 3},
				ID:     1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
		},
		{
			desc:   "Custom version - Version non specified",
			client: versionClient,
			reqMsg: &RequestMsg{
				Method: "testMethod",
				Params: []int{1, 2, 3},
				ID:     1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.1","method":"testMethod","params":[1,2,3],"id":1}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if len(tt.expectedReqBody) > 0 {
				m := testutils.RequestMatcher(t, "", tt.expectedReqBody)

				header := make(http.Header)
				header.Set("Content-Type", "application/json")

				transport.EXPECT().RoundTrip(m).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"jsonrpc": "2.0","id":1, "result":"abcd"}`))),
					Header:     header,
				}, nil)
			}

			_, err := tt.client.Do(tt.reqMsg)
			require.NoError(t, err, "Do must not error")
		})
	}
}

func TestValidateID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)

	client := ValidateID(NewHTTPClient(&http.Client{Transport: transport}))

	tests := []struct {
		desc string

		reqMsg          *RequestMsg
		expectedReqBody []byte

		respBody []byte

		expectedError error
	}{
		{
			desc: "Same ID",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respBody:        []byte(`{"jsonrpc": "2.0","id":1, "result":1}`),
		},
		{
			desc: "Distinct ID",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respBody:        []byte(`{"jsonrpc": "2.0","id":5, "result":5}`),
			expectedError: &ErrorMsg{
				Code:    -32003,
				Message: "Invalid downstream JSON-RPC response",
				Data: map[string]interface{}{
					"message": "response id does not match request id",
				},
			},
		},
		{
			desc: "Distinct type ID",
			reqMsg: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedReqBody: []byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":1}`),
			respBody:        []byte(`{"jsonrpc": "2.0","id":"1", "result":5}`),
			expectedError: &ErrorMsg{
				Code:    -32003,
				Message: "Invalid downstream JSON-RPC response",
				Data: map[string]interface{}{
					"message": "json: cannot unmarshal string into Go value of type int",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			m := testutils.RequestMatcher(t, "", tt.expectedReqBody)
			header := make(http.Header)
			header.Set("Content-Type", "application/json")

			transport.EXPECT().RoundTrip(m).Return(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(tt.respBody)),
				Header:     header,
			}, nil)

			_, err := client.Do(tt.reqMsg)
			if tt.expectedError != nil {
				require.Error(t, err, "Do must error")
				assert.Equal(t, tt.expectedError, err, "Error must be valid")
			} else {
				require.NoError(t, err, "Do must not error")
			}
		})
	}
}
