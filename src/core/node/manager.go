package node

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
)

// Manager allows to manage multiple stores
type Manager interface {
	// Load manifest
	// If any error occurs it is attached to the corresponding Message
	Load(ctx context.Context, mnfsts ...*manifest.Manifest) error

	// GetSecretStore by name
	GetNode(ctx context.Context, name string) (Node, error)

	// List stores
	List(ctx context.Context) ([]string, error)
}
