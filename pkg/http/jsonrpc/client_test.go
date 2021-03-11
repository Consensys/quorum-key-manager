package jsonrpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type requestMatcher struct {
	URLPath string
	Body    []byte
}

func (m requestMatcher) Matches(x interface{}) bool {
	req, ok := x.(*http.Request)
	if !ok {
		return false
	}

	b := make([]byte, req.ContentLength-1)
	_, _ = io.ReadFull(req.Body, b)

	fmt.Printf("%v\n", string(b))
	fmt.Printf("%v\n", string(req.URL.Path))
	return req.URL.Path == m.URLPath && bytes.Equal(b, m.Body)

}

func (m requestMatcher) String() string {
	return ""
}

func TestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)
	req, _ := http.NewRequest(http.MethodPost, "www.example.com", nil)
	client := NewClientFromRequest(nil, "3.0").WithTransport(transport).WithRequest(req)

	assert.Equal(t, "3.0", client.Version(), "Versin should be correct")

	m := requestMatcher{
		URLPath: "www.example.com",
		Body:    []byte(`{"jsonrpc":"3.0","method":"testMethod","params":[1,2,3],"id":1}`),
	}

	respBody := []byte(`{"jsonrpc": "1.0","id": "25", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	resp, err := client.Call(context.Background(), "testMethod", []int{1, 2, 3})
	require.Error(t, err, "Call should error")
	assert.Equal(t, "test error message", err.Error(), "Call should error")
	require.IsType(t, new(ErrorMsg), err, "Error should have correct type")
	assert.Equal(t, -32600, err.(*ErrorMsg).Code, "Error code should be correct")
	assert.Equal(t, "1.0", resp.Version(), "Version should be correct")
}
