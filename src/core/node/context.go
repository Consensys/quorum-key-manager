package node

import (
	"context"
)

type ctxKeyType string

const (
	ctxNodeKey ctxKeyType = "node"
)

func FromContext(ctx context.Context) Node {
	n, ok := ctx.Value(ctxNodeKey).(Node)
	if !ok {
		return nil
	}

	return n
}

func WithNode(ctx context.Context, n Node) context.Context {
	return context.WithValue(ctx, ctxNodeKey, n)
}
