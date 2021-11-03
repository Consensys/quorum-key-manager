package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Destroy(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("permanently deleting secret")

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionDestroy, Resource: entities.ResourceSecret})
	if err != nil {
		return err
	}

	_, err = c.db.GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Secrets) error {
		err = dbtx.Purge(ctx, id)
		if err != nil {
			return err
		}

		err = c.store.Destroy(ctx, id)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("secret was permanently deleted")
	return nil
}
