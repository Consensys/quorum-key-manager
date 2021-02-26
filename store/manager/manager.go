package manager

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/store"
)

// Manager allows to manage multiple stores
type Manager interface {
	// Load manifests and performs associated actions (such as creating stores)
	// If any error occurs it is attached to the corresponding Message
	Load(ctx context.Context, mnfsts ...*manifest.Message)

	// Get store by name
	Get(ctx context.Context, name string) (store.Store, error)

	// List stores
	List(ctx context.Context, kind string) ([]store.Store, error)
}
