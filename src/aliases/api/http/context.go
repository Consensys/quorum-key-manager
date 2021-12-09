package http

import "context"

type contextKey struct{}

func WithRegistryName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, contextKey{}, name)
}

func RegistryNameFromContext(ctx context.Context) string {
	name, ok := ctx.Value(contextKey{}).(string)
	if ok {
		return name
	}

	return ""
}
