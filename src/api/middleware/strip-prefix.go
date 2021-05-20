package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/handlers"
)

const StoreURLPlaceholder = "{.+}"

// Modified version of http.StripPrefix() to append a tail backslash in case prefix exact match with URL.path
func StripPrefix(prefix string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlPrefix := prefix
		if strings.Contains(prefix, StoreURLPlaceholder) {
			storeID := r.Context().Value(handlers.StoreContextID)
			urlPrefix = strings.Replace(urlPrefix, StoreURLPlaceholder, storeID.(string), 1)
		}

		if p := strings.TrimPrefix(r.URL.Path, urlPrefix); len(p) < len(r.URL.Path) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			if p == "" {
				p = "/"
			}

			r2.URL.Path = p
			h.ServeHTTP(w, r2)
		} else {
			http.NotFound(w, r)
		}
	})
}
