package stores

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authorizator"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/keys"
)

func (c *Connector) GetKeys(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.KeyStore, error) {
	permissions := c.authManager.UserPermissions(userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

	store, err := c.getKeyStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("key store found successfully", "store_name", storeName)
	return keys.NewConnector(store, c.db.Keys(storeName), resolver, c.logger), nil
}

func (c *Connector) getKeyStore(_ context.Context, storeName string, resolver auth.Authorizator) (stores.KeyStore, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if bundle, ok := c.keys[storeName]; ok {
		if err := resolver.CheckAccess(bundle.allowedTenants); err != nil {
			return nil, err
		}

		if store, ok := bundle.store.(stores.KeyStore); ok {
			return store, nil
		}
	}

	errMessage := "key store was not found"
	c.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}
