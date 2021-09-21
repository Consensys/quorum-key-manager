package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrClientQuit = errors.New("client is closed")
	ErrClientStop = errors.New("client is stopped")
)

// WebSocketClient is a connector to a jsonrpc server
type WebSocketClient struct {
	conn *websocket.Conn

	writeTimeout time.Duration

	liveOps map[string]*operation

	todos     chan *operation
	failedOps chan *operation

	readResp chan *ResponseMsg
	readErr  chan error

	stop    chan struct{}
	closing chan struct{}
	close   chan struct{}

	err error

	errors chan error
}

// NewWebsocketClient creates a new jsonrpc HTTPClient from an HTTP HTTPClient
func NewWebsocketClient(conn *websocket.Conn) *WebSocketClient {
	return &WebSocketClient{
		conn:         conn,
		writeTimeout: 10 * time.Second,
		liveOps:      make(map[string]*operation),
		todos:        make(chan *operation),
		failedOps:    make(chan *operation),
		readResp:     make(chan *ResponseMsg, 20),
		readErr:      make(chan error),
		closing:      make(chan struct{}),
		close:        make(chan struct{}),
		stop:         make(chan struct{}),
		errors:       make(chan error, 1),
	}
}

func (c *WebSocketClient) Start(context.Context) error {
	go c.read()
	go c.manageOp()
	return nil
}

// Stop the client from receiving new messages
// It finishes reading all messages before quiting
func (c *WebSocketClient) Stop(ctx context.Context) error {
	// Close stop channel
	close(c.stop)

	// Set connection readline to Now so NextReader return instantaneously
	_ = c.conn.SetReadDeadline(time.Now())

	// Blocks until close signal
	select {
	case <-c.close:
		return c.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *WebSocketClient) Errors() <-chan error {
	return c.errors
}

// Do sends a jsonrpc request over the underlying HTTP client and returns a jsonrpc response
func (c *WebSocketClient) Do(reqMsg *RequestMsg) (*ResponseMsg, error) {
	err := reqMsg.Validate()
	if err != nil {
		return nil, err
	}

	rawID, _ := json.Marshal(reqMsg.ID)
	op := &operation{
		c:    c,
		id:   string(rawID),
		msg:  reqMsg,
		sent: make(chan error),
		resp: make(chan *ResponseMsg),
	}

	err = c.send(op)
	if err != nil {
		return nil, DownstreamError(err)
	}

	// dispatch has accepted the request and will close the channel when it quits.
	respMsg, err := op.wait()
	if err != nil {
		return nil, DownstreamError(err)
	}

	return respMsg, nil
}

func (c *WebSocketClient) send(op *operation) error {
	ctx := op.msg.Context()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.closing:
		return ErrClientQuit
	case c.todos <- op:
		err := <-op.sent
		return err
	}
}

func (c *WebSocketClient) read() {
	for {
		_, r, err := c.conn.NextReader()
		if err != nil {
			select {
			case <-c.stop:
				c.readErr <- ErrClientStop
			default:
				c.readErr <- err
			}
			return
		}

		respMsg := new(ResponseMsg)
		err = json.NewDecoder(r).Decode(respMsg)
		if err != nil {
			continue
		}

		err = respMsg.Validate()
		if err != nil {
			continue
		}

		c.readResp <- respMsg
	}
}

func (c *WebSocketClient) manageOp() {
	for {
		select {
		case op := <-c.todos:
			// Register op
			c.addOp(op)

			// Compute write deadline
			ctx := op.msg.Context()
			deadline, ok := ctx.Deadline()
			if !ok {
				deadline = time.Now().Add(c.writeTimeout)
			}

			// Write Message
			_ = op.c.conn.SetWriteDeadline(deadline)
			err := op.c.conn.WriteJSON(op.msg)
			op.sent <- err

			if err != nil {
				select {
				case op.c.failedOps <- op:
				case <-op.c.closing:
				}
			}
		case msg := <-c.readResp:
			c.handleRespMsg(msg)
		case err := <-c.readErr:
			// Indicate that we start closing
			close(c.closing)

			// Stop processing msgs
			close(c.readErr)
			close(c.readResp)

			// Finish processing all responses that were already received
			for msg := range c.readResp {
				c.handleRespMsg(msg)
			}

			// Cancel remaining operations
			c.cancelAllOps(err)

			close(c.failedOps)
			if err != ErrClientStop {
				c.err = err
				c.errors <- err
			}
			close(c.close)
			close(c.errors)

			return
		case op := <-c.failedOps:
			c.removeOp(op)
		}
	}
}

func (c *WebSocketClient) addOp(op *operation) {
	c.liveOps[op.id] = op
}

func (c *WebSocketClient) removeOp(op *operation) {
	delete(c.liveOps, op.id)
}

// cancelAllOps cancel all current ops
func (c *WebSocketClient) cancelAllOps(err error) {
	for _, op := range c.liveOps {
		c.removeOp(op)
		op.err = err
		close(op.resp)
	}
}

func (c *WebSocketClient) handleRespMsg(respMsg *ResponseMsg) {
	rawID, _ := json.Marshal(respMsg.ID)
	op, ok := c.liveOps[string(rawID)]
	if !ok {
		return
	}

	op.resp <- respMsg
	c.removeOp(op)
}

type operation struct {
	c *WebSocketClient

	id  string
	msg *RequestMsg

	sent chan error

	resp chan *ResponseMsg
	err  error
}

func (op *operation) wait() (*ResponseMsg, error) {
	ctx := op.msg.Context()

	select {
	case <-ctx.Done():
		select {
		case op.c.failedOps <- op:
		case <-op.c.closing:
		}

		return nil, ctx.Err()
	case resp := <-op.resp:
		return resp, op.err
	}
}
