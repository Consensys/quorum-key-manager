package postgres

import (
	"context"
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"

	"github.com/consensys/quorum-key-manager/src/stores/database"
)

type Secrets struct {
	storeID string
	logger  log.Logger
	client  postgres.Client
}

var _ database.Secrets = &Secrets{}

func NewSecrets(storeID string, db postgres.Client, logger log.Logger) *Secrets {
	return &Secrets{
		storeID: storeID,
		logger:  logger,
		client:  db,
	}
}

func (s Secrets) RunInTransaction(ctx context.Context, persist func(dbtx database.Secrets) error) error {
	return s.client.RunInTransaction(ctx, func(dbTx postgres.Client) error {
		s.client = dbTx
		return persist(&s)
	})
}

func (s *Secrets) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	var err error
	if version == "" {
		version, err = s.GetLatestVersion(ctx, id, false)
		if err != nil {
			return nil, err
		}
	}

	item := &models.Secret{ID: id, Version: version, StoreID: s.storeID}
	err = s.client.SelectPK(ctx, item)
	if err != nil {
		errMessage := "failed to get secrets"
		s.logger.With("id", id).With("version", version).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return item.ToEntity(), nil
}

func (s *Secrets) GetDeleted(ctx context.Context, id string) (*entities.Secret, error) {
	version, err := s.GetLatestVersion(ctx, id, true)
	if err != nil {
		return nil, err
	}

	item := &models.Secret{ID: id, Version: version, StoreID: s.storeID}
	err = s.client.SelectDeletedPK(ctx, item)
	if err != nil {
		errMessage := "failed to get deleted secrets"
		s.logger.With("id", id).With("version", version).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return item.ToEntity(), nil
}

func (s *Secrets) GetLatestVersion(ctx context.Context, id string, isDeleted bool) (string, error) {
	var version string
	var err error
	if !isDeleted {
		err = s.client.QueryOne(ctx, &version,
			"SELECT version FROM secrets WHERE id = ? AND store_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1", id, s.storeID)
	} else {
		err = s.client.QueryOne(ctx, &version,
			"SELECT version FROM secrets WHERE id = ? AND store_id = ? AND deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT 1", id, s.storeID)
	}

	if err != nil {
		errMessage := "failed to get latest secret version"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return "", errors.FromError(err).SetMessage(errMessage)
	}

	return version, nil
}

func (s *Secrets) ListVersions(ctx context.Context, id string, isDeleted bool) ([]string, error) {
	var versions = []string{}
	var err error
	if !isDeleted {
		err = s.client.Query(ctx, &versions,
			"SELECT array_agg(version ORDER BY created_at ASC) FROM secrets WHERE id = ? AND store_id = ? AND deleted_at IS NULL", id, s.storeID)
	} else {
		err = s.client.Query(ctx, &versions,
			"SELECT array_agg(version ORDER BY created_at ASC) FROM secrets WHERE id = ? AND store_id = ? AND deleted_at IS NOT NULL", id, s.storeID)
	}

	if err != nil {
		errMessage := "failed to list secret versions"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return versions, nil
}

func (s *Secrets) GetAll(ctx context.Context) ([]*entities.Secret, error) {
	var itemModels []*models.Secret

	err := s.client.SelectWhere(ctx, &itemModels, "store_id = ?", s.storeID)
	if err != nil {
		errMessage := "failed to get all secrets"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var items []*entities.Secret
	for _, item := range itemModels {
		items = append(items, item.ToEntity())
	}

	return items, nil
}

func (s *Secrets) GetAllDeleted(ctx context.Context) ([]*entities.Secret, error) {
	var itemModels []*models.Secret

	err := s.client.SelectDeletedWhere(ctx, &itemModels, "store_id = ?", s.storeID)
	if err != nil {
		errMessage := "failed to get all deleted secrets"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var items []*entities.Secret
	for _, key := range itemModels {
		items = append(items, key.ToEntity())
	}

	return items, nil
}

func (s *Secrets) Add(ctx context.Context, secret *entities.Secret) (*entities.Secret, error) {
	itemModel := models.NewSecret(secret)
	itemModel.StoreID = s.storeID
	itemModel.CreatedAt = time.Now()
	itemModel.UpdatedAt = time.Now()

	err := s.client.Insert(ctx, itemModel)
	if err != nil {
		errMessage := "failed to add secret"
		s.logger.With("id", itemModel.ID).With("version", itemModel.Version).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return itemModel.ToEntity(), nil
}

func (s *Secrets) Update(ctx context.Context, secret *entities.Secret) (*entities.Secret, error) {
	itemModel := models.NewSecret(secret)
	itemModel.StoreID = s.storeID
	itemModel.UpdatedAt = time.Now()

	err := s.client.UpdatePK(ctx, itemModel)
	if err != nil {
		errMessage := "failed to update secret"
		s.logger.With("id", itemModel.ID).With("version", itemModel.Version).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return itemModel.ToEntity(), nil
}

func (s *Secrets) Delete(ctx context.Context, id string) error {
	err := s.client.DeleteWhere(ctx, &models.Secret{ID: id, StoreID: s.storeID}, "id = ?", id)
	if err != nil {
		errMessage := "failed to delete secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Secrets) Restore(ctx context.Context, id string) error {
	err := s.client.UndeleteWhere(ctx, &models.Secret{}, "id = ? AND store_id = ?", id, s.storeID)
	if err != nil {
		errMessage := "failed to restore secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Secrets) Purge(ctx context.Context, id string) error {
	err := s.client.ForceDeleteWhere(ctx, &models.Secret{ID: id, StoreID: s.storeID}, "id = ?", id)
	if err != nil {
		errMessage := "failed to permanently delete secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}
