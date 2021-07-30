package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/stores/store/database"
)

func (c Connector) Delete(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("deleting key")

	err := c.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		err := dbtx.Keys().Delete(ctx, id)
		if err != nil {
			return err
		}

		err = c.store.Delete(ctx, id)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("key deleted successfully")
	return nil
}
