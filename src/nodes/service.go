package nodes

import (
	"context"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/nodes/node"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock

// Service allows managing nodes
type Service interface {
	// Node returns a node by name
	Node(ctx context.Context, name string, userInfo *authtypes.UserInfo) (node.Node, error)

	// List returns a list of nodes
	List(ctx context.Context, userInfo *authtypes.UserInfo) ([]string, error)
}
