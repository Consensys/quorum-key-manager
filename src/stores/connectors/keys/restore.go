package keys

import (
	"context"
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

		return c.store.Undelete(ctx, id)
	})
	if err != nil {
		return err
	}

	logger.Info("key restored successfully")
	return nil
}
