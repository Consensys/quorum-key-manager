package jsonrpcapi

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
)

// Proxy proxies the request to session downstream JSON-RPC
func Proxy(rw jsonrpc.ResponseWriter, req *jsonrpc.Request) {
	// Extract node session from context
	sess := node.SessionFromContext(req.Request().Context())

	// Proxy request
	sess.ProxyRPC().ServeRPC(rw, req)
}
