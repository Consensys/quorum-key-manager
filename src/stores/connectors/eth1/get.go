package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Get(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	logger := c.logger.With("address", addr)

	acc, err := c.db.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	logger.Debug("ethereum account retrieved successfully")
	return acc, nil
}

func (c Connector) GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error) {
	logger := c.logger.With("address", addr)

	acc, err := c.db.GetDeleted(ctx, addr)
	if err != nil {
		return nil, err
	}

	logger.Debug("deleted ethereum account retrieved successfully")
	return acc, nil
}
