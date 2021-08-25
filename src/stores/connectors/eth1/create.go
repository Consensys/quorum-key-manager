package eth1

import "C"
import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETH1Account, error) {
	logger := c.logger.With("id", id)
	logger.Debug("creating ethereum account")

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceEth1Account})
	if err != nil {
		return nil, err
	}

	key, err := c.store.Create(ctx, id, eth1Algo, attr)
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Add(ctx, newEth1Account(key, attr))
	if err != nil {
		// @TODO Ensure key is destroyed if we fail to insert in DB
		return nil, err
	}

	logger.With("address", acc.Address, "key_id", acc.KeyID).Info("ethereum account created successfully")
	return acc, nil
}
