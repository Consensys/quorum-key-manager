package node

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/node"
)

type ctxKeyType string

const (
	ctxSessionKey ctxKeyType = "session"
)

func SessionFromContext(ctx context.Context) node.Session {
	n, ok := ctx.Value(ctxSessionKey).(node.Session)
	if !ok {
		return nil
	}

	return n
}

func WithSession(ctx context.Context, n node.Session) context.Context {
	return context.WithValue(ctx, ctxSessionKey, n)
}
