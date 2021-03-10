package api

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
)

func newHTTPMiddleware(bcknd core.Backend) func(http.Handler) http.Handler {
	// TODO: implement the sequence of middlewares to apply before routing
	return func(h http.Handler) http.Handler {
		return h
	}
}
