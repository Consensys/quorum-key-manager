package stores

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/keys"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
)

func (c *Connector) Key(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.KeyStore, error) {
	permissions := c.roles.UserPermissions(ctx, userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

	store, err := c.getKeyStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("key store found successfully", "store_name", storeName)
	return keys.NewConnector(store, c.db.Keys(storeName), resolver, c.logger), nil
}

func (c *Connector) getKeyStore(ctx context.Context, storeName string, resolver auth.Authorizator) (stores.KeyStore, error) {
	storeInfo, err := c.getStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	if storeInfo.StoreType != entities.KeyStoreType {
		errMessage := "not a key store"
		c.logger.Error(errMessage, "store_name", storeName)
		return nil, errors.NotFoundError(errMessage)
	}

	return storeInfo.Store.(stores.KeyStore), nil
}
