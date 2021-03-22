package jsonrpc

import (
	"context"
	"net/http"
	"sync/atomic"

	httpclient "github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/client"
)

var defaultVersion = "2.0"

type ClientConfig struct {
	Version string             `json:"version,omitempty"`
	HTTP    *httpclient.Config `json:"http,omitempty"`
}

func (cfg *ClientConfig) SetDefault() {
	if cfg.HTTP == nil {
		cfg.HTTP = new(httpclient.Config)
	}

	cfg.HTTP.SetDefault()

	if cfg.Version == "" {
		cfg.Version = defaultVersion
	}
}

// Client is a connector to a jsonrpc server
type Client struct {
	cfg *ClientConfig

	client httpclient.Client
}

func NewClient(cfg *ClientConfig, client httpclient.Client) (*Client, error) {
	if cfg == nil {
		cfg = new(ClientConfig)
	}

	cfg.SetDefault()

	if client == nil {
		var err error
		client, err = httpclient.New(cfg.HTTP, nil)
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		cfg:    cfg,
		client: client,
	}, nil
}

// Version returns jsonrpc version
func (c *Client) Version() string {
	return c.cfg.Version
}

func (c *Client) Do(req *Request) (*Response, error) {
	// write request body
	err := req.WriteBody()
	if err != nil {
		return nil, err
	}

	httpResp, err := c.client.Do(req.Request())
	if err != nil {
		return nil, err
	}

	// Create response and reads body
	resp := NewResponse(httpResp)
	err = resp.ReadBody()
	if err != nil {
		return resp, err
	}

	err = resp.Error()
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Client) Caller(req *http.Request) Caller {
	return &caller{
		client:  c,
		req:     req,
		version: c.cfg.Version,
	}
}

type Caller interface {
	Call(ctx context.Context, method string, params interface{}) (*Response, error)
}

type caller struct {
	client *Client

	req *http.Request

	version   string
	idCounter uint32
}

// Call sends a JSON-RPC request over underlying http.Transport

// Returns an http.Response which body as already been consumed in the jsonrpc.ResponseMsg

// It returns an error in following scenarios
// - underlying transport failed to roundtrip
// - response status code is not 2XX
// - response body is an invalid JSON-RPC response
// - JSON-RPC response is failed (in which case it returns the jsonrpc.ErrorMsg)
func (c *caller) Call(ctx context.Context, method string, params interface{}) (*Response, error) {
	req := c.newRequest(ctx).WithMethod(method).WithParams(params)

	return c.client.Do(req)
}

func (c *caller) newRequest(ctx context.Context) *Request {
	return NewRequest(c.req.Clone(ctx)).WithID(c.nextID()).WithVersion(c.version)
}

func (c *caller) nextID() int {
	return int(atomic.AddUint32(&c.idCounter, 1))
}
