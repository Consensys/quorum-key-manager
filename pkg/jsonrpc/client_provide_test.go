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

func TestProvide(t *testing.T) {
	type TestParam struct {
		Value string `json:"value"`
	}

	type TestResult struct {
		Value string `json:"value"`
	}

	type TestService struct {
		CtxInput_NoOutput        func(Caller) func(context.Context)                                  // nolint
		NoInput_NoOutput         func(Caller) func()                                                 // nolint
		NonCtxInput_NoOutput     func(Caller) func(int)                                              // nolint
		MultiInput_NoOutput      func(Caller) func(context.Context, int, string)                     // nolint
		NoInput_ErrorOutput      func(Caller) func() error                                           // nolint
		NoInput_IntOutput        func(Caller) func() int                                             // nolint
		NoInput_IntErrorOutput   func(Caller) func() (int, error)                                    // nolint
		StructInput_StructOutput func(Caller) func(context.Context, *TestParam) (*TestResult, error) // nolint
		AllTags                  func(Caller) func()                                                 `method:"testMethod" namespace:"eth"`
		MethodTag                func(Caller) func()                                                 `method:"testMethod"`
		NamespaceTag             func(Caller) func()                                                 `namespace:"eth"`
		ObjectTag                func(Caller) func(*TestParam)                                       `object:"-"`
	}

	srv := new(TestService)
	err := Provide(srv)
	require.NoError(t, err, "Provide must not error")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)
	client := NewClient(&http.Client{Transport: transport})
	req, _ := http.NewRequest(http.MethodPost, "www.example.com", nil)

	// Empty ID and version client
	cllr := NewCaller(WithVersion("2.0")(client), NewRequest(req))

	// CtxInput_NoOutput
	m := testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"CtxInput_NoOutput","params":[],"id":null}`),
	)
	respBody := []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	srv.CtxInput_NoOutput(cllr)(context.Background())

	// NoInput_NoOutput
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_NoOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	srv.NoInput_NoOutput(cllr)()

	// NonCtxInput_NoOutput
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"NonCtxInput_NoOutput","params":[278],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	srv.NonCtxInput_NoOutput(cllr)(278)

	// MultiInput_NoOutput
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"MultiInput_NoOutput","params":[278,"hello world"],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	srv.MultiInput_NoOutput(cllr)(context.Background(), 278, "hello world")

	// NoInput_ErrorOutput
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_ErrorOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	err = srv.NoInput_ErrorOutput(cllr)()
	require.Error(t, err)
	require.IsType(t, err, new(ErrorMsg))
	require.Equal(t, -32600, err.(*ErrorMsg).Code)

	// NoInput_IntOutput
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_IntOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": 45}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	res := srv.NoInput_IntOutput(cllr)()
	assert.Equal(t, 45, res, "NoInput_IntOutput result should match")

	// NoInput_IntErrorOutput
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_IntErrorOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": 38}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	res, err = srv.NoInput_IntErrorOutput(cllr)()
	require.NoError(t, err, "NoInput_IntErrorOutput must not error")
	assert.Equal(t, 38, res, "NoInput_IntErrorOutput result should match")

	// StructInput_StructOutput
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"StructInput_StructOutput","params":[{"value":"test-req-value"}],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": {"value":"test-resp-value"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	var res2 *TestResult
	res2, err = srv.StructInput_StructOutput(cllr)(context.Background(), &TestParam{Value: "test-req-value"})
	require.NoError(t, err, "StructInput_StructOutput must not error")
	assert.Equal(t, "test-resp-value", res2.Value, "StructInput_StructOutput result should match")

	// AllTags
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"eth_testMethod","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	srv.AllTags(cllr)()

	// MethodTag
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"testMethod","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	srv.MethodTag(cllr)()

	// NamespaceTag
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"eth_namespaceTag","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	srv.NamespaceTag(cllr)()

	// ObjectTag
	m = testutils.RequestMatcher(
		t,
		"www.example.com",
		[]byte(`{"jsonrpc":"2.0","method":"ObjectTag","params":{"value":"test-req-value"},"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": {"value":"test-resp-value"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}, nil)

	srv.ObjectTag(cllr)(&TestParam{Value: "test-req-value"})
}
