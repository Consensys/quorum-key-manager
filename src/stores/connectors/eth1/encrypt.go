package eth1

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (c Connector) Encrypt(ctx context.Context, addr ethcommon.Address, data []byte) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	acc, err := c.db.Get(ctx, addr.Hex())
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
