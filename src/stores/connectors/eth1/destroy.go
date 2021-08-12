package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Destroy(ctx context.Context, addr string) error {
	logger := c.logger.With("address", addr)
	logger.Debug("destroying ethereum account")

	acc, err := c.db.GetDeleted(ctx, addr)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.ETH1Accounts) error {
		err = c.db.Purge(ctx, addr)
		if err != nil {
			return err
		}

		err = c.store.Destroy(ctx, acc.KeyID)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})

	logger.Info("ethereum account was permanently deleted")
	return nil
}
