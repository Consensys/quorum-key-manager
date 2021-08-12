package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Destroy(ctx context.Context, id, version string) error {
	logger := c.logger.With("id", id).With("id", id).With("version", version)
	logger.Debug("destroying key")

	_, err := c.db.GetDeleted(ctx, id, version)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Secrets) error {
		err = c.db.Purge(ctx, id, version)
		if err != nil {
			return err
		}

		err = c.store.Destroy(ctx, id, version)
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
