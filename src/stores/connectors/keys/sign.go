package keys

import "context"

func (c Connector) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := c.logger.With("id", id)

	key, err := c.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	result, err := c.store.Sign(ctx, id, data, key.Algo)
	if err != nil {
		return nil, err
	}

	logger.Debug("payload signed successfully")
	return result, nil
}
