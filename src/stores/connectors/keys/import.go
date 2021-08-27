package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/auth/types"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c Connector) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := c.logger.With("id", id).With("algorithm", alg.Type).With("curve", alg.EllipticCurve)
	logger.Debug("importing key")

	err := c.authorizator.CheckPermission(&types.Operation{Action: types.ActionWrite, Resource: types.ResourceKey})
	if err != nil {
		return nil, err
	}
	
	if !isSupportedAlgo(alg){
		errMessage := "invalid or not supported elliptic curve and signing algorithm combination for importing"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	key, err := c.store.Import(ctx, id, privKey, alg, attr)
	if err != nil && errors.IsAlreadyExistsError(err) {
		key, err = c.store.Get(ctx, id)
	}
	if err != nil {
		return nil, err
	}

	key, err = c.db.Add(ctx, key)
	if err != nil {
		return nil, err
	}

	logger.Info("key imported successfully")
	return key, nil
}
