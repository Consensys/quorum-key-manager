package keys

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Sign(ctx context.Context, id string, data []byte, algo *entities.Algorithm) ([]byte, error) {
	logger := c.logger.With("id", id)

	err := c.authorizator.Check(&types.Operation{Action: types.ActionSign, Resource: types.ResourceKey})
	if err != nil {
		return nil, err
	}

	var result []byte
	if algo != nil {
		result, err = c.store.Sign(ctx, id, data, algo)
	} else {
		var key *entities.Key
		key, err = c.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		result, err = c.store.Sign(ctx, id, data, key.Algo)
	}

	if err != nil {
		return nil, err
	}

	logger.Debug("payload signed successfully")
	return result, nil
}
