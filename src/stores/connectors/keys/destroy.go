package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

func (c Connector) Destroy(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("destroying key")

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceKey})
	if err != nil {
		return err
	}

	_, err = c.Get(ctx, id)
	if err == nil {
		logger.WithError(err).Error("try to destroy an exiting key, must be deleted first")
		return errors.StatusConflictError("key %s must be deleted first", id)
	}

	_, err = c.db.GetDeleted(ctx, id)
	if err == nil {
		err = c.db.RunInTransaction(ctx, func(dbtx database.Keys) error {
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
	}

	logger.Info("key was permanently deleted")
	return nil
}
