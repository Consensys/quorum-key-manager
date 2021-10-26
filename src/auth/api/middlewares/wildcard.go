package middlewares

import (
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"net/http"
)

func WildcardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r.Clone(WithUserInfo(r.Context(), entities.NewWildcardUser())))
	})
}
