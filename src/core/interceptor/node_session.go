package interceptor

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	nodemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/node-manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
)

// NodeSessionMiddleware is a JSON-RPC middleware that attaches
// a node Session to every incoming request
type NodeSessionMiddleware struct {
	nodes nodemanager.Manager
}

// NewNodeSessionMiddleware creates a Node Middleware
func NewNodeSessionMiddleware(nodes nodemanager.Manager) *NodeSessionMiddleware {
	return &NodeSessionMiddleware{
		nodes: nodes,
	}
}

func (m *NodeSessionMiddleware) ServeRPC(rw jsonrpc.ResponseWriter, req *jsonrpc.Request, next jsonrpc.Handler) {
	// so far we support we use only a single default node
	n, err := m.nodes.Node(req.Request().Context(), "default")
	if err != nil {
		_ = jsonrpc.WriteError(rw, err)
		return
	}

	// Get node session for the request
	session, err := n.Session(req)
	if err != nil {
		_ = jsonrpc.WriteError(rw, err)
		return
	}

	// Execute next handler with session attached to request context
	next.ServeRPC(rw, req.WithContext(node.WithSession(req.Request().Context(), session)))
}

func (m *NodeSessionMiddleware) Next(h jsonrpc.Handler) jsonrpc.Handler {
	return jsonrpc.HandlerFunc(func(rw jsonrpc.ResponseWriter, req *jsonrpc.Request) {
		m.ServeRPC(rw, req, h)
	})
}
