package connectors

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/policy"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
)

type KeyConnector struct {
	store    keys.Store
	logger   log.Logger
	resolver *policy.Resolver
}

var _ keys.Store = KeyConnector{}

func NewKeyConnector(store keys.Store, resolvr *policy.Resolver, logger log.Logger) *KeyConnector {
	return &KeyConnector{
		store:    store,
		logger:   logger,
		resolver: resolvr,
	}
}


func (c KeyConnector) Info(ctx context.Context) (*entities.StoreInfo, error) {
	result, err := c.store.Info(ctx)
	if err != nil {
		c.logger.WithError(err).Error("failed to fetch keystore info")
		return nil, err
	}
	
	c.logger.Debug("fetched keystore info successfully")
	return result, nil
}

func (c KeyConnector) Create(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := c.logger.With("id", id).With("algorithm", alg.Type).With("curve", alg.EllipticCurve)
	logger.Debug("creating key")
	result, err := c.store.Create(ctx, id, alg, attr)
	if err != nil {
		logger.WithError(err).Error("failed to create key")
		return nil, err
	}
	
	logger.Info("key created successfully")
	return result, nil
}

func (c KeyConnector) Import(ctx context.Context, id string, privKey []byte, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := c.logger.With("id", id).With("algorithm", alg.Type).With("curve", alg.EllipticCurve)
	logger.Debug("importing key")
	
	result, err := c.store.Import(ctx, id, privKey, alg, attr)
	if err != nil {
		logger.WithError(err).Error("failed to import key")
		return nil, err
	}
	
	logger.Info("key imported successfully")
	return result, nil
}

func (c KeyConnector) Get(ctx context.Context, id string) (*entities.Key, error) {
	logger := c.logger.With("id", id)
	result, err := c.store.Get(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to get key")
		return nil, err
	}
	
	logger.Debug("key retrieved successfully")
	return result, nil
}

func (c KeyConnector) List(ctx context.Context) ([]string, error) {
	result, err := c.store.List(ctx)
	if err != nil {
		c.logger.WithError(err).Error("failed to list keys")
		return nil, err
	}
	
	c.logger.Debug("keys listed successfully")
	return result, nil
}

func (c KeyConnector) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	logger := c.logger.With("id", id)
	logger.Debug("updating key")
	
	result, err := c.store.Update(ctx, id, attr)
	if err != nil {
		logger.WithError(err).Error("failed to update key")
		return nil, err
	}
	
	logger.Info("key updated successfully")
	return result, nil
}

func (c KeyConnector) Delete(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("deleting key")
	
	err := c.store.Delete(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to delete key")
		return err
	}
	
	logger.Info("key deleted successfully")
	return nil
}

func (c KeyConnector) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	logger := c.logger.With("id", id)
	
	result, err := c.store.GetDeleted(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to get deleted key")
		return nil, err
	}
	
	logger.Debug("deleted key retrieved successfully")
	return result, nil
}

func (c KeyConnector) ListDeleted(ctx context.Context) ([]string, error) {
	result, err := c.store.ListDeleted(ctx)
	if err != nil {
		c.logger.Error("failed to list deleted keys")
		return nil, err
	}
	
	c.logger.Debug("deleted keys listed successfully")
	return result, nil
}

func (c KeyConnector) Undelete(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("restoring key")
	
	err := c.store.Undelete(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to restore key")
		return err
	}
	
	logger.Info("key restored successfully")
	return nil
}

func (c KeyConnector) Destroy(ctx context.Context, id string) error {
	logger := c.logger.With("id", id)
	logger.Debug("destroying key")
	
	err := c.store.Destroy(ctx, id)
	if err != nil {
		c.logger.WithError(err).Error("failed to permanently delete key")
		return err
	}
	
	logger.Info("key was permanently deleted")
	return nil
}

func (c KeyConnector) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := c.logger.With("id", id)
	
	result, err := c.store.Sign(ctx, id, data)
	if err != nil {
		c.logger.WithError(err).Error("failed to sign payload")
		return nil, err
	}
	
	logger.Debug("payload signed successfully")
	return result, nil
}

func (c KeyConnector) Verify(ctx context.Context, pubKey, data, sig []byte, algo *entities.Algorithm) error {
	err := c.store.Verify(ctx, pubKey, data, sig, algo)
	if err != nil {
		c.logger.WithError(err).Error("failed to verify data")
		return err
	}
	
	c.logger.Debug("data verified successfully")
	return nil
}

func (c KeyConnector) Encrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := c.logger.With("id", id)
	result, err := c.store.Encrypt(ctx, id, data)
	if err != nil {
		c.logger.WithError(err).Error("failed to encrypt data")
		return nil, err
	}
	
	logger.Debug("data encrypted successfully")
	return result, nil
}

func (c KeyConnector) Decrypt(ctx context.Context, id string, data []byte) ([]byte, error) {
	logger := c.logger.With("id", id)
	result, err := c.store.Decrypt(ctx, id, data)
	if err != nil {
		c.logger.WithError(err).Error("failed to decrypt data")
		return nil, err
	}
	
	logger.Debug("data decrypted successfully")
	return result, nil
}


