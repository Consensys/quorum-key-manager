package jsonrpc

import (
	"net/http"
)

type Handler interface {
	Serve(*Context)
}

type HandlerFunc func(*Context)

func (f HandlerFunc) Serve(hctx *Context) {
	f(hctx)
}

// ToHTTPHandler wraps a jsonrpc.Handler into a http.Handler
func ToHTTPHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		hctx, ok := fromRequest(req, rw)
		if !ok {
			hctx = newContext(rw, req)
		}
		h.Serve(hctx)
	})
}

// FromHTTPHandler wraps a http.Handler into a jsonrpc.Handler
func FromHTTPHandler(h http.Handler) Handler {
	return HandlerFunc(func(hctx *Context) {
		h.ServeHTTP(hctx.Writer(), hctx.Request())
	})
}
