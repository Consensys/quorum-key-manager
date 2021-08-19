package secrets

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := c.logger.With("id", id)
	logger.Debug("creating secret")

	err := c.authorizator.Check(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	secret, err := c.store.Set(ctx, id, value, attr)
	if err != nil {
		return nil, err
	}

	_, err = c.db.Add(ctx, secret)
	if err != nil {
		// @TODO Ensure secret is destroyed if we fail to insert in DB
		return nil, err
	}

	logger.Info("secret created successfully", "version", secret.Metadata.Version)
	return secret, nil
}
