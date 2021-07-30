package keys

import "context"

func (c Connector) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := c.logger.With("id", id)

	result, err := c.store.Encrypt(ctx, id, data)
	if err != nil {
		return nil, err
	}

	logger.Debug("data encrypted successfully")
	return result, nil
}
