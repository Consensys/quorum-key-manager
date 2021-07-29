package keys

import "context"

func (c Connector) List(ctx context.Context) ([]string, error) {
	ids := []string{}
	keysRetrieved, err := c.db.Keys().GetAll(ctx)
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
	ids := []string{}
	keysRetrieved, err := c.db.Keys().GetAllDeleted(ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range keysRetrieved {
		ids = append(ids, key.ID)
	}

	c.logger.Debug("deleted keys listed successfully")
	return ids, nil
}
