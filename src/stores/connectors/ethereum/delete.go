package eth

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (c Connector) Delete(ctx context.Context, addr ethcommon.Address) error {
	logger := c.logger.With("address", addr.Hex())
	logger.Debug("deleting ethereum account")

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount})
	if err != nil {
		return err
	}

	acc, err := c.db.Get(ctx, addr.Hex())
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.ETHAccounts) error {
		err = dbtx.Delete(ctx, addr.Hex())
		if err != nil {
			return err
		}

		err = c.store.Delete(ctx, acc.KeyID)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("ethereum account deleted successfully")
	return nil
}
