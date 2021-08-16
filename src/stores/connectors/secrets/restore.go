package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Restore(ctx context.Context, id, version string) error {
	logger := c.logger.With("id", id).With("version", version)
	logger.Debug("restoring secret")

	secret, err := c.db.GetDeleted(ctx, id, version)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Secrets) error {
		err = c.db.Restore(ctx, secret.ID, version)
		if err != nil {
			return err
		}

		err = c.store.Restore(ctx, secret.ID, version)
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
