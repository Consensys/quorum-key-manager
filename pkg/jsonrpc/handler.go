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

// DefaultRWHandler is an utility middleware that attaches request ID and Version to ResponseWriter
// so when developper has not to bother with response ID and Version when writing response
func DefaultRWHandler(h Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, msg *RequestMsg) {
		h.ServeRPC(RWWithVersion(msg.Version)(RWWithID(msg.ID)(rw)), msg)
	})
}

func NotSupportedVersion(rw ResponseWriter, msg *RequestMsg) {
	_ = WriteError(rw, NotSupporteVersionError(msg.Version))
}

// NotSupportedVersionHandler returns a simple handler
// that replies to each request with a not supported version request error
func NotSupportedVersionHandler() Handler { return HandlerFunc(NotSupportedVersion) }

// InvalidMethod replies to the request with an invalid method error
func InvalidMethod(rw ResponseWriter, msg *RequestMsg) {
	_ = WriteError(rw, InvalidMethodError(msg.Method))
}

// InvalidMethod returns a simple handler
// that replies to each request with an invalid method error
func InvalidMethodHandler() Handler { return HandlerFunc(InvalidMethod) }

// MethodNotFound replies to the request with a method not found error
func MethodNotFound(rw ResponseWriter, msg *RequestMsg) {
	_ = WriteError(rw, MethodNotFoundError())
}

// InvalidMethod returns a simple handler
// that replies to each request with an invalid method error
func MethodNotFoundHandler() Handler { return HandlerFunc(MethodNotFound) }

// NotImplementedMethod replies to the request with an not implemented error
func NotImplementedMethod(rw ResponseWriter, msg *RequestMsg) {
	_ = WriteError(rw, NotImplementedMethodError(msg.Method))
}

// NotImplementedMethodHandler returns a simple handler
// that replies to each request with an invalid method error
func NotImplementedMethodHandler() Handler { return HandlerFunc(NotImplementedMethod) }

// InvalidParamsHandler returns a simple handler
// that replies to each request with an invalid parameters error
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
			WithField("id", fmt.Sprintf("%v", msg.ID)).
			WithField("method", msg.Method).
			Info("serve JSON-RPC request")
		h.ServeRPC(rw, msg)
	})
}
