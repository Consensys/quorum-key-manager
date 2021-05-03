package jsonrpcapi

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
)

// New creates a http.Handler to be served on JSON-RPC
func New(bcknd core.Backend) http.Handler {
	return http.HandleFunc(func(rw http.ResponseWriter, msg *http.Request){
		n, err := bcknd.NodeManager().Node(req.Context(), "default")
		if err != nil {
			http.NotFound(rw, req)
			return
		}
		n.ServeHTTP(rw, err)
	})
}
