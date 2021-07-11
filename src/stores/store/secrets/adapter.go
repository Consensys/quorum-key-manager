package secrets

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

type Adapter struct {
	store  Store
	logger log.Logger
}

var _ Store = Adapter{}

func NewAdapter(store Store, logger log.Logger) *Adapter {
	return &Adapter{
		store:  store,
		logger: logger,
	}
}

func (s Adapter) Info(ctx context.Context) (*entities.StoreInfo, error) {
	return s.store.Info(ctx)
}

func (s Adapter) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := s.logger.With("id", id)
	logger.Debug("creating secret")
	result, err := s.store.Set(ctx, id, value, attr)
	if err != nil {
		logger.WithError(err).Error("failed to set secret")
	}

	logger.Info("secret set successfully")
	return result, nil
}

func (s Adapter) Get(ctx context.Context, id string, version string) (*entities.Secret, error) {
	logger := s.logger.With("id", id)
	result, err := s.store.Get(ctx, id, version)
	if err != nil {
		logger.WithError(err).Error("failed to get secret")
	}

	logger.Debug("secret retrieved successfully")
	return result, nil
}

func (s Adapter) List(ctx context.Context) ([]string, error) {
	result, err := s.store.List(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to list secrets")
	}

	s.logger.Debug("secrets listed successfully")
	return result, nil
}

func (s Adapter) Delete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	err := s.store.Delete(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to delete secret")
	}

	logger.Info("secret was deleted")
	return nil
}

func (s Adapter) GetDeleted(ctx context.Context, id string) (*entities.Secret, error) {
	logger := s.logger.With("id", id)
	result, err := s.store.GetDeleted(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to get deleted secrets")
	}

	logger.Debug("retrieve delete secret successfully")
	return result, nil
}

func (s Adapter) ListDeleted(ctx context.Context) ([]string, error) {
	return s.store.ListDeleted(ctx)
}

func (s Adapter) Undelete(ctx context.Context, id string) error {
	return s.store.Undelete(ctx, id)
}

func (s Adapter) Destroy(ctx context.Context, id string) error {
	return s.store.Destroy(ctx, id)
}
