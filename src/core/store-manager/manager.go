package storemanager

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
)

//go:generate mockgen -source=manager.go -destination=mocks/manager.go -package=mocks

// StoreManager allows to manage multiple stores
type StoreManager interface {
	// Load manifests and performs associated actions (such as creating stores)
	// If any error occurs it is attached to the corresponding Message
	Load(ctx context.Context, mnfsts ...*manifest.Manifest) error

	// GetSecretStore by name
	GetSecretStore(ctx context.Context, name string) (secrets.SecretStore, error)

	// GetKeyStore by name
	GetKeyStore(ctx context.Context, name string) (keys.KeyStore, error)

	// GetAccountStore by name
	GetAccountStore(ctx context.Context, name string) (accounts.Store, error)

	// List stores
	List(ctx context.Context, kind manifest.Kind) ([]string, error)
}
