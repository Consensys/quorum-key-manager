package middleware

import (
	"net/http"
	"strings"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/middleware"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/handlers"
)

const StoreURLPlaceholder = "{.+}"

// Modified version of http.StripPrefix() to append a tail backslash in case prefix exact match with URL.path
func StripPrefix(prefix string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			urlPrefix := prefix
			if strings.Contains(prefix, StoreURLPlaceholder) {
				storeID := r.Context().Value(handlers.StoreContextID)
				urlPrefix = strings.Replace(urlPrefix, StoreURLPlaceholder, storeID.(string), 1)
			}

			middleware.StripPrefix(urlPrefix)(h)
		})
	}
}
