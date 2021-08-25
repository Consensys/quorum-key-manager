package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Destroy(ctx context.Context, id, version string) error {
	logger := c.logger.With("id", id, "version", version)
	logger.Debug("permanently deleting secret")

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceSecret})
	if err != nil {
		return err
	}

	_, err = c.db.GetDeleted(ctx, id, version)
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
