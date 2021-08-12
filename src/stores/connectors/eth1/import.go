package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETH1Account, error) {
	logger := c.logger.With("id", id)
	logger.Debug("importing ethereum account")

	key, err := c.store.Import(ctx, id, privKey, eth1Algo, attr)
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Add(ctx, parseKey(key, attr))
	if err != nil {
		return nil, err
	}

	logger.With("address", acc.Address).
		With("key_id", acc.KeyID).
		Info("ethereum account imported successfully")

	return acc, nil
}
