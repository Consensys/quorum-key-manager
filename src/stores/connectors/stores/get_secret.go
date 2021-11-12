package stores

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/secrets"
)

func (c *Connector) Secret(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.SecretStore, error) {
	permissions := c.roles.UserPermissions(ctx, userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

	store, err := c.getSecretStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("secret store found successfully", "store_name", storeName)
	return secrets.NewConnector(store, c.db.Secrets(storeName), resolver, c.logger), nil
}

func (c *Connector) getSecretStore(ctx context.Context, storeName string, resolver auth.Authorizator) (stores.SecretStore, error) {
	storeInfo, err := c.getStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	if storeInfo.StoreType != entities.SecretStoreType {
		errMessage := "not a secret store"
		c.logger.Error(errMessage, "store_name", storeName)
		return nil, errors.NotFoundError(errMessage)
	}

	return storeInfo.Store.(stores.SecretStore), nil
}
