package types

import "context"

type ctxKey string

var reqCtxKey ctxKey = "req"

func UserContextFromContext(ctx context.Context) *UserContext {
	if reqCtx, ok := ctx.Value(reqCtxKey).(*UserContext); ok {
		return reqCtx
	}
	return nil
}

func WithUserContext(ctx context.Context, reqCtx *UserContext) context.Context {
	return context.WithValue(ctx, reqCtxKey, reqCtx)
}
