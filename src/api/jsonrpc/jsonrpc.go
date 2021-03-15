package jsonrpcapi

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
)

// New creates a http.Handler to be served on JSON-RPC
func New(bcknd core.Backend) http.Handler {
	// Only JSON-RPC v2 is supported
	router := jsonrpc.NewRouter().DefaultHandler(jsonrpc.NotSupportedVersionHandler())
	v2Router := router.Version("2.0").Subrouter().DefaultHandler(jsonrpc.NotImplementedHandler())

	// Silence personal
	v2Router.MethodPrefix("personal_").Handle(jsonrpc.InvalidMethodHandler())

	// Wrap handler into middlewares
	handler := jsonrpc.LoggedHandler(router)

	return jsonrpc.ToHTTPHandler(handler)
}
