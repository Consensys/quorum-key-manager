package jsonrpc

import (
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
		rpcReq := RequestFromContext(req.Context())
		if rpcReq == nil {
			rpcReq = NewRequest(req)
			rpcReq.req = req.WithContext(WithRequest(req.Context(), rpcReq))
		}

		rpcRw, ok := rw.(ResponseWriter)
		if !ok {
			rpcRw = NewResponseWriter(rw)
		}

		h.ServeRPC(rpcRw.WithVersion(rpcReq.Version()).WithID(rpcReq.ID()), rpcReq)
	})
}

// FromHTTPHandler wraps a http.Handler into a jsonrpc.Handler
func FromHTTPHandler(h http.Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		h.ServeHTTP(rw, req.Request())
	})
}
