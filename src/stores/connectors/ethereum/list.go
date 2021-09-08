package eth

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/ethereum/go-ethereum/common"
)

func (c Connector) List(ctx context.Context, limit, offset uint64) ([]common.Address, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	strAddr, err := c.db.SearchAddresses(ctx, false, limit, offset)
	if err != nil {
		return nil, err
	}

	addrs := []common.Address{}
	for _, addr := range strAddr {
		addrs = append(addrs, common.HexToAddress(addr))
	}

	c.logger.Debug("ethereum accounts listed successfully")
	return addrs, nil
}

func (c Connector) ListDeleted(ctx context.Context, limit, offset uint64) ([]common.Address, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	strAddr, err := c.db.SearchAddresses(ctx, true, limit, offset)
	if err != nil {
		return nil, err
	}

	addrs := []common.Address{}
	for _, addr := range strAddr {
		addrs = append(addrs, common.HexToAddress(addr))
	}

	c.logger.Debug("deleted ethereum accounts listed successfully")
	return addrs, nil
}
