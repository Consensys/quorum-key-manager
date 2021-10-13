package stores

import (
	"context"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/src/auth/authorizator"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Connector) List(ctx context.Context, kind manifest.Kind, userInfo *authtypes.UserInfo) ([]string, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.list(ctx, kind, userInfo), nil
}

func (c *Connector) ListAllAccounts(ctx context.Context, userInfo *authtypes.UserInfo) ([]common.Address, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	var accs []common.Address
	for _, storeName := range c.list(ctx, manifest.Ethereum, userInfo) {
		store, err := c.getEthStore(ctx, storeName, userInfo)
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

func (c *Connector) listStores(list map[string]*storeBundle, kind manifest.Kind, userInfo *authtypes.UserInfo) []string {
	var storeNames []string
	for k, storeBundle := range list {
		permissions := c.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, storeBundle.logger)

		if err := resolver.CheckAccess(storeBundle.manifest.AllowedTenants); err != nil {
			continue
		}

		if kind == "" || storeBundle.manifest.Kind == kind {
			storeNames = append(storeNames, k)
		}
	}

	return storeNames
}

func (c *Connector) list(_ context.Context, kind manifest.Kind, userInfo *authtypes.UserInfo) []string {
	var storeNames []string
	switch kind {
	case "":
		storeNames = append(
			append(c.listStores(c.secrets, kind, userInfo), c.listStores(c.keys, kind, userInfo)...), c.listStores(c.ethAccounts, kind, userInfo)...)
	case manifest.HashicorpSecrets, manifest.AKVSecrets, manifest.AWSSecrets:
		storeNames = c.listStores(c.secrets, kind, userInfo)
	case manifest.AKVKeys, manifest.HashicorpKeys, manifest.AWSKeys:
		storeNames = c.listStores(c.keys, kind, userInfo)
	case manifest.Ethereum:
		storeNames = c.listStores(c.ethAccounts, kind, userInfo)
	}

	return storeNames
}
