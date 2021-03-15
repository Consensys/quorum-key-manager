package jsonrpc

import (
	"fmt"
	"net/http"
)

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
			// if no JSON-RPC request is found then creates on and attached to http.Request context
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
		_ = req.WriteBody()
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
// that replies to each request with a not supported request error
func NotSupportedHandler() Handler { return HandlerFunc(NotSupported) }

// InvalidMethod replies to the request with an invalid method error
func InvalidMethod(rw ResponseWriter, req *Request) {
	_ = rw.WriteError(&ErrorMsg{
		Code:    -32601,
		Message: fmt.Sprintf("invalid method %q", req.Method()),
	})
}

// InvalidMethod returns a simple handler
// that replies to each request with an invalid method error
func MethodNotFoundHandler() Handler { return HandlerFunc(InvalidMethod) }
