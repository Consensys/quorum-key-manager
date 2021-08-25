package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

func (c Connector) List(ctx context.Context) ([]string, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	ids := []string{}
	items, err := c.db.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		ids = append(ids, item.ID)
	}

	c.logger.Debug("secrets listed successfully")
	return ids, nil
}

func (c Connector) ListDeleted(ctx context.Context) ([]string, error) {
	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionRead, Resource: types.ResourceSecret})
	if err != nil {
		return nil, err
	}

	ids := []string{}
	items, err := c.db.GetAllDeleted(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		ids = append(ids, item.ID)
	}

	c.logger.Debug("deleted secrets listed successfully")
	return ids, nil
}
