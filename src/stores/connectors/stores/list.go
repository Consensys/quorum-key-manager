package stores

import (
	"context"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/src/auth/authorizator"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Connector) List(_ context.Context, storeType manifest.StoreType, userInfo *authtypes.UserInfo) ([]string, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	var storeNames []string
	switch storeType {
	case "":
		storeNames = append(c.listStores(c.secrets, storeType, userInfo), c.listStores(c.keys, storeType, userInfo)...)
	case manifest.HashicorpSecrets, manifest.AKVSecrets, manifest.AWSSecrets:
		storeNames = c.listStores(c.secrets, storeType, userInfo)
	case manifest.AKVKeys, manifest.HashicorpKeys, manifest.AWSKeys, manifest.Ethereum:
		storeNames = c.listStores(c.keys, storeType, userInfo)
	}

	return storeNames, nil
}

func (c *Connector) ListAllAccounts(ctx context.Context, userInfo *authtypes.UserInfo) ([]common.Address, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	var accs []common.Address
	stores, err := c.List(ctx, manifest.Ethereum, userInfo)
	if err != nil {
		return nil, err
	}

	for _, storeName := range stores {
		store, err := c.GetEthereum(ctx, storeName, userInfo)
		if err != nil {
			return nil, err
		}

		storeAccs, err := store.List(ctx, 0, 0)
		if err != nil {
			return nil, err
		}
		accs = append(accs, storeAccs...)
	}

	return accs, nil
}

func (c *Connector) listStores(list map[string]storeBundle, storeType manifest.StoreType, userInfo *authtypes.UserInfo) []string {
	var storeNames []string
	for k, storeBundle := range list {
		permissions := c.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

		if err := resolver.CheckAccess(storeBundle.allowedTenants); err != nil {
			continue
		}

		if storeType == "" || storeBundle.storeType == storeType {
			storeNames = append(storeNames, k)
		}
	}

	return storeNames
}
