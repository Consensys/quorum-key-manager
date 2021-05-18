package storemanager

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=manager.go -destination=mocks/manager.go -package=mock

// StoreManager allows to manage multiple stores
type StoreManager interface {
	// Load manifests and performs associated actions (such as creating stores)
	// If any error occurs it is attached to the corresponding Message
	Load(ctx context.Context, mnfsts ...*manifest.Manifest) error

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
