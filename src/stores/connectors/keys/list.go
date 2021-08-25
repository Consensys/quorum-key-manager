package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

func (c Connector) List(ctx context.Context) ([]string, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey})
	if err != nil {
		return nil, err
	}

	ids := []string{}
	keysRetrieved, err := c.db.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range keysRetrieved {
		ids = append(ids, key.ID)
	}

	c.logger.Debug("keys listed successfully")
	return ids, nil
}

func (c Connector) ListDeleted(ctx context.Context) ([]string, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceKey})
	if err != nil {
		return nil, err
	}

	ids := []string{}
	keysRetrieved, err := c.db.GetAllDeleted(ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range keysRetrieved {
		ids = append(ids, key.ID)
	}

	c.logger.Debug("deleted keys listed successfully")
	return ids, nil
}
