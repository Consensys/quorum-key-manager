package middlewares

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

type contextKey struct{}

func UserInfoFromContext(ctx context.Context) *entities.UserInfo {
	if reqCtx, ok := ctx.Value(contextKey{}).(*entities.UserInfo); ok {
		return reqCtx
	}
	return nil
}

func WithUserInfo(ctx context.Context, reqCtx *entities.UserInfo) context.Context {
	return context.WithValue(ctx, contextKey{}, reqCtx)
}
