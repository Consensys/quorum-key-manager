package request

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/header"
)

func HeadersPreparer(h func(http.Header) error) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		return req, h(req.Header)
	})
}

// Headers sets or deletes custom request headers
func Headers(overides map[string][]string) Preparer {
	return HeadersPreparer(header.Overide(overides))
}
