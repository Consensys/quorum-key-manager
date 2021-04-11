package jsonrpcapi

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	nodemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/node-manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
)

// Proxy proxies the request to session downstream JSON-RPC
func Proxy(rw jsonrpc.ResponseWriter, req *jsonrpc.Request) {
	node.SessionFromContext(req.Request().Context()).ProxyRPC().ServeRPC(rw, req)
}

// New creates a http.Handler to be served on JSON-RPC
func New(bcknd core.Backend) http.Handler {
	// Only JSON-RPC v2 is supported
	router := jsonrpc.NewRouter().DefaultHandler(jsonrpc.NotSupportedVersionHandler())
	v2Router := router.Version("2.0").Subrouter().DefaultHandler(jsonrpc.HandlerFunc(Proxy))

	// Silence personal
	v2Router.MethodPrefix("personal_").Handle(jsonrpc.MethodNotFoundHandler())

	// Wrap handler into middlewares
	handler := nodemanager.NewMiddleware(bcknd.NodeManager()).Next(router)
	handler = jsonrpc.LoggedHandler(handler)

	return jsonrpc.ToHTTPHandler(handler)
}
