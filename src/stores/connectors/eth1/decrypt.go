package eth1

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/ethereum/go-ethereum/common"
)

func (c Connector) Decrypt(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	logger := c.logger.With("address", addr.Hex())

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionEncrypt, Resource: types.ResourceEth1Account})
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Get(ctx, addr.Hex())
	if err != nil {
		return nil, err
	}

	result, err := c.store.Decrypt(ctx, acc.KeyID, data)
	if err != nil {
		return nil, err
	}

	logger.Debug("data decrypted successfully")
	return result, nil
}
