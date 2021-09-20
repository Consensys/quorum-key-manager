package jsonrpc

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/http/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvideCaller(t *testing.T) {
	type TestParam struct {
		Value string `json:"value"`
	}

	type TestResult struct {
		Value string `json:"value"`
	}

	type TestService struct {
		CtxInput_NoOutput        func(Client) func(context.Context)                                  // nolint
		NoInput_NoOutput         func(Client) func()                                                 // nolint
		NonCtxInput_NoOutput     func(Client) func(int)                                              // nolint
		MultiInput_NoOutput      func(Client) func(context.Context, int, string)                     // nolint
		NoInput_ErrorOutput      func(Client) func() error                                           // nolint
		NoInput_IntOutput        func(Client) func() int                                             // nolint
		NoInput_IntErrorOutput   func(Client) func() (int, error)                                    // nolint
		StructInput_StructOutput func(Client) func(context.Context, *TestParam) (*TestResult, error) // nolint
		AllTags                  func(Client) func()                                                 `method:"testMethod" namespace:"eth"`
		MethodTag                func(Client) func()                                                 `method:"testMethod"`
		NamespaceTag             func(Client) func()                                                 `namespace:"eth"`
		ObjectTag                func(Client) func(*TestParam)                                       `object:"-"`
		ByteInput_ByteOutput     func(Client) func([]byte) []byte                                    // nolint
		SliceInput_SliceOutput   func(Client) func([]string) []int                                   // nolint
		MapInput_MapOutput       func(Client) func(map[string][]string) map[string][]string          // nolint
	}

	srv := new(TestService)
	err := ProvideCaller(srv)
	require.NoError(t, err, "ProvideCaller must not error")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transport := testutils.NewMockRoundTripper(ctrl)
	client := WithVersion("2.0")(NewHTTPClient(&http.Client{Transport: transport}))

	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	// CtxInput_NoOutput
	m := testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"CtxInput_NoOutput","params":[],"id":null}`),
	)
	respBody := []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	srv.CtxInput_NoOutput(client)(context.Background())

	// NoInput_NoOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_NoOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	srv.NoInput_NoOutput(client)()

	// NonCtxInput_NoOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"NonCtxInput_NoOutput","params":[278],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	srv.NonCtxInput_NoOutput(client)(278)

	// MultiInput_NoOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"MultiInput_NoOutput","params":[278,"hello world"],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	srv.MultiInput_NoOutput(client)(context.Background(), 278, "hello world")

	// NoInput_ErrorOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_ErrorOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	err = srv.NoInput_ErrorOutput(client)()
	require.Error(t, err)

	require.IsType(t, new(ErrorMsg), err)
	require.Equal(t, -32600, err.(*ErrorMsg).Code)

	// NoInput_IntOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_IntOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": 45}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	res := srv.NoInput_IntOutput(client)()
	assert.Equal(t, 45, res, "NoInput_IntOutput result should match")

	// NoInput_IntErrorOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_IntErrorOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": 38}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	res, err = srv.NoInput_IntErrorOutput(client)()
	require.NoError(t, err, "NoInput_IntErrorOutput must not error")
	assert.Equal(t, 38, res, "NoInput_IntErrorOutput result should match")

	// NoInput_IntErrorOutput (error output)
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"NoInput_IntErrorOutput","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code":-32601,"message":"Method not found"},"id": "1"}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	_, err = srv.NoInput_IntErrorOutput(client)()
	require.Error(t, err)
	require.Equal(t, err.Error(), "Method not found")

	// StructInput_StructOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"StructInput_StructOutput","params":[{"value":"test-req-value"}],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": {"value":"test-resp-value"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	var res2 *TestResult
	res2, err = srv.StructInput_StructOutput(client)(context.Background(), &TestParam{Value: "test-req-value"})
	require.NoError(t, err, "StructInput_StructOutput must not error")
	assert.Equal(t, "test-resp-value", res2.Value, "StructInput_StructOutput result should match")

	// StructInput_StructOutput passing nil arg
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"StructInput_StructOutput","params":[null],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": {"value":"test-resp-value"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)
	res2, err = srv.StructInput_StructOutput(client)(context.Background(), nil)
	require.NoError(t, err, "StructInput_StructOutput must not error")
	assert.Equal(t, "test-resp-value", res2.Value, "StructInput_StructOutput result should match")

	// AllTags
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"eth_testMethod","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	srv.AllTags(client)()

	// MethodTag
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"testMethod","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	srv.MethodTag(client)()

	// NamespaceTag
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"eth_namespaceTag","params":[],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "error": {"code": -32600, "message":"test error message"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	srv.NamespaceTag(client)()

	// ObjectTag
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"ObjectTag","params":{"value":"test-req-value"},"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result": {"value":"test-resp-value"}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	srv.ObjectTag(client)(&TestParam{Value: "test-req-value"})

	// ByteInput_ByteOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"ByteInput_ByteOutput","params":["YWJjZA=="],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result":"ZWZnaA=="}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	b := srv.ByteInput_ByteOutput(client)([]byte("abcd"))
	assert.Equal(t, "efgh", string(b), "ByteInput_ByteOutput result should be correct")

	// SliceInput_SliceOutput   func(Client) func([]string) []int                                   // nolint
	// MapInput_MapOutput       func(Client) func(map[string][]string) map[string][]string          // nolint

	// SliceInput_SliceOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"SliceInput_SliceOutput","params":[["abcd"]],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result":[5]}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	resInts := srv.SliceInput_SliceOutput(client)([]string{"abcd"})
	assert.Equal(t, []int{5}, resInts, "SliceInput_SliceOutput result should be correct")

	// MapInput_MapOutput
	m = testutils.RequestMatcher(
		t,
		"",
		[]byte(`{"jsonrpc":"2.0","method":"MapInput_MapOutput","params":[{"key1":["value1"]}],"id":null}`),
	)
	respBody = []byte(`{"jsonrpc": "1.0", "result":{"key2":["value2"],"key3":["value3","value4"]}}`)
	transport.EXPECT().RoundTrip(m).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
		Header:     header,
	}, nil)

	mapStrs := srv.MapInput_MapOutput(client)(map[string][]string{"key1": {"value1"}})
	assert.Equal(t, map[string][]string{"key2": {"value2"}, "key3": {"value3", "value4"}}, mapStrs, "SliceInput_SliceOutput result should be correct")
}
