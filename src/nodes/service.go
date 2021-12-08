package nodes

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock

// Nodes Service allows managing nodes
type Nodes interface {
	// Create creates a new node
	Create(ctx context.Context, name string, config *proxynode.Config, allowedTenants []string, userInfo *entities.UserInfo) error

	// Get returns a node by name
	Get(ctx context.Context, name string, userInfo *entities.UserInfo) (*proxynode.Node, error)

	// List returns a list of nodes
	List(ctx context.Context, userInfo *entities.UserInfo) ([]string, error)
}
