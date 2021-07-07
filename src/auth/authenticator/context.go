package authenticator

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

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

// UserContext is a set of data attached to every incoming request
type UserContext struct {
	// UserInfo records user information
	UserInfo *types.UserInfo
}

func NewUserContext(userInfo *types.UserInfo) *UserContext {
	return &UserContext{userInfo}
}
