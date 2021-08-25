package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := c.logger.With("id", id, "version", version)

	err := c.authorizator.Check(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	secret, err := c.db.Get(ctx, id, version)
	if err != nil {
		return nil, err
	}

	secretVault, err := c.store.Get(ctx, id, version)
	if err != nil {
		return nil, err
	}
	secret.Value = secretVault.Value

	logger.Debug("secret retrieved successfully")
	return secret, nil
}

func (c Connector) GetDeleted(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := c.logger.With("id", id, "version", version)

	err := c.authorizator.Check(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	secret, err := c.db.GetDeleted(ctx, id, version)
	if err != nil {
		return nil, err
	}

	logger.Debug("deleted secret retrieved successfully")
	return secret, nil
}
