package keys

import (
	"context"
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

		return c.store.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	logger.Info("key deleted successfully")
	return nil
}
