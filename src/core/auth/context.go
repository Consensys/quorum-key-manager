package auth

import "context"

type ctxKey string

var authKey ctxKey = "auth"

// WithAuth attaches an Auth to context
func WithAuth(ctx context.Context, auth *Auth) context.Context {
	return context.WithValue(ctx, authKey, auth)
}

// FromContext looks for an Auth attached to context
// If no Auth is attached, it returns a NonAuthenticatedAuth
func FromContext(ctx context.Context) *Auth {
	auth, ok := ctx.Value(authKey).(*Auth)
	if !ok {
		return NotAuthenticatedAuth
	}
	return auth
}
