package nodemanager

import (
	"context"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/nodes/node"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows managing multiple stores
type Manager interface {
	// Node return by name
	Node(ctx context.Context, name string, userInfo *authtypes.UserInfo) (node.Node, error)

	// List stores
	List(ctx context.Context, userInfo *authtypes.UserInfo) ([]string, error)
}
