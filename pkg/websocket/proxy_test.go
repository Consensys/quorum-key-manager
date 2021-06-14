package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/consensysquorum/quorum-key-manager/pkg/http/request"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/response"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		http.Error(w, reason.Error(), status)
	},
	HandshakeTimeout: 30 * time.Second,
}

var dialer = &websocket.Dialer{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 30 * time.Second,
}

type nextMsg struct {
	typ int
	msg []byte
	err error
}

type bckndHandler struct {
	in  chan *nextMsg
	out chan interface{}

	stop chan struct{}
}

func (h bckndHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("test-server: Upgrade: %v\n", err)
		return
	}

	clientConn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("test-server: Handle close code=%v text=%v\n", code, text)
		message := websocket.FormatCloseMessage(code, text)
		_ = clientConn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	go func() {
		<-h.stop
		clientConn.Close()
	}()

	go func() {
		defer clientConn.Close()
		defer close(h.in)
		for {

			nextMsg := new(nextMsg)
			nextMsg.typ, nextMsg.msg, nextMsg.err = clientConn.ReadMessage()
			fmt.Printf("test-server: ReadMessage: %v %v %v\n", nextMsg.typ, nextMsg.msg, nextMsg.err)
			h.in <- nextMsg
			if nextMsg.err != nil {
				return
			}
		}
	}()

	go func() {
		defer clientConn.Close()
		for out := range h.out {
			select {
			case <-h.stop:
				return
			default:
			}
			err = clientConn.WriteJSON(out)
			fmt.Printf("test-server: WriteJSON: %v\n", err)
			if err != nil {
				return
			}
		}
	}()
}

func createProxyServer(uri string) (*httptest.Server, *Proxy) {
	prep, _ := request.Proxy(&request.ProxyConfig{Addr: uri})
	modif := response.Proxy(&response.ProxyConfig{})
	cfg := (&ProxyConfig{}).SetDefault()

	prx := NewProxy(cfg)
	prx.Upgrader = upgrader
	prx.Dialer = dialer
	prx.ReqPreparer = prep
	prx.RespModifier = modif

	_ = prx.Start(context.TODO())

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		prx.ServeHTTP(w, req)
	}))

	prx.RegisterServerShutdown(srv.Config)

	return srv, prx
}

func TestProxy(t *testing.T) {
	h := bckndHandler{
		in:  make(chan *nextMsg, 3),
		out: make(chan interface{}),
	}

	backSrv := httptest.NewServer(h)
	defer backSrv.Close()

	proxySrv, _ := createProxyServer(backSrv.URL)
	defer proxySrv.Close()

	proxyAddr := proxySrv.Listener.Addr().String()
	clientConn, _, err := dialer.Dial(fmt.Sprintf("ws://%v", proxyAddr), nil)
	require.NoError(t, err, "Dial must not error")
	defer clientConn.Close()

	pongs := make(chan string, 1)
	clientConn.SetPongHandler(func(data string) error {
		pongs <- data
		return nil
	})

	// Client message is passed to Backend
	err = clientConn.WriteJSON("test client message")
	require.NoError(t, err, "WriteJSON must not error")

	next := <-h.in
	require.NoError(t, next.err, "NextMessage must not error")
	assert.Equal(t, websocket.TextMessage, next.typ, "Type should be correct")

	var s string
	err = json.Unmarshal(next.msg, &s)
	require.NoError(t, err, "Unmarshal must not error")
	assert.Equal(t, "test client message", s, "Message should be correct")

	// Send a ping and wait to make sure pong has time to roundtrip
	err = clientConn.WriteControl(websocket.PingMessage, []byte(`test control`), time.Now().Add(time.Second))
	require.NoError(t, err, "WriteJSON must not error")
	time.Sleep(500 * time.Millisecond)

	// Backend message is passed to client
	h.out <- "test server message"

	err = clientConn.ReadJSON(&s)
	require.NoError(t, err, "ReadJSON must not error")
	assert.Equal(t, "test server message", s, "Message should be correct")

	// Pong message has been passed
	pong := <-pongs
	assert.Equal(t, "test control", pong, "Pong message should be correct")
}

func TestCloseClientNormal(t *testing.T) {
	h := bckndHandler{
		in:  make(chan *nextMsg, 3),
		out: make(chan interface{}),
	}

	backSrv := httptest.NewServer(h)
	defer backSrv.Close()

	proxySrv, _ := createProxyServer(backSrv.URL)
	defer proxySrv.Close()

	proxyAddr := proxySrv.Listener.Addr().String()
	clientConn, _, err := dialer.Dial(fmt.Sprintf("ws://%v", proxyAddr), nil)
	require.NoError(t, err, "Dial must not error")
	defer clientConn.Close()

	err = clientConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "test close"), time.Now().Add(time.Second))
	require.NoError(t, err, "WriteControl must not error")

	next := <-h.in
	require.Error(t, next.err, "NextMessage must error")
	assert.Equal(t, &websocket.CloseError{Code: websocket.CloseNormalClosure, Text: "test close"}, next.err, "Error should be correct")
}

func TestCloseClient(t *testing.T) {
	h := bckndHandler{
		in:  make(chan *nextMsg, 3),
		out: make(chan interface{}),
	}

	backSrv := httptest.NewServer(h)
	defer backSrv.Close()

	proxySrv, _ := createProxyServer(backSrv.URL)
	defer proxySrv.Close()

	proxyAddr := proxySrv.Listener.Addr().String()
	clientConn, _, err := dialer.Dial(fmt.Sprintf("ws://%v", proxyAddr), nil)
	require.NoError(t, err, "Dial must not error")
	clientConn.Close()

	next := <-h.in
	require.Error(t, next.err, "NextMessage must error")
	assert.Equal(t, &websocket.CloseError{Code: websocket.CloseGoingAway}, next.err, "Error should be correct")
}

func TestProxyCloseServer(t *testing.T) {
	h := bckndHandler{
		in:   make(chan *nextMsg, 3),
		out:  make(chan interface{}),
		stop: make(chan struct{}),
	}

	backSrv := httptest.NewServer(h)
	defer backSrv.Close()

	proxySrv, _ := createProxyServer(backSrv.URL)
	defer proxySrv.Close()

	proxyAddr := proxySrv.Listener.Addr().String()
	clientConn, _, err := dialer.Dial(fmt.Sprintf("ws://%v", proxyAddr), nil)
	require.NoError(t, err, "Dial must not error")
	defer clientConn.Close()

	// Keeps client connection active
	go func() {
		for i := 0; i < 10; i++ {
			_ = clientConn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Simulate server closing server
	close(h.stop)

	var s string
	err = clientConn.ReadJSON(&s)
	require.Error(t, err, "ReadJSON must error")
	assert.Equal(t, &websocket.CloseError{Code: websocket.CloseGoingAway}, err, "Message should be correct")
}

func TestProxyStop(t *testing.T) {
	h := bckndHandler{
		in:   make(chan *nextMsg, 3),
		out:  make(chan interface{}),
		stop: make(chan struct{}),
	}

	backSrv := httptest.NewServer(h)
	defer backSrv.Close()

	proxySrv, prx := createProxyServer(backSrv.URL)

	proxyAddr := proxySrv.Listener.Addr().String()
	clientConn, _, err := dialer.Dial(fmt.Sprintf("ws://%v", proxyAddr), nil)
	require.NoError(t, err, "Dial must not error")
	defer clientConn.Close()

	// Keeps connection active
	go func() {
		for i := 0; i < 10; i++ {
			_ = clientConn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Close server
	proxySrv.Close()
	_ = proxySrv.Config.Shutdown(context.TODO())

	var s string
	err = clientConn.ReadJSON(&s)
	require.Error(t, err, "ReadJSON must error")
	assert.Equal(t, &websocket.CloseError{Code: websocket.CloseGoingAway}, err, "Message should be correct")

	next := <-h.in
	require.Error(t, next.err, "NextMessage must error")
	assert.Equal(t, &websocket.CloseError{Code: websocket.CloseGoingAway}, next.err, "Error should be correct")

	// Wait for proxy to complete
	<-prx.Done()
}
