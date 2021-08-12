package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := c.logger.With("id", id)

	var err error
	_, err = c.db.Get(ctx, id, version)
	if err != nil {
		return nil, err
	}

	secret, err := c.store.Get(ctx, id, version)
	if err != nil {
		return nil, err
	}

	logger.Debug("secret retrieved successfully")
	return secret, nil
}

func (c Connector) GetDeleted(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := c.logger.With("id", id).With("version", version)

	var err error
	dbSecret, err := c.db.GetDeleted(ctx, id, version)
	if err != nil {
		return nil, err
	}

	// Obtain secret value from vault
	secret, err := c.store.GetDeleted(ctx, id, dbSecret.Metadata.Version)
	if err != nil {
		if !errors.IsNotSupportedError(err) {
			return nil, err
		}

		// In case it is not supported we use DB value
		secret = dbSecret
	}

	logger.Debug("deleted secret retrieved successfully")
	return secret, nil
}
