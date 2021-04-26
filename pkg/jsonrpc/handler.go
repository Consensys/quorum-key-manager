package jsonrpc

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock

// Handler is and JSON-RPC handler to be used in a JSON-RPC server
// It provides the JSON-RPC abstraction over http.Handler interface
type Handler interface {
	ServeRPC(ResponseWriter, *RequestMsg)
}

type HandlerFunc func(ResponseWriter, *RequestMsg)

func (f HandlerFunc) ServeRPC(rw ResponseWriter, msg *RequestMsg) {
	f(rw, msg)
}

// DefaultRWHandler is an utility middleware that attaches msguest ID and Version to ResponseWriter
// so when developper has not to bother with response ID and Version when writing response
func DefaultRWHandler(h Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, msg *RequestMsg) {
		h.ServeRPC(RWWithVersion(msg.Version)(RWWithID(msg.ID)(rw)), msg)
	})
}

// // ToHTTPHandler wraps a jsonrpc.Handler into a http.Handler
// func ToHTTPHandler(h Handler) http.Handler {
// 	h = DefaultRWHandler(h)
// 	return http.HandlerFunc(func(rw http.ResponseWriter, msg *http.Request) {
// 		rpwRW := NewResponseWriter(rw)

// 		// extract JSON-RPC msguest from context
// 		rpcReq := RequestFromContext(msg.Context())
// 		if rpcReq == nil {
// 			// if no JSON-RPC msguest is found then creates one and attached to http.Request context
// 			rpcReq = NewRequest(msg)
// 			err := rpcReq.ReadBody()
// 			if err != nil {
// 				_ = WriteError(rpwRW, InvalidRequest(err))
// 				return
// 			}
// 			rpcReq.msg = msg.WithContext(WithRequest(msg.Context(), rpcReq))
// 		} else {
// 			// if found update http.Request
// 			rpcReq.msg = msg
// 		}

// 		// Serve
// 		h.ServeRPC(rpwRW, rpcReq)
// 	})
// }

// // FromHTTPHandler wraps a http.Handler into a jsonrpc.Handler
// func FromHTTPHandler(h http.Handler) Handler {
// 	return HandlerFunc(func(rw ResponseWriter, msg *RequestMsg) {
// 		// Write JSON-RPC msguest message into msguest body
// 		_ = msg.WriteBody()

// 		// Serve HTTP msguest
// 		h.ServeHTTP(rw.(common.WriterWrapper).Writer().(http.ResponseWriter), msg.Request())
// 	})
// }

func NotSupportedVersion(rw ResponseWriter, msg *RequestMsg) {
	_ = WriteError(rw, NotSupporteVersionError(msg.Version))
}

// NotSupportedVersionHandler returns a simple handler
// that replies to each msguest with a not supported version msguest error
func NotSupportedVersionHandler() Handler { return HandlerFunc(NotSupportedVersion) }

// InvalidMethod replies to the msguest with an invalid method error
func InvalidMethod(rw ResponseWriter, msg *RequestMsg) {
	_ = WriteError(rw, InvalidMethodError(msg.Method))
}

// InvalidMethod returns a simple handler
// that replies to each msguest with an invalid method error
func InvalidMethodHandler() Handler { return HandlerFunc(InvalidMethod) }

// MethodNotFound replies to the msguest with a method not found error
func MethodNotFound(rw ResponseWriter, msg *RequestMsg) {
	_ = WriteError(rw, MethodNotFoundError())
}

// InvalidMethod returns a simple handler
// that replies to each msguest with an invalid method error
func MethodNotFoundHandler() Handler { return HandlerFunc(MethodNotFound) }

// NotImplementedMethod replies to the msguest with an not implemented error
func NotImplementedMethod(rw ResponseWriter, msg *RequestMsg) {
	_ = WriteError(rw, NotImplementedMethodError(msg.Method))
}

// NotImplementedMethodHandler returns a simple handler
// that replies to each msguest with an invalid method error
func NotImplementedMethodHandler() Handler { return HandlerFunc(NotImplementedMethod) }

// InvalidParamsHandler returns a simple handler
// that replies to each msguest with an invalid parameters error
func InvalidParamsHandler(err error) Handler {
	return HandlerFunc(func(rw ResponseWriter, msg *RequestMsg) {
		_ = WriteError(rw, InvalidParamsError(err))
	})
}

// LoggedHandler
func LoggedHandler(h Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, msg *RequestMsg) {
		log.FromContext(msg.Context()).
			WithField("version", msg.Version).
			WithField("id", fmt.Sprintf("%s", msg.ID)).
			WithField("method", msg.Method).
			Info("serve JSON-RPC msguest")
		h.ServeRPC(rw, msg)
	})
}
