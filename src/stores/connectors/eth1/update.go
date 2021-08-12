package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Update(ctx context.Context, addr string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	logger := c.logger.With("address", addr)
	logger.Debug("updating ethereum account")

	acc, err := c.db.Get(ctx, addr)
	if err != nil {
		return nil, err
	}
	acc.Tags = attr.Tags

	err = c.db.RunInTransaction(ctx, func(dbtx database.ETH1Accounts) error {
		var derr error
		acc, derr = c.db.Update(ctx, acc)
		if derr != nil {
			return derr
		}

		_, derr = c.store.Update(ctx, addr, attr)
		if derr != nil && !errors.IsNotSupportedError(derr) {
			return derr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Info("ethereum account updated successfully")
	return acc, nil
}
