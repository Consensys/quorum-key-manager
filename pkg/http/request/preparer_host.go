package request

import (
	"net/http"
)

// Host set request Host
// This is useful when proxying request (to set the request host to match downstream server)

// If host is nil then it sets the Host to request URL host
func Host(host *string) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if host == nil {
			req.Host = req.URL.Host
		} else {
			req.Host = *host
		}
		return req, nil
	})
}
