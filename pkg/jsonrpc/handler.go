package jsonrpc

import (
	"fmt"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
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

// DefaultRWHandler is an utility middleware that attaches request ID and Version to ResponseWriter
// so when developper has not to bother with response ID and Version when writing response
func DefaultRWHandler(h Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		h.ServeRPC(RWWithVersion(req.Version())(RWWithID(req.ID())(rw)), req)
	})
}

// ToHTTPHandler wraps a jsonrpc.Handler into a http.Handler
func ToHTTPHandler(h Handler) http.Handler {
	h = DefaultRWHandler(h)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rpwRW := NewResponseWriter(rw)

		// extract JSON-RPC request from context
		rpcReq := RequestFromContext(req.Context())
		if rpcReq == nil {
			// if no JSON-RPC request is found then creates one and attached to http.Request context
			rpcReq = NewRequest(req)
			err := rpcReq.ReadBody()
			if err != nil {
				_ = WriteError(rpwRW, ParseError(err))
				return
			}
			rpcReq.req = req.WithContext(WithRequest(req.Context(), rpcReq))
		} else {
			// if found update http.Request
			rpcReq.req = req
		}

		// Serve
		h.ServeRPC(rpwRW, rpcReq)
	})
}

// FromHTTPHandler wraps a http.Handler into a jsonrpc.Handler
func FromHTTPHandler(h http.Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		// Write JSON-RPC request message into request body
		_ = req.WriteBody()

		// Serve HTTP request
		h.ServeHTTP(rw.(common.WriterWrapper).Writer().(http.ResponseWriter), req.Request())
	})
}

func NotSupportedVersion(rw ResponseWriter, req *Request) {
	_ = WriteError(rw, NotSupporteVersionError(req.Version()))
}

// NotSupportedVersionHandler returns a simple handler
// that replies to each request with a not supported version request error
func NotSupportedVersionHandler() Handler { return HandlerFunc(NotSupportedVersion) }

// InvalidMethod replies to the request with an invalid method error
func InvalidMethod(rw ResponseWriter, req *Request) {
	_ = WriteError(rw, InvalidMethodError(req.Method()))
}

// InvalidMethod returns a simple handler
// that replies to each request with an invalid method error
func InvalidMethodHandler() Handler { return HandlerFunc(InvalidMethod) }

// MethodNotFound replies to the request with a method not found error
func MethodNotFound(rw ResponseWriter, req *Request) {
	_ = WriteError(rw, MethodNotFoundError())
}

// InvalidMethod returns a simple handler
// that replies to each request with an invalid method error
func MethodNotFoundHandler() Handler { return HandlerFunc(MethodNotFound) }

// NotImplementedMethod replies to the request with an not implemented error
func NotImplementedMethod(rw ResponseWriter, req *Request) {
	_ = WriteError(rw, NotImplementedMethodError(req.Method()))
}

// NotImplementedMethodHandler returns a simple handler
// that replies to each request with an invalid method error
func NotImplementedMethodHandler() Handler { return HandlerFunc(NotImplementedMethod) }

// InvalidParamsHandler returns a simple handler
// that replies to each request with an invalid parameters error
func InvalidParamsHandler(err error) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		_ = WriteError(rw, InvalidParamsError(err))
	})
}

// LoggedHandler
func LoggedHandler(h Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		log.FromContext(req.Request().Context()).
			WithField("version", req.Version()).
			WithField("id", fmt.Sprintf("%s", req.ID())).
			WithField("method", req.Method()).
			Info("serve JSON-RPC request")
		h.ServeRPC(rw, req)
	})
}
