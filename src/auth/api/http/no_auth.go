package http

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

type NoAuth struct {
}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (m *NoAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r.WithContext(WithUserInfo(r.Context(), entities.NewWildcardUser())))
	})
}
