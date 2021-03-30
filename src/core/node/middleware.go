package node

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/jsonrpc"
)

type Middleware struct {
	mngr Manager
}

func NewMiddleware(mngr Manager) *Middleware {
	return &Middleware{
		mngr: mngr,
	}
}

func (m *Middleware) ServeRPC(rw jsonrpc.ResponseWriter, req *jsonrpc.Request, next jsonrpc.Handler) {
	node, err := m.mngr.GetNode(req.Request().Context(), "") // so far we support only one node
	if err != nil {
		rw.WriteError(err)
		return
	}

	next.ServeRPC(rw, node.WrapRequest(req))
}

func (m *Middleware) Next(h Handler) Handler {
	return jsonrpc.HandlerFunc(func(rw jsonrpc.ResponseWriter, req *jsonrpc.Request) {
		m.ServeRPC(rw, req, h)
	})
}
