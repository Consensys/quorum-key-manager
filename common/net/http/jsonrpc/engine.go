package jsonrpc

import (
	"net/http"
	"sync"
)

// Engine is an HTTP handler capable of serving JSON-RPC request
type Engine struct {
	ctxPool *sync.Pool

	handler Handler
}

// NewEngine creates a new engine capable to serve JSON-RPC requests
func NewEngine(h Handler) *Engine {
	return &Engine{
		ctxPool: &sync.Pool{
			New: func() interface{} {
				return new(Context)
			},
		},
		handler: h,
	}
}

// ServeHTTP matches http.Handler interface
func (engine *Engine) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hctx := engine.ctxPool.Get().(*Context)
	hctx.reset(rw, req)

	engine.handler.Serve(hctx)

	engine.ctxPool.Put(hctx)
}
