package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Get(ctx context.Context, id string) (*entities.Key, error) {
	logger := c.logger.With("id", id)

	key, err := c.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Debug("key retrieved successfully")
	return key, nil
}

func (c Connector) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	logger := c.logger.With("id", id)

	key, err := c.db.GetDeleted(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Debug("deleted key retrieved successfully")
	return key, nil
}
