package storemanager

import (
	"context"

	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/eth1"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/secrets"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows to manage multiple stores
type Manager interface {
	// GetSecretStore by name
	GetSecretStore(ctx context.Context, name string) (secrets.Store, error)

	// GetKeyStore by name
	GetKeyStore(ctx context.Context, name string) (keys.Store, error)

	// GetEth1Store by name
	GetEth1Store(ctx context.Context, name string) (eth1.Store, error)

	// GetEth1StoreByAddr
	GetEth1StoreByAddr(ctx context.Context, addr ethcommon.Address) (eth1.Store, error)

	// List stores
	List(ctx context.Context, kind manifest.Kind) ([]string, error)

	// ListAllAccounts list all accounts from all stores
	ListAllAccounts(context.Context) ([]*entities.ETH1Account, error)
}