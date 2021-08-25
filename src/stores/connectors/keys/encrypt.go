package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

func (c Connector) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := c.logger.With("id", id)

	err := c.authorizator.Check(&types.Operation{Action: types.ActionEncrypt, Resource: types.ResourceKey})
	if err != nil {
		return nil, err
	}

	result, err := c.store.Encrypt(ctx, id, data)
	if err != nil {
		return nil, err
	}

	logger.Debug("data encrypted successfully")
	return result, nil
}
