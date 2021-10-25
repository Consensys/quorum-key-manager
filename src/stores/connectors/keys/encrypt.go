package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (c Connector) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := c.logger.With("id", id)

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceKey})
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
