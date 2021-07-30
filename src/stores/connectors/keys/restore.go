package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/stores/store/database"
)

func (c Connector) Restore(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("restoring key")

	key, err := c.db.Keys().GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		derr := c.db.Keys().Restore(ctx, key)
		if derr != nil {
			return derr
		}

		err = c.store.Undelete(ctx, id)
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
