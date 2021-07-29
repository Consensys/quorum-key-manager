package keys

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
)

func (c Connector) Destroy(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("destroying key")

	_, err := c.db.Keys().GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		derr := c.db.Keys().Purge(ctx, id)
		if derr != nil {
			return derr
		}

		return c.store.Destroy(ctx, id)
	})
	if err != nil {
		return err
	}

	logger.Info("key was permanently deleted")
	return nil
}
