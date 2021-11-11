package stores

import (
	"context"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=stores.go -destination=mock/stores.go -package=mock

type Stores interface {
	// CreateEthereum creates an ethereum store
	CreateEthereum(_ context.Context, name, keyStore string, allowedTenants []string, userInfo *auth.UserInfo) error

	// CreateKey creates a key store
	CreateKey(_ context.Context, name, vault, secretStore string, allowedTenants []string, userInfo *auth.UserInfo) error

	// CreateSecret creates a secret store
	CreateSecret(_ context.Context, name, vault string, allowedTenants []string, userInfo *auth.UserInfo) error

	// ImportEthereum import ethereum accounts from the vault into an ethereum store
	ImportEthereum(ctx context.Context, name string, userInfo *auth.UserInfo) error

	// ImportKeys import keys from the vault into a key store
	ImportKeys(ctx context.Context, storeName string, userInfo *auth.UserInfo) error

	// ImportSecrets import secrets from the vault into a secret store
	ImportSecrets(ctx context.Context, storeName string, userInfo *auth.UserInfo) error

	// Secret get secret store by name
	Secret(ctx context.Context, storeName string, userInfo *auth.UserInfo) (SecretStore, error)

	// Key get key store by name
	Key(ctx context.Context, storeName string, userInfo *auth.UserInfo) (KeyStore, error)

	// Ethereum get ethereum store by name
	Ethereum(ctx context.Context, storeName string, userInfo *auth.UserInfo) (EthStore, error)

	// EthereumByAddr gets ethereum store by address
	EthereumByAddr(ctx context.Context, addr common.Address, userInfo *auth.UserInfo) (EthStore, error)

	// List stores
	List(ctx context.Context, storeType string, userInfo *auth.UserInfo) ([]string, error)

	// ListAllAccounts list all accounts from all stores
	ListAllAccounts(ctx context.Context, userInfo *auth.UserInfo) ([]common.Address, error)
}
