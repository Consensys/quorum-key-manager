package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

// StripPrefix is a modified version of http.StripPrefix() to append a tail backslash in case prefix exact match with URL.path
func StripPrefix(prefix string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
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
}
