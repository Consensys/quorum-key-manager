package stores

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Connector) List(_ context.Context, storeType manifest.StoreType, userInfo *authtypes.UserInfo) ([]string, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	var storeNames []string
	for k, storeInfo := range c.stores {
		if storeType != "" && storeInfo.StoreType != storeType {
			continue
		}

		permissions := c.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

		if err := resolver.CheckAccess(storeInfo.AllowedTenants); err != nil {
			continue
		}

		storeNames = append(storeNames, k)
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
		store, err := c.Ethereum(ctx, storeName, userInfo)
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
