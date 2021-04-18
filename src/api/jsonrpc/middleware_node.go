package jsonrpcapi

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	nodemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/node-manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
)

// NodeMiddleware is a JSON-RPC middleware that attaches
// a node Session to every incoming request
type NodeMiddleware struct {
	mngr nodemanager.Manager
}

// NewNodeMiddleware creates a Node Middleware
func NewNodeMiddleware(mngr nodemanager.Manager) *NodeMiddleware {
	return &NodeMiddleware{
		mngr: mngr,
	}
}

func (m *NodeMiddleware) ServeRPC(rw jsonrpc.ResponseWriter, req *jsonrpc.Request, next jsonrpc.Handler) {
	// so far we support we use only a single default node
	n, err := m.mngr.Node(req.Request().Context(), "default")
	if err != nil {
		_ = rw.WriteError(err)
		return
	}

	// Get node session for the request
	session, err := n.Session(req)
	if err != nil {
		_ = rw.WriteError(err)
		return
	}

	// Execute next handler with session attached to request context
	next.ServeRPC(rw, req.WithContext(node.WithSession(req.Request().Context(), session)))
}

func (m *NodeMiddleware) Next(h jsonrpc.Handler) jsonrpc.Handler {
	return jsonrpc.HandlerFunc(func(rw jsonrpc.ResponseWriter, req *jsonrpc.Request) {
		m.ServeRPC(rw, req, h)
	})
}
