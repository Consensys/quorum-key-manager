package stores

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authorizator"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/secrets"
)

func (c *Connector) Secret(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.SecretStore, error) {
	permissions := c.authManager.UserPermissions(userInfo)
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

	if storeInfo.StoreType != manifest.Secrets {
		errMessage := "not a secret store"
		c.logger.Error(errMessage, "store_name", storeName)
		return nil, errors.NotFoundError(errMessage)
	}

	return storeInfo.Store.(stores.SecretStore), nil
}
