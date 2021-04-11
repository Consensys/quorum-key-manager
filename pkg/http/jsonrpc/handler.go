package jsonrpc

import (
	"fmt"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

//go:generate mockgen -source=handler.go -destination=handler_mock.go -package=jsonrpc

// Handler is and JSON-RPC handler to be used in a JSON-RPC server
// It provides the JSON-RPC abstraction over http.Handler interface
type Handler interface {
	ServeRPC(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeRPC(rw ResponseWriter, req *Request) {
	f(rw, req)
}

// ToHTTPHandler wraps a jsonrpc.Handler into a http.Handler
func ToHTTPHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// extract JSON-RPC request from context
		rpcReq := RequestFromContext(req.Context())
		if rpcReq == nil {
			// if no JSON-RPC request is found then creates one and attached to http.Request context
			rpcReq = NewRequest(req)
			_ = rpcReq.ReadBody()
			rpcReq.req = req.WithContext(WithRequest(req.Context(), rpcReq))
		} else {
			// if found update http.Request
			rpcReq.req = req
		}

		rpcRw, ok := rw.(ResponseWriter)
		if !ok {
			rpcRw = NewResponseWriter(rw).WithVersion(rpcReq.Version()).WithID(rpcReq.ID())
		}

		// Serve
		h.ServeRPC(rpcRw, rpcReq)
	})
}

// FromHTTPHandler wraps a http.Handler into a jsonrpc.Handler
func FromHTTPHandler(h http.Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		// Write JSON-RPC request message into request body
		_ = req.WriteBody()

		// Serve HTTP request
		h.ServeHTTP(rw, req.Request())
	})
}

// NotSupported replies to the request with a not supported request error
func NotSupported(rw ResponseWriter, req *Request) {
	_ = rw.WriteError(&ErrorMsg{
		Message: "not supported",
	})
}

// NotSupportedHandler returns a simple handler
// that replies to each request with a not supported version request error
func NotSupportedHandler() Handler { return HandlerFunc(NotSupported) }

func NotSupportedVersion(rw ResponseWriter, req *Request) {
	_ = rw.WriteError(&ErrorMsg{
		Message: fmt.Sprintf("not supported version %v", req.Version()),
	})
}

// NotSupportedVersionHandler returns a simple handler
// that replies to each request with a not supported version request error
func NotSupportedVersionHandler() Handler { return HandlerFunc(NotSupportedVersion) }

// InvalidMethod replies to the request with an invalid method error
func InvalidMethod(rw ResponseWriter, req *Request) {
	_ = rw.WriteError(&ErrorMsg{
		Code:    -32601,
		Message: fmt.Sprintf("invalid method %q", req.Method()),
	})
}

// InvalidMethod returns a simple handler
// that replies to each request with an invalid method error
func InvalidMethodHandler() Handler { return HandlerFunc(InvalidMethod) }

// NotImplementedMethod replies to the request with an not implemented error
func NotImplementedMethod(rw ResponseWriter, req *Request) {
	_ = rw.WriteError(&ErrorMsg{
		Code:    -32601,
		Message: fmt.Sprintf("not implemented method %q", req.Method()),
	})
}

// NotImplementedMethodHandler returns a simple handler
// that replies to each request with an invalid method error
func NotImplementedMethodHandler() Handler { return HandlerFunc(NotImplementedMethod) }

// LoggedHandler
func LoggedHandler(h Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		logger := log.FromContext(req.Request().Context())
		logger.
			WithField("version", req.Version()).
			WithField("id", req.ID()).
			WithField("method", req.Method()).
			Info("serve JSON-RPC request")
		h.ServeRPC(rw, req)
	})
}
