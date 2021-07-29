package storemanager

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/connectors"

	"github.com/consensys/quorum-key-manager/src/auth/types"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/eth1"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows to manage multiple stores
type Manager interface {
	// GetSecretStore by name
	GetSecretStore(ctx context.Context, name string, userInfo *types.UserInfo) (secrets.Store, error)

	// GetKeyStore by name
	GetKeyStore(ctx context.Context, name string, userInfo *types.UserInfo) (connectors.KeysConnector, error)

	// GetEth1Store by name
	GetEth1Store(ctx context.Context, name string, userInfo *types.UserInfo) (eth1.Store, error)

	// GetEth1StoreByAddr gets a eth1 store by address
	GetEth1StoreByAddr(ctx context.Context, addr ethcommon.Address, userInfo *types.UserInfo) (eth1.Store, error)

	// List stores
	List(ctx context.Context, kind manifest.Kind, userInfo *types.UserInfo) ([]string, error)

	// ListAllAccounts list all accounts from all stores
	ListAllAccounts(ctx context.Context, userInfo *types.UserInfo) ([]*entities.ETH1Account, error)
}
