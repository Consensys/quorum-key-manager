package eth1

import (
	"context"
)

func (c Connector) Encrypt(ctx context.Context, addr string, data []byte) ([]byte, error) {
	logger := c.logger.With("address", addr)

	acc, err := c.db.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	result, err := c.store.Encrypt(ctx, acc.KeyID, data)
	if err != nil {
		return nil, err
	}

	logger.Debug("data encrypted successfully")
	return result, nil
}
