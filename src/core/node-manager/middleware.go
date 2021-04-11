package nodemanager

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
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
	// so far we support we use default node
	n, err := m.mngr.Node(req.Request().Context(), "")
	if err != nil {
		_ = rw.WriteError(err)
		return
	}

	// Get session for the request
	session, err := n.Session(req)
	if err != nil {
		_ = rw.WriteError(err)
		return
	}

	// Execute next handler with session attached to request context
	next.ServeRPC(rw, req.WithContext(node.WithSession(req.Request().Context(), session)))
}

func (m *Middleware) Next(h jsonrpc.Handler) jsonrpc.Handler {
	return jsonrpc.HandlerFunc(func(rw jsonrpc.ResponseWriter, req *jsonrpc.Request) {
		m.ServeRPC(rw, req, h)
	})
}
