package stores

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authorizator"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/stores"
	eth "github.com/consensys/quorum-key-manager/src/stores/connectors/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/keys"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/secrets"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Connector) GetSecretStore(_ context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.SecretStore, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if storeBundle, ok := c.secrets[storeName]; ok {
		permissions := c.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, storeBundle.logger)

		if err := resolver.CheckAccess(storeBundle.manifest.AllowedTenants); err != nil {
			return nil, err
		}

		if store, ok := storeBundle.store.(stores.SecretStore); ok {
			return secrets.NewConnector(store, c.db.Secrets(storeName), resolver, storeBundle.logger), nil
		}
	}

	errMessage := "secret store was not found"
	c.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}

func (c *Connector) GetKeyStore(_ context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.KeyStore, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if storeBundle, ok := c.keys[storeName]; ok {
		permissions := c.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, storeBundle.logger)

		if err := resolver.CheckAccess(storeBundle.manifest.AllowedTenants); err != nil {
			return nil, err
		}

		if store, ok := storeBundle.store.(stores.KeyStore); ok {
			return keys.NewConnector(store, c.db.Keys(storeName), resolver, storeBundle.logger), nil
		}
	}

	errMessage := "key store was not found"
	c.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}

func (c *Connector) GetEthStore(ctx context.Context, name string, userInfo *authtypes.UserInfo) (stores.EthStore, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.getEthStore(ctx, name, userInfo)
}

func (c *Connector) GetEthStoreByAddr(ctx context.Context, addr common.Address, userInfo *authtypes.UserInfo) (stores.EthStore, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	for _, storeName := range c.list(ctx, manifest.Ethereum, userInfo) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			acc, err := c.getEthStore(ctx, storeName, userInfo)
			if err != nil {
				return nil, err
			}

			_, err = acc.Get(ctx, addr)
			if err == nil {
				// CheckPermission if account exists in store and returns it
				_, err := acc.Get(ctx, addr)
				if err == nil {
					return acc, nil
				}
				return acc, nil
			}
		}
	}

	errMessage := "account was not found"
	c.logger.Error(errMessage, "account", addr.Hex())
	return nil, errors.InvalidParameterError(errMessage)
}

func (c *Connector) getEthStore(_ context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.EthStore, error) {
	if storeBundle, ok := c.ethAccounts[storeName]; ok {
		permissions := c.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, storeBundle.logger)

		if err := resolver.CheckAccess(storeBundle.manifest.AllowedTenants); err != nil {
			return nil, err
		}

		if store, ok := storeBundle.store.(stores.KeyStore); ok {
			return eth.NewConnector(store, c.db.ETHAccounts(storeName), resolver, storeBundle.logger), nil
		}
	}

	errMessage := "account store was not found"
	c.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}
