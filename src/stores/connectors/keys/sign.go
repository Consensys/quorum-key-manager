package keys

import (
	"context"

	authentities "github.com/consensys/quorum-key-manager/src/auth/entities"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Sign(ctx context.Context, id string, data []byte, algo *entities.Algorithm) ([]byte, error) {
	logger := c.logger.With("id", id)

	err := c.authorizator.CheckPermission(&authentities.Operation{Action: authentities.ActionSign, Resource: authentities.ResourceKey})
	if err != nil {
		return nil, err
	}

	if algo == nil {
		key, derr := c.db.Get(ctx, id)
		if derr != nil {
			return nil, derr
		}

		algo = key.Algo
	}

	result, err := c.store.Sign(ctx, id, data, algo)
	if err != nil {
		return nil, err
	}

	logger.Debug("payload signed successfully")
	return result, nil
}
