package request

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/header"
)

// Headers sets or deletes custom request headers
func Headers(overides map[string][]string) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		header.Overide(req.Header, overides)
		return req, nil
	})
}
