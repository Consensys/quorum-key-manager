package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
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
			err = clientConn.WriteJSON(out)
			fmt.Printf("test-server: WriteJSON: %v\n", err)
			if err != nil {
				return
			}
		}
	}()
}

func createProxyServer(uri string) *httptest.Server {
	prep, _ := request.Proxy(&request.ProxyConfig{Addr: uri})
	modif := response.Proxy(&response.ProxyConfig{})
	prx := &Proxy{
		Upgrader:               upgrader,
		Dialer:                 dialer,
		Interceptor:            Forward,
		PingPongTimeout:        time.Second,
		WriteControlMsgTimeout: time.Second,
		ReqPreparer:            prep,
		RespModifier:           modif,
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		prx.ServeHTTP(w, req)
	}))
}

func TestProxy(t *testing.T) {
	h := bckndHandler{
		in:  make(chan *nextMsg, 3),
		out: make(chan interface{}),
	}

	backSrv := httptest.NewServer(h)
	defer backSrv.Close()

	proxySrv := createProxyServer(backSrv.URL)
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

func TestProxyCloseClientNormal(t *testing.T) {
	h := bckndHandler{
		in:  make(chan *nextMsg, 3),
		out: make(chan interface{}),
	}

	backSrv := httptest.NewServer(h)
	defer backSrv.Close()

	proxySrv := createProxyServer(backSrv.URL)
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

func TestProxyCloseClient(t *testing.T) {
	h := bckndHandler{
		in:  make(chan *nextMsg, 3),
		out: make(chan interface{}),
	}

	backSrv := httptest.NewServer(h)
	defer backSrv.Close()

	proxySrv := createProxyServer(backSrv.URL)
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
		in:  make(chan *nextMsg, 3),
		out: make(chan interface{}),
	}

	backSrv := httptest.NewServer(h)

	proxySrv := createProxyServer(backSrv.URL)
	defer proxySrv.Close()

	proxyAddr := proxySrv.Listener.Addr().String()
	clientConn, _, err := dialer.Dial(fmt.Sprintf("ws://%v", proxyAddr), nil)
	require.NoError(t, err, "Dial must not error")
	defer clientConn.Close()
	// Make sure client conn will not timeout
	_ = clientConn.SetReadDeadline(time.Now().Add(10 * time.Second))

	// Close server and waits for server to close
	backSrv.Close()

	// Keeps client connection active
	go func() {
		for i := 0; i < 3; i++ {
			_ = clientConn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Waits for server to close
	time.Sleep(200 * time.Millisecond)

	var s string
	err = clientConn.ReadJSON(&s)
	require.Error(t, err, "ReadJSON must error")
	assert.Equal(t, &websocket.CloseError{Code: websocket.CloseGoingAway}, err, "Message should be correct")
}
