package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		http.Error(w, reason.Error(), status)
	},
}

var dialer = websocket.Dialer{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 30 * time.Second,
}

type handler struct {
	*testing.T
	reqs  chan *RequestMsg
	resps chan *ResponseMsg
}

func newServer(t *testing.T) (*httptest.Server, handler) {
	h := handler{
		T:     t,
		reqs:  make(chan *RequestMsg),
		resps: make(chan *ResponseMsg),
	}
	return httptest.NewServer(h), h
}

func (t handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("test-server: Upgrade: %v\n", err)
		return
	}

	go func() {
		defer conn.Close()
		for {
			reqMsg := new(RequestMsg)
			e := conn.ReadJSON(reqMsg)
			fmt.Printf("test-server: ReadJSON: resp=%v err=%v\n", *reqMsg, e)
			if e != nil {
				return
			}
		}
	}()

	go func() {
		defer conn.Close()
		for resp := range t.resps {
			err = conn.WriteJSON(resp)
			fmt.Printf("test-server: WriteJSON: %v\n", err)
			if err != nil {
				return
			}
		}
	}()
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}

type TestCase struct {
	req *RequestMsg

	err  error
	resp *ResponseMsg

	expectedErr  error
	expectedResp *ResponseMsg

	done chan struct{}
}

func (cas *TestCase) isDone() bool {
	after := time.After(50 * time.Millisecond)
	select {
	case <-cas.done:
		return true
	case <-after:
		return false
	}
}

func (cas *TestCase) assert(t *testing.T) {
	if cas.expectedErr != nil {
		require.Error(t, cas.err, "Case must error")
		assert.Equal(t, cas.expectedErr, cas.err, "Error must be correct")
	} else {
		require.NoError(t, cas.err, "Case must error")
	}

	if cas.expectedResp == nil {
		require.Nil(t, cas.resp, "Response must be nil")
	} else {
		expectedResp, _ := json.Marshal(cas.expectedResp)
		resp, _ := json.Marshal(cas.resp)
		assert.Equal(t, expectedResp, resp, "Resp must match")
	}
}

func TestWebSocketClient(t *testing.T) {
	s, h := newServer(t)
	defer s.Close()

	url := makeWsProto(s.URL)
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		require.NoError(t, err, "Dial must not error")
	}
	defer conn.Close()

	client := NewWebsocketClient(conn)
	err = client.Start(context.TODO())
	require.NoError(t, err, "Start should not error")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	cases := []*TestCase{
		&TestCase{
			req: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod0",
				Params:  []int{1, 2, 3},
				ID:      1,
			},
			expectedResp: &ResponseMsg{
				Version: "2.0",
				Result:  "abcd",
				ID:      1,
			},
			done: make(chan struct{}),
		},
		&TestCase{
			req: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod1",
				Params:  []int{1, 2, 3},
				ID:      2,
			},
			expectedResp: &ResponseMsg{
				Version: "2.0",
				Result:  "abcd",
				ID:      2,
			},
			done: make(chan struct{}),
		},
		&TestCase{
			req: &RequestMsg{
				Version: "2.0",
				Method:  "testMethod2",
				Params:  []int{1, 2, 3},
				ID:      "abcd",
			},
			expectedResp: &ResponseMsg{
				Version: "2.0",
				Result:  "abcd",
				ID:      "abcd",
			},
			done: make(chan struct{}),
		},
		&TestCase{
			req: (&RequestMsg{
				Version: "2.0",
				Method:  "testMethod3",
				Params:  []int{1, 2, 3},
				ID:      "ef",
			}).WithContext(timeoutCtx),
			expectedErr: &ErrorMsg{
				Code:    -32000,
				Message: "Downstream error",
				Data: map[string]interface{}{
					"message": "Client Closed Connection",
					"status":  499,
				},
			},
			done: make(chan struct{}),
		},
	}

	for _, cas := range cases {
		go func(cas *TestCase) {
			defer func() { close(cas.done) }()
			cas.resp, cas.err = client.Do(cas.req)
		}(cas)
	}

	time.Sleep(100 * time.Millisecond)

	// Return response in distinct order
	h.resps <- cases[2].expectedResp
	h.resps <- cases[0].expectedResp

	require.True(t, cases[0].isDone(), "Case 0 must have completed")
	cases[0].assert(t)

	require.True(t, cases[2].isDone(), "Case 2 must have completed")
	cases[0].assert(t)

	require.False(t, cases[1].isDone(), "Case 1 must not have completed")

	// Sleep to trigger deadline on case 3
	time.Sleep(200 * time.Millisecond)
	require.True(t, cases[3].isDone(), "Case 3 must have completed")
	cases[3].assert(t)

	h.resps <- cases[1].expectedResp

	require.True(t, cases[1].isDone(), "Case 1 must have completed")
	cases[1].assert(t)

	err = client.Stop(context.TODO())
	require.NoError(t, err, "Stop must not error")

	close(h.resps)
}
