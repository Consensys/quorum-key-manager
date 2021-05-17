package proxynode

import (
	"context"
)

type ctxKeyType string

const (
	ctxSessionKey ctxKeyType = "session"
)

func SessionFromContext(ctx context.Context) Session {
	n, ok := ctx.Value(ctxSessionKey).(Session)
	if !ok {
		return nil
	}

	return n
}

func WithSession(ctx context.Context, n Session) context.Context {
	return context.WithValue(ctx, ctxSessionKey, n)
}
