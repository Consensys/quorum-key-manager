package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync/atomic"
)

var version = "2.0"

// Client is a connector to a jsonrpc server
type Client struct {
	transport http.RoundTripper

	version string

	req *http.Request

	idCounter uint32
}

// WithTransport copies client ant attaches transport
func (c *Client) WithTransport(t http.RoundTripper) *Client {
	cpy := &Client{
		transport: t,
		version:   c.version,
		idCounter: c.idCounter,
	}

	if c.req != nil {
		cpy.req = c.req.Clone(c.req.Context())
	}

	return cpy
}

// WithRequest copies client and attaches request to it
func (c *Client) WithRequest(req *http.Request) *Client {
	cpy := &Client{
		transport: c.transport,
		version:   c.version,
		req:       req,
	}

	reqMsg := RequestMsgFromContext(req.Context())
	if reqMsg.Version != "" {
		cpy.version = reqMsg.Version
		_ = json.Unmarshal(reqMsg.ID, &cpy.idCounter)
	}

	return cpy
}

// Version returns jsonrpc version
func (client *Client) Version() string {
	if client.version != "" {
		return client.version
	}

	return version
}

// Call sends a JSON-RPC request over underlying http.Transport

// Returns an http.Response which body as already been consumed in the jsonrpc.ResponseMsg

// It returns an error in following scenarios
// - underlyng transport failed to roundtrip
// - response status code is not 2XX
// - response body is an invalid JSON-RPC response
// - JSON-RPC response is failed (in which case it returns the jsonrpc.Error)
func (client *Client) Call(ctx context.Context, method string, params interface{}) (*Response, error) {
	// Create and prepare request
	reqMsg := client.newRequestMsg(method)
	if err := reqMsg.WithParams(params); err != nil {
		return nil, err
	}
	req := client.newRequest(ctx)
	prepareRequestBody(req, reqMsg)

	// Execute HTTP call
	resp, err := client.do(req)
	if err != nil {
		return nil, err
	}

	// Check for HTTP error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("%v (code=%v)", resp.Status, resp.StatusCode)
	}

	defer resp.Body.Close()

	// Read response body
	resp.Msg = new(ResponseMsg)
	if err := newClientCodec(nil, resp.Body).ReadResponse(resp.Msg); err != nil {
		return resp, err
	}

	// Validate response JSON-RPC message
	if err := resp.Msg.Validate(); err != nil {
		return resp, err
	}

	return resp, resp.Msg.Error
}

func (client *Client) do(req *http.Request) (*Response, error) {
	resp, err := client.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return &Response{Response: resp}, nil
}

func (client *Client) newRequestMsg(method string) *RequestMsg {
	return &RequestMsg{
		Version: client.Version(),
		Method:  method,
		ID:      client.nextID(),
	}
}

func (client *Client) newRequest(ctx context.Context) *http.Request {
	return client.req.Clone(ctx)
}

func prepareRequestBody(req *http.Request, msg *RequestMsg) *http.Request {
	// Write request msg to buffer
	buf := new(bytes.Buffer)
	_ = newClientCodec(buf, nil).WriteRequest(msg)

	// Set request body with buffer
	req.Body = ioutil.NopCloser(buf)
	req.ContentLength = int64(buf.Len())

	return req
}

func (client *Client) nextID() json.RawMessage {
	id := atomic.AddUint32(&client.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}
