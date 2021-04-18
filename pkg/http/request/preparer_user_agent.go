package request

import (
	"net/http"
)

// UserAgent sets User-Agent header
func UserAgent(agent string) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if agent != "" {
			req.Header.Set("User-Agent", agent)
		}

		return req, nil
	})
}
