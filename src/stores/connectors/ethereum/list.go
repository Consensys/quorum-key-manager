package eth

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/ethereum/go-ethereum/common"
)

func (c Connector) List(ctx context.Context) ([]common.Address, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	accs, err := c.db.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	addrs := []common.Address{}
	for _, acc := range accs {
		addrs = append(addrs, acc.Address)
	}

	c.logger.Debug("ethereum accounts listed successfully")
	return addrs, nil
}

func (c Connector) ListDeleted(ctx context.Context) ([]common.Address, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	accs, err := c.db.GetAllDeleted(ctx)
	if err != nil {
		return nil, err
	}

	addrs := []common.Address{}
	for _, acc := range accs {
		addrs = append(addrs, acc.Address)
	}

	c.logger.Debug("deleted ethereum accounts listed successfully")
	return addrs, nil
}
