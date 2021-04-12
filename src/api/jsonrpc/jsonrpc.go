package jsonrpcapi

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
)

// New creates a http.Handler to be served on JSON-RPC
func New(bcknd core.Backend) http.Handler {
	// Only JSON-RPC v2 is supported
	router := jsonrpc.NewRouter().DefaultHandler(jsonrpc.NotSupportedVersionHandler())
	v2Router := router.Version("2.0").Subrouter().DefaultHandler(jsonrpc.HandlerFunc(Proxy))

	// Silence JSON-RPC personal
	v2Router.MethodPrefix("personal_").Handle(jsonrpc.MethodNotFoundHandler())

	// Wrap router into middlewares

	// Node Manager Middleware is responsible to attach Node session to context
	handler := NewNodeMiddleware(bcknd.NodeManager()).Next(router)

	handler = jsonrpc.LoggedHandler(handler)

	return jsonrpc.ToHTTPHandler(handler)
}
