package keys

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

func (c Connector) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := c.logger.With("id", id).With("algorithm", alg.Type).With("curve", alg.EllipticCurve)
	logger.Debug("importing key")

	key, err := c.store.Import(ctx, id, privKey, alg, attr)
	if err != nil {
		return nil, err
	}

	key, err = c.db.Keys().Add(ctx, key)
	if err != nil {
		return nil, err
	}

	logger.Info("key imported successfully")
	return key, nil
}
