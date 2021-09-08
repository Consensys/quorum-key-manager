package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

func (c Connector) List(ctx context.Context, limit, offset uint64) ([]string, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey})
	if err != nil {
		return nil, err
	}

	ids, err := c.db.SearchIDs(ctx, false, limit, offset)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("keys listed successfully")
	return ids, nil
}

func (c Connector) ListDeleted(ctx context.Context, limit, offset uint64) ([]string, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey})
	if err != nil {
		return nil, err
	}

	ids, err := c.db.SearchIDs(ctx, true, limit, offset)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("deleted keys listed successfully")
	return ids, nil
}
