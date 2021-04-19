package jsonrpc

import (
	"context"
)

// Caller is an interface for a JSON-RPC caller
type Caller interface {
	Call(ctx context.Context, method string, params interface{}) (*Response, error)
}

type caller struct {
	client Client
	req    *Request
}

func NewCaller(c Client, req *Request) Caller {
	return &caller{
		client: c,
		req:    req,
	}
}

// Call sends a JSON-RPC request over underlying http.Transport

// Returns an http.Response which body as already been consumed in the jsonrpc.ResponseMsg

// It returns an error in following scenarios
// - underlying transport failed to roundtrip
// - response status code is not 2XX
// - response body is an invalid JSON-RPC response
// - JSON-RPC response is failed (in which case it returns the jsonrpc.ErrorMsg)
func (c *caller) Call(ctx context.Context, method string, params interface{}) (*Response, error) {
	req := RequestFromContext(ctx)
	if req == nil {
		req = c.req
	}

	return c.client.Do(req.Clone(ctx).WithMethod(method).WithParams(params))
}
