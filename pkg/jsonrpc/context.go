package jsonrpc

import (
	"context"
)

type ctxKey string

var (
	reqCtxKey ctxKey = "req"
)

// WithRequest attaches a Request to context
func WithRequest(ctx context.Context, req *Request) context.Context {
	return context.WithValue(ctx, reqCtxKey, req)
}

// RequestMsgFromContext looks for a RequestMsg attached to context
func RequestFromContext(ctx context.Context) *Request {
	req, ok := ctx.Value(reqCtxKey).(*Request)
	if !ok {
		return nil
	}
	return req
}
