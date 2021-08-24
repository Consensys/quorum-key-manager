package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Restore(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("restoring key")

	err := c.authorizator.Check(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceKey})
	if err != nil {
		return err
	}

	_, err = c.db.GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Keys) error {
		err = c.db.Restore(ctx, id)
		if err != nil {
			return err
		}

		err = c.store.Restore(ctx, id)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support restoring, we only restore in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("key restored successfully")
	return nil
}
