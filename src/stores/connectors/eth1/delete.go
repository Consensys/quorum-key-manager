package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func (c Connector) Delete(ctx context.Context, addr string) error {
	logger := c.logger.With("address", addr)
	logger.Debug("deleting ethereum account")

	acc, err := c.db.Get(ctx, addr)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.ETH1Accounts) error {
		err = c.db.Delete(ctx, addr)
		if err != nil {
			return err
		}

		err = c.store.Delete(ctx, acc.KeyID)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})

	logger.Info("ethereum account deleted successfully")
	return nil
}
