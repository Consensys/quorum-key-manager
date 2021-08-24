package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (c Connector) Destroy(ctx context.Context, addr ethcommon.Address) error {
	logger := c.logger.With("address", addr.Hex())
	logger.Debug("destroying ethereum account")

	err := c.authorizator.Check(&types.Operation{Action: types.ActionDestroy, Resource: types.ResourceEth1Account})
	if err != nil {
		return err
	}

	acc, err := c.db.GetDeleted(ctx, addr.Hex())
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.ETH1Accounts) error {
		err = c.db.Purge(ctx, addr.Hex())
		if err != nil {
			return err
		}

		err = c.store.Destroy(ctx, acc.KeyID)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("ethereum account was permanently deleted")
	return nil
}
