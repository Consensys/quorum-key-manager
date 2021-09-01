package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := c.logger.With("id", id, "version", version)

	if version == "" {
		var err error
		version, err = c.db.GetLatestVersion(ctx, id, false)
		errMsg := "failed to fetch latest secret version"
		logger.WithError(err).Error(errMsg)
		return nil, errors.FromError(err).SetMessage(errMsg)
	}

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret})
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

func (c Connector) GetDeleted(ctx context.Context, id string) (*entities.Secret, error) {
	logger := c.logger.With("id", id)

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	secret, err := c.db.GetDeleted(ctx, id, "")
	if err != nil {
		return nil, err
	}

	logger.Debug("deleted secret retrieved successfully")
	return secret, nil
}
