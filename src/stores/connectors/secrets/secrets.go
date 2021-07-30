package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/policy"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"
)

type SecretConnector struct {
	store    secrets.Store
	logger   log.Logger
	resolver *policy.Resolver
}

var _ secrets.Store = &SecretConnector{}

func NewSecretConnector(store secrets.Store, resolvr *policy.Resolver, logger log.Logger) *SecretConnector {
	return &SecretConnector{
		store:    store,
		logger:   logger,
		resolver: resolvr,
	}
}

func (c SecretConnector) Info(ctx context.Context) (*entities.StoreInfo, error) {
	return c.store.Info(ctx)
}

func (c SecretConnector) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := c.logger.With("id", id)
	logger.Debug("creating secret")

	result, err := c.store.Set(ctx, id, value, attr)
	if err != nil {
		return nil, err
	}

	logger.Info("secret created successfully")
	return result, nil
}

func (c SecretConnector) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := c.logger.With("id", id)

	result, err := c.store.Get(ctx, id, version)
	if err != nil {
		return nil, err
	}

	logger.Debug("secret retrieved successfully")
	return result, nil
}

func (c SecretConnector) List(ctx context.Context) ([]string, error) {
	result, err := c.store.List(ctx)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("secrets listed successfully")
	return result, nil
}

func (c SecretConnector) Delete(ctx context.Context, id, version string) error {
	logger := c.logger.With("id", id).With("version", version)

	err := c.store.Delete(ctx, id, version)
	if err != nil {
		return err
	}

	logger.Info("secret deleted successfully")
	return nil
}

func (c SecretConnector) GetDeleted(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := c.logger.With("id", id).With("version", version)

	result, err := c.store.GetDeleted(ctx, id, version)
	if err != nil {
		return nil, err
	}

	logger.Debug("deleted secret retrieved successfully")
	return result, nil
}

func (c SecretConnector) ListDeleted(ctx context.Context) ([]string, error) {
	logger := c.logger.With("id")

	result, err := c.store.ListDeleted(ctx)
	if err != nil {
		return nil, err
	}

	logger.Debug("deleted secret listed successfully")
	return result, nil
}

func (c SecretConnector) Restore(ctx context.Context, id, version string) error {
	logger := c.logger.With("id", id).With("version", version)

	err := c.store.Restore(ctx, id, version)
	if err != nil {
		return err
	}

	logger.Debug("secret restored successfully")
	return nil
}

func (c SecretConnector) Destroy(ctx context.Context, id, version string) error {
	logger := c.logger.With("id", id)

	err := c.store.Destroy(ctx, id, version)
	if err != nil {
		return err
	}

	logger.Debug("secret permanently deleted")
	return nil
}
