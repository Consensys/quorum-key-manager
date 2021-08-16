package secrets

import (
	"context"
)

func (c Connector) List(ctx context.Context) ([]string, error) {
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
