package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Restore(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("restoring secret")

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionDelete, Resource: types.ResourceSecret})
	if err != nil {
		return err
	}
	
	_, err = c.Get(ctx, id, "")
	if err == nil {
		return nil
	}

	secret, err := c.db.GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Secrets) error {
		err = dbtx.Restore(ctx, secret.ID)
		if err != nil {
			return err
		}

		err = c.store.Restore(ctx, secret.ID)
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
