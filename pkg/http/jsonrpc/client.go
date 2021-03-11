package jsonrpc

import (
	"context"
	"net/http"
	"sync/atomic"
)

var defaultVersion = "2.0"

// Client is a connector to a jsonrpc server
type Client struct {
	version string

	transport http.RoundTripper

	req *http.Request

	idCounter uint32
}

// NewClient creates a client
func NewClient(addr, version string) (*Client, error) {
	req, err := http.NewRequest(http.MethodPost, addr, nil)
	if err != nil {
		return nil, err
	}

	return NewClientFromRequest(req, version), nil
}

func NewClientFromRequest(req *http.Request, version string) *Client {
	if version == "" {
		version = defaultVersion
	}

	return &Client{
		version:   version,
		transport: http.DefaultTransport,
		req:       req,
	}
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

	return cpy
}

// Version returns jsonrpc version
func (c *Client) Version() string {
	return c.version
}

// Call sends a JSON-RPC request over underlying http.Transport

// Returns an http.Response which body as already been consumed in the jsonrpc.ResponseMsg

// It returns an error in following scenarios
// - underlying transport failed to roundtrip
// - response status code is not 2XX
// - response body is an invalid JSON-RPC response
// - JSON-RPC response is failed (in which case it returns the jsonrpc.ErrorMsg)
func (c *Client) Call(ctx context.Context, method string, params interface{}) (*Response, error) {
	req := c.newRequest(ctx).
		WithVersion(c.Version()).
		WithMethod(method).
		WithParams(params).
		WithID(c.nextID())

	return c.do(req)
}

func (client *Client) newRequest(ctx context.Context) *Request {
	return NewRequest(client.req.Clone(ctx))
}

func (c *Client) do(req *Request) (*Response, error) {
	// write request body
	err := req.WriteBody()
	if err != nil {
		return nil, err
	}

	httpResp, err := c.transport.RoundTrip(req.Request())
	if err != nil {
		return nil, err
	}

	// Create response and reads body
	resp := NewResponse(httpResp)
	err = resp.ReadBody()
	if err != nil {
		return resp, err
	}

	return resp, resp.Error()
}

func (c *Client) nextID() int {
	return int(atomic.AddUint32(&c.idCounter, 1))
}
