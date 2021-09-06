package stores

import (
	"context"

	authtype "github.com/consensys/quorum-key-manager/src/auth/types"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows to manage multiple stores
type Manager interface {
	// GetSecretStore by name
	GetSecretStore(ctx context.Context, name string, userInfo *authtype.UserInfo) (SecretStore, error)

	// GetKeyStore by name
	GetKeyStore(ctx context.Context, name string, userInfo *authtype.UserInfo) (KeyStore, error)

	// GetEthStore by name
	GetEthStore(ctx context.Context, name string, userInfo *authtype.UserInfo) (EthStore, error)

	// GetEthStoreByAddr gets a eth store by address
	GetEthStoreByAddr(ctx context.Context, addr ethcommon.Address, userInfo *authtype.UserInfo) (EthStore, error)

	// List stores
	List(ctx context.Context, kind manifest.Kind, userInfo *authtype.UserInfo) ([]string, error)

	// ListAllAccounts list all accounts from all stores
	ListAllAccounts(ctx context.Context, userInfo *authtype.UserInfo) ([]ethcommon.Address, error)
}
