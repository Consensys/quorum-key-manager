package jsonrpc

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)

	client := NewClient(&http.Client{Transport: transport})

	req, _ := http.NewRequest(http.MethodPost, "www.example.com", nil)
	ctx := WithRequest(context.Background(), NewRequest(req))

	// Empty ID and version client
	cllr := NewCaller(WithVersion("")(WithIncrementalID(nil)(client)), nil)

	m := testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"testMethod","params":[1,2,3],"id":"1"}`),
	)
	respBody := []byte(`{"jsonrpc": "1.0","id": "1", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	resp, err := cllr.Call(ctx, "testMethod", []int{1, 2, 3})
	require.Error(t, err, "Call should error")
	assert.Equal(t, "test error message", err.Error(), "Call should error")
	require.IsType(t, new(ErrorMsg), err, "Error should have correct type")
	assert.Equal(t, -32600, err.(*ErrorMsg).Code, "Error code should be correct")
	assert.Equal(t, "1.0", resp.Version(), "Version should be correct")

	// Non Empty ID client caller
	cllr = NewCaller(WithVersion("3.0")(WithIncrementalID("abcd")(client)), nil)

	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"3.0","method":"testMethod","params":[1,2,3],"id":"abcd.1"}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0","id": "abcd.1", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	resp, err = cllr.Call(ctx, "testMethod", []int{1, 2, 3})
	require.Error(t, err, "Call should error")
	assert.Equal(t, "test error message", err.Error(), "Call should error")
	require.IsType(t, new(ErrorMsg), err, "Error should have correct type")
	assert.Equal(t, -32600, err.(*ErrorMsg).Code, "Error code should be correct")
	assert.Equal(t, "1.0", resp.Version(), "Version should be correct")
}
