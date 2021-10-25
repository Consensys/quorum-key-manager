package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Restore(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("restoring secret")

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceSecret})
	if err != nil {
		return err
	}

	// If secret already exists, exit without any action
	_, err = c.Get(ctx, id, "")
	if err == nil {
		return nil
	}

	secret, err := c.db.GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Secrets) error {
		err = dbtx.Restore(ctx, secret.ID)
		if err != nil {
			return err
		}

		err = c.store.Restore(ctx, secret.ID)
		if err != nil && !errors.IsNotSupportedError(err) {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("secret restored successfully")
	return nil
}
