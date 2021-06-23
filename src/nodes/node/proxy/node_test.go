package proxynode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/consensysquorum/quorum-key-manager/pkg/log/testutils"

	"github.com/golang/mock/gomock"

	"github.com/consensysquorum/quorum-key-manager/pkg/http/request"
	"github.com/consensysquorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertResponse(t *testing.T, resp *jsonrpc.ResponseMsg, expectedVersion string, expectedID, expectedRes interface{}) {
	assert.Equal(t, expectedVersion, resp.Version, "Version should be correct")

	id := reflect.New(reflect.TypeOf(expectedID))
	err := resp.UnmarshalID(id.Interface())
	require.NoError(t, err, "UnmarshalID must not error")
	assert.Equal(t, expectedID, id.Elem().Interface(), "ID should be correct")

	res := reflect.New(reflect.TypeOf(expectedRes))
	err = resp.UnmarshalResult(res.Interface())
	require.NoError(t, err, "UnmarshalResult must not error")
	assert.Equal(t, expectedRes, res.Elem().Interface(), "Result should be correct")
}

func TestRPCNodeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	n, err := New(cfg, testutils.NewMockLogger(ctrl))
	require.NoError(t, err, "New must not error")

	err = n.Start(context.Background())
	require.NoError(t, err, "Start must not error")
	defer func() { _ = n.Stop(context.Background()) }()

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		http.Error(w, reason.Error(), status)
	},
}

var dialer = &websocket.Dialer{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 30 * time.Second,
}

func TestNodeWebSocket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rpcServer := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			conn, err := upgrader.Upgrade(rw, req, nil)
			if err != nil {
				fmt.Printf("test-server: Upgrade: %v\n", err)
				return
			}

			go func() {
				defer conn.Close()
				for {
					reqMsg := new(jsonrpc.RequestMsg)
					e := conn.ReadJSON(reqMsg)
					b, _ := json.Marshal(reqMsg)
					fmt.Printf("test-server: ReadJSON: msg=%v err=%v\n", string(b), e)
					if e != nil {
						return
					}

					w, err := conn.NextWriter(websocket.TextMessage)
					if err != nil {
						fmt.Printf("test-server: NextWriter: err=%v\n", e)
						continue
					}

					rpcRw := jsonrpc.NewResponseWriter(w)

					jsonrpc.DefaultRWHandler(jsonrpc.HandlerFunc(func(rpcRw jsonrpc.ResponseWriter, msg *jsonrpc.RequestMsg) {
						err = jsonrpc.WriteResult(rpcRw, msg.Params)
					})).ServeRPC(rpcRw, reqMsg)

					w.Close()
				}
			}()
		}),
	)
	defer rpcServer.Close()

	cfg := (&Config{
		RPC: &DownstreamConfig{
			Addr: rpcServer.URL,
		},
	}).SetDefault()

	n, err := New(cfg, testutils.NewMockLogger(ctrl))
	require.NoError(t, err, "New must not error")

	err = n.Start(context.Background())
	require.NoError(t, err, "Start must not error")

	proxySrv := httptest.NewServer(n)
	defer proxySrv.Close()

	proxyAddr := proxySrv.Listener.Addr().String()
	clientConn, _, err := dialer.Dial(fmt.Sprintf("ws://%v", proxyAddr), nil)
	require.NoError(t, err, "Dial must not error")
	defer clientConn.Close()

	reqMsg := new(jsonrpc.RequestMsg).WithVersion("2.0").WithMethod("testMethod").WithID("test-id").WithParams("test-message")
	err = clientConn.WriteJSON(reqMsg)
	require.NoError(t, err, "WriteJSON must not error")

	respMsg := new(jsonrpc.ResponseMsg)
	err = clientConn.ReadJSON(respMsg)
	require.NoError(t, err, "ReadJSON must not error")
	assertResponse(t, respMsg, "2.0", "test-id", "test-message")

	done := make(chan struct{})
	go func() {
		_ = n.Stop(context.Background())
		close(done)
	}()
	err = clientConn.ReadJSON(respMsg)
	require.Error(t, err, "ReadJSON must error")
	<-done
}

func TestCustomTesseraHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	n, err := New(cfg, testutils.NewMockLogger(ctrl))
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

	err = n.Start(context.Background())
	require.NoError(t, err, "Start must not error")
	defer func() { _ = n.Stop(context.Background()) }()

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
