package request

import (
	"net/http"
	"net/url"
)

// ExtractURI parses request URI, updates request URL and optionnaly reset URI
func ExtractURI(reset bool) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if req.RequestURI != "" {
			u, err := url.ParseRequestURI(req.RequestURI)
			if err != nil {
				return req, err
			}

			req.URL.Path = u.Path
			req.URL.RawPath = u.RawPath
			req.URL.RawQuery = u.RawQuery
			if reset {
				req.RequestURI = ""
			}
		}

		return req, nil
	})
}
