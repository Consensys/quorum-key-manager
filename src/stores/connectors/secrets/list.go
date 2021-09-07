package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

func (c Connector) List(ctx context.Context, limit, offset int) ([]string, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	ids, err := c.db.ListIDs(ctx, limit, offset, false)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("secrets listed successfully")
	return ids, nil
}

func (c Connector) ListDeleted(ctx context.Context, limit, offset int) ([]string, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	ids, err := c.db.ListIDs(ctx, limit, offset, true)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("deleted secrets listed successfully")
	return ids, nil
}
