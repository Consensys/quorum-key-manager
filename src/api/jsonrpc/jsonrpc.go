package jsonrpcapi

import (
	"encoding/json"
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

	// TODO: when implementing interceptors below is an example of declaring a JSON-RPC route
	// v2Router.MethodPrefix("eth_sendTransaction").Handle(jsonrpc.SendTransationHander())

	// Wrap router into middlewares

	handler := jsonrpc.LoggedHandler(router)

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Create ResponseWriter
		rpcRw := jsonrpc.NewResponseWriter(rw)

		// Parse request body
		msg := new(jsonrpc.RequestMsg)
		err := json.NewDecoder(req.Body).Decode(msg)
		req.Body.Close()
		if err != nil {
			_ = rpcRw.WriteError(jsonrpc.ParseError(err))
			return
		}

		// Serve
		handler.ServeRPC(rpcRw, msg.WithContext(req.Context()))
	})
}
