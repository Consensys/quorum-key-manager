package keys

import "context"

func (c Connector) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := c.logger.With("id", id)

	result, err := c.store.Decrypt(ctx, id, data)
	if err != nil {
		return nil, err
	}

	logger.Debug("data decrypted successfully")
	return result, nil
}
