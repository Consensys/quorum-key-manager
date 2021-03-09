package storemanager

import (
	"context"

	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest/loader"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
)

// Manager allows to manage multiple stores
type Manager interface {
	// Load manifests and performs associated actions (such as creating stores)
	// If any error occurs it is attached to the corresponding Message
	Load(ctx context.Context, mnfsts ...*manifestloader.Message)

	// GetSecretStore by name
	GetSecretStore(ctx context.Context, name string) (secrets.Store, error)

	// GetKeyStore by name
	GetKeyStore(ctx context.Context, name string) (keys.Store, error)

	// GetAccountStore by name
	GetAccountStore(ctx context.Context, name string) (accounts.Store, error)

	// List stores
	List(ctx context.Context, kind string) ([]string, error)
}
