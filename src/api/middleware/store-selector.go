package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/handlers"
)

func StoreSelector(storePath string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, storePath) {
			h.ServeHTTP(w, r)
			return
		}

		pieces := strings.Split(r.URL.Path[1:], "/")
		ctx := context.WithValue(r.Context(), handlers.StoreContextID, pieces[1])
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
