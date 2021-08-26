package request

import (
	"net/http"
	"net/url"
	"path"
)

// ExtractURI parses request URI, updates request URL and optional reset URI
func ExtractURI(reset bool) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if req.RequestURI == "" {
			return req, nil
		}

		u, err := url.ParseRequestURI(req.RequestURI)
		if err != nil {
			return req, err
		}

		if req.URL.Path != u.Path {
			req.URL.RawPath = path.Join(fistNotEmpty(req.URL.RawPath, req.URL.Path), fistNotEmpty(u.RawPath, u.Path))
			req.URL.Path = path.Join(req.URL.Path, u.Path)

		}

		req.URL.RawQuery = u.RawQuery
		if reset {
			req.RequestURI = ""
			req.Host = ""
		}

		return req, nil
	})
}

func fistNotEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return ""
}
