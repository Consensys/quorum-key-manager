package proxynode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPCNodeHTTP(t *testing.T) {
	rpcServer := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// Create ResponseWriter
			rpcRw := jsonrpc.NewResponseWriter(rw)

			// Parse request body
			msg := new(jsonrpc.RequestMsg)
			err := json.NewDecoder(req.Body).Decode(msg)
			req.Body.Close()
			if err != nil {
				_ = jsonrpc.WriteError(rpcRw, jsonrpc.ParseError(err))
				return
			}

			jsonrpc.DefaultRWHandler(jsonrpc.HandlerFunc(func(rpcRw jsonrpc.ResponseWriter, msg *jsonrpc.RequestMsg) {
				_ = jsonrpc.WriteResult(rpcRw, msg.Params)
			})).ServeRPC(rpcRw, msg)

		}),
	)
	defer rpcServer.Close()

	cfg := (&Config{
		RPC: &DownstreamConfig{
			Addr: rpcServer.URL,
		},
	}).SetDefault()

	n, err := New(cfg)
	require.NoError(t, err, "New must not error")

	// Test Node
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	msg := new(jsonrpc.RequestMsg).WithVersion("2.0").WithMethod("testMethod").WithID("test-id").WithParams("test-message")
	_ = request.WriteJSON(req, msg)

	rec := httptest.NewRecorder()
	n.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "StatusCode should be OK")
	expectedRespBody := []byte(`{"jsonrpc":"2.0","result":"test-message","error":null,"id":"test-id"}`)
	assert.Equal(t, expectedRespBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "WriteMsg should write correct body")
}

func TestCustomTesseraHandler(t *testing.T) {
	privTxMngrServer := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			_, _ = rw.Write([]byte(`{"key":"q80="}`))
		}),
	)
	defer privTxMngrServer.Close()

	cfg := (&Config{
		PrivTxManager: &DownstreamConfig{
			Addr: privTxMngrServer.URL,
		},
	}).SetDefault()
	b, _ := json.Marshal(cfg)
	t.Logf(string(b))

	n, err := New(cfg)
	require.NoError(t, err, "New must not error")

	n.Handler = jsonrpc.DefaultRWHandler(
		jsonrpc.HandlerFunc(func(rw jsonrpc.ResponseWriter, msg *jsonrpc.RequestMsg) {
			key, e := SessionFromContext(msg.Context()).ClientPrivTxManager().StoreRaw(context.Background(), []byte{}, "")
			if err != nil {
				_ = jsonrpc.WriteError(rw, e)
			} else {
				_ = jsonrpc.WriteResult(rw, key)
			}
		}),
	)

	// Test Node
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	msg := new(jsonrpc.RequestMsg).WithVersion("2.0").WithMethod("testMethod").WithID("test-id").WithParams("test-message")
	_ = request.WriteJSON(req, msg)

	rec := httptest.NewRecorder()
	n.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "StatusCode should be OK")
	expectedRespBody := []byte(`{"jsonrpc":"2.0","result":"q80=","error":null,"id":"test-id"}`)
	assert.Equal(t, expectedRespBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "WriteMsg should write correct body")
}
