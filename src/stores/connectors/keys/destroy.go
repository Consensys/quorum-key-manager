package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Destroy(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("destroying key")

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceKey})
	if err != nil {
		return err
	}

	_, err = c.db.GetDeleted(ctx, id)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Keys) error {
		err = dbtx.Purge(ctx, id)
		if err != nil {
			return err
		}

		err = c.store.Destroy(ctx, id)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	logger.Info("key was permanently deleted")
	return nil
}
