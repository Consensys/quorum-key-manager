package http_middlewares

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

func WildcardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r.Clone(WithUserInfo(r.Context(), entities.NewWildcardUser())))
	})
}
