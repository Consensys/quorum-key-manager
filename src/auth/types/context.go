package types

import "context"

type ctxKey string

var reqCtxKey ctxKey = "req"

func RequestContextFromContext(ctx context.Context) *RequestContext {
	if reqCtx, ok := ctx.Value(reqCtxKey).(*RequestContext); ok {
		return reqCtx
	}
	return nil
}

func WithRequestContext(ctx context.Context, reqCtx *RequestContext) context.Context {
	return context.WithValue(ctx, reqCtxKey, reqCtx)
}
