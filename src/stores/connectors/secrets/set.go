package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := c.logger.With("id", id)
	logger.Debug("creating secret")

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	secret, err := c.store.Set(ctx, id, value, attr)
	if err != nil && errors.IsAlreadyExistsError(err) {
		secret, err = c.store.Get(ctx, id, "")
	}
	if err != nil {
		return nil, err
	}

	_, err = c.db.Add(ctx, secret)
	if err != nil {
		return nil, err
	}

	logger.Info("secret created successfully", "version", secret.Metadata.Version)
	return secret, nil
}
