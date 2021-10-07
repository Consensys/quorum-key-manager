package stores

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	auth "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=stores.go -destination=mock/stores.go -package=mock

type Stores interface {
	// Create create a store given a manifest
	Create(ctx context.Context, mnf *manifest.Manifest) error

	// GetSecretStore get secret store by name
	GetSecretStore(ctx context.Context, storeName string, userInfo *auth.UserInfo) (SecretStore, error)

	// GetKeyStore get key store by name
	GetKeyStore(ctx context.Context, storeName string, userInfo *auth.UserInfo) (KeyStore, error)

	// GetEthStore get ethereum store by name
	GetEthStore(ctx context.Context, storeName string, userInfo *auth.UserInfo) (EthStore, error)

	// GetEthStoreByAddr gets ethereum store by address
	GetEthStoreByAddr(ctx context.Context, addr common.Address, userInfo *auth.UserInfo) (EthStore, error)

	// List stores
	List(ctx context.Context, kind manifest.Kind, userInfo *auth.UserInfo) ([]string, error)

	// ListAllAccounts list all accounts from all stores
	ListAllAccounts(ctx context.Context, userInfo *auth.UserInfo) ([]common.Address, error)
}
