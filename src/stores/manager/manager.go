package storemanager

import (
	"context"
	manifest2 "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/types"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	eth12 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1"
	keys2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys"
	secrets2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/secrets"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows to manage multiple stores
type Manager interface {
	// GetSecretStore by name
	GetSecretStore(ctx context.Context, name string) (secrets2.Store, error)

	// GetKeyStore by name
	GetKeyStore(ctx context.Context, name string) (keys2.Store, error)

	// GetEth1Store by name
	GetEth1Store(ctx context.Context, name string) (eth12.Store, error)

	// GetEth1StoreByAddr
	GetEth1StoreByAddr(ctx context.Context, addr ethcommon.Address) (eth12.Store, error)

	// List stores
	List(ctx context.Context, kind manifest2.Kind) ([]string, error)

	// ListAllAccounts list all accounts from all stores
	ListAllAccounts(context.Context) ([]*entities2.ETH1Account, error)
}
