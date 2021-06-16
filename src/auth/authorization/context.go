package authorization

import "context"

type ctxKey string

var resCtxKey ctxKey = "resolver"

func WithResolver(ctx context.Context, res *Resolver) context.Context {
	return context.WithValue(ctx, resCtxKey, res)
}

func ResolverFromContext(ctx context.Context) *Resolver {
	if resolver, ok := ctx.Value(resCtxKey).(*Resolver); ok {
		return resolver
	}
	return nil
}
