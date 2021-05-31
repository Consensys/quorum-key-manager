package handlers

import "context"

type ctxKeyType string

const storeNameCtxKey ctxKeyType = "storeName"

func WithStoreName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, storeNameCtxKey, name)
}

func StoreNameFromContext(ctx context.Context) string {
	name, ok := ctx.Value(storeNameCtxKey).(string)
	if ok {
		return name
	}

	return ""
}
