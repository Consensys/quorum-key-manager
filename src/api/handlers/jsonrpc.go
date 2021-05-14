package handlers

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
)

// New creates a http.Handler to be served on JSON-RPC
func NewJsonRPCHandler(bcknd core.Backend) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		n, err := bcknd.NodeManager().Node(req.Context(), "default")
		if err != nil {
			http.NotFound(rw, req)
			return
		}
		n.ServeHTTP(rw, req)
	})
}
