package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Restore(ctx context.Context, id, version string) error {
	logger := c.logger.With("id", id, "version", version)
	logger.Debug("restoring secret")

	err := c.authorizator.Check(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceSecret})
	if err != nil {
		return err
	}

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
