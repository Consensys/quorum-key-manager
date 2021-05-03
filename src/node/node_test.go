package node

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertResponse(t *testing.T, resp *jsonrpc.Response, expectedVersion string, expectedID, expectedRes interface{}) {
	assert.Equal(t, expectedVersion, resp.Version(), "Version should be correct")

	id := reflect.New(reflect.TypeOf(expectedID))
	err := resp.UnmarshalID(id.Interface())
	require.NoError(t, err, "UnmarshalID must not error")
	assert.Equal(t, expectedID, id.Elem().Interface(), "ID should be correct")

	res := reflect.New(reflect.TypeOf(expectedRes))
	err = resp.UnmarshalResult(res.Interface())
	require.NoError(t, err, "UnmarshalResult must not error")
	assert.Equal(t, expectedRes, res.Elem().Interface(), "Result should be correct")
}

func TestNodeRPC(t *testing.T) {
	rpcServer := httptest.NewServer(jsonrpc.ToHTTPHandler(
		jsonrpc.HandlerFunc(func(rw jsonrpc.ResponseWriter, req *jsonrpc.Request) {
			_ = rw.WriteResult(req.Params())
		}),
	))
	defer rpcServer.Close()

	privTxMngrServer := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			_, _ = rw.Write([]byte(`All good`))
		}),
	)
	defer privTxMngrServer.Close()

	cfg := (&Config{
		RPC: &DownstreamConfig{
			Addr: rpcServer.URL,
		},
		PrivTxManager: &DownstreamConfig{
			Addr: privTxMngrServer.URL,
		},
	}).SetDefault()

	n, err := New(cfg)
	require.NoError(t, err, "New must not error")

	// Test ClientRPC
	req := jsonrpc.NewRequest(&http.Request{Header: make(http.Header)}).WithVersion("2.0").WithMethod("testMethod").WithID("test-id1").WithParams("test-message1")
	resp, err := n.ClientRPC().Do(req)
	require.NoError(t, err, "Do must not error")
	assertResponse(t, resp, "2.0", "test-id1", "test-message1")

	// TestClient ProxyRPC
	req = jsonrpc.NewRequest(&http.Request{Header: make(http.Header)}).WithVersion("2.0").WithMethod("testMethod").WithID("test-id2").WithParams("test-message2")

	rec := httptest.NewRecorder()
	rw := jsonrpc.NewResponseWriter(rec)

	n.ProxyRPC().ServeRPC(rw, req)

	expectedRespBody := []byte(`{"jsonrpc":"2.0","result":"test-message2","error":null,"id":"test-id2"}`)
	assert.Equal(t, expectedRespBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "WriteMsg should write correct body")

	// Test Session
	session, err := n.Session(jsonrpc.NewRequest(&http.Request{Header: make(http.Header)}).WithVersion("1.0").WithID("test-session-id"))
	require.NoError(t, err, "Session must not error")

	resp, err = session.CallerRPC().Call(context.Background(), "testMethod", "test-message3")
	require.NoError(t, err, "Call must not error")

	assertResponse(t, resp, "1.0", "test-session-id.1", "test-message3")

	// Call a second time to ensure id is correctly incremented
	resp, err = session.CallerRPC().Call(context.Background(), "testMethod", "test-message4")
	require.NoError(t, err, "Call must not error")

	assertResponse(t, resp, "1.0", "test-session-id.2", "test-message4")

	session.Close()

	// Test ClientPrivTxManager
	httpResp, err := n.ClientPrivTxManager().Do(&http.Request{Header: make(http.Header)})
	require.NoError(t, err, "Do must not error")

	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, httpResp.Body)
	assert.Equal(t, []byte(`All good`), buf.Bytes(), "Response body should be correct")

	// Test ProxyPrivTxManager
	rec = httptest.NewRecorder()
	n.ProxyPrivTxManager().ServeHTTP(rec, &http.Request{Header: make(http.Header)})
	assert.Equal(t, []byte(`All good`), rec.Body.Bytes(), "ServeHTTP should write correct body")

}
