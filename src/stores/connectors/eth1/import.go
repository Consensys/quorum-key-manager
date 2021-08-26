package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETH1Account, error) {
	logger := c.logger.With("id", id)
	logger.Debug("importing ethereum account")

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceEth1Account})
	if err != nil {
		return nil, err
	}

	key, err := c.store.Import(ctx, id, privKey, eth1Algo, attr)
	if err != nil && errors.IsAlreadyExistsError(err) {
		key, err = c.store.Get(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Add(ctx, newEth1Account(key, attr))
	if err != nil {
		return nil, err
	}

	logger.With("address", acc.Address, "key_id", acc.KeyID).Info("ethereum account imported successfully")
	return acc, nil
}
