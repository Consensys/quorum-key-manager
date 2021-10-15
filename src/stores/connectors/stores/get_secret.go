package stores

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authorizator"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/secrets"
)

func (c *Connector) GetSecrets(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.SecretStore, error) {
	permissions := c.authManager.UserPermissions(userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

	store, err := c.getSecretStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("secret store found successfully", "store_name", storeName)
	return secrets.NewConnector(store, c.db.Secrets(storeName), resolver, c.logger), nil
}

func (c *Connector) getSecretStore(_ context.Context, storeName string, resolver auth.Authorizator) (stores.SecretStore, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if bundle, ok := c.keys[storeName]; ok {
		if err := resolver.CheckAccess(bundle.allowedTenants); err != nil {
			return nil, err
		}

		if store, ok := bundle.store.(stores.SecretStore); ok {
			return store, nil
		}
	}

	errMessage := "secret store was not found"
	c.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}
