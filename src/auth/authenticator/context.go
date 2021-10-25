package authenticator

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

type ctxKey string

var reqCtxKey ctxKey = "req"

func UserContextFromContext(ctx context.Context) *UserContext {
	if reqCtx, ok := ctx.Value(reqCtxKey).(*UserContext); ok {
		return reqCtx
	}
	return nil
}

func UserInfoContextFromContext(ctx context.Context) *entities.UserInfo {
	userCtx := UserContextFromContext(ctx)
	if userCtx == nil {
		return nil
	}
	return userCtx.UserInfo
}

func WithUserContext(ctx context.Context, reqCtx *UserContext) context.Context {
	return context.WithValue(ctx, reqCtxKey, reqCtx)
}

// UserContext is a set of data attached to every incoming request
type UserContext struct {
	// UserInfo records user information
	UserInfo *entities.UserInfo
}

func NewUserContext(userInfo *entities.UserInfo) *UserContext {
	return &UserContext{userInfo}
}
