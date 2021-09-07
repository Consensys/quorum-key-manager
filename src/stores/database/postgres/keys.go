package postgres

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/src/stores/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"

	"github.com/consensys/quorum-key-manager/src/stores/database"
)

type Keys struct {
	storeID string
	logger  log.Logger
	client  postgres.Client
}

var _ database.Keys = &Keys{}

func NewKeys(storeID string, db postgres.Client, logger log.Logger) *Keys {
	return &Keys{
		storeID: storeID,
		logger:  logger,
		client:  db,
	}
}

func (k Keys) RunInTransaction(ctx context.Context, persist func(dbtx database.Keys) error) error {
	return k.client.RunInTransaction(ctx, func(dbTx postgres.Client) error {
		k.client = dbTx
		return persist(&k)
	})
}

func (k *Keys) Get(ctx context.Context, id string) (*entities.Key, error) {
	key := &models.Key{ID: id, StoreID: k.storeID}

	err := k.client.SelectPK(ctx, key)
	if err != nil {
		errMessage := "failed to get key"
		k.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return key.ToEntity(), nil
}

func (k *Keys) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	key := &models.Key{ID: id, StoreID: k.storeID}

	err := k.client.SelectDeletedPK(ctx, key)
	if err != nil {
		errMessage := "failed to get deleted key"
		k.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return key.ToEntity(), nil
}

func (k *Keys) GetAll(ctx context.Context) ([]*entities.Key, error) {
	var keyModels []*models.Key

	err := k.client.SelectWhere(ctx, &keyModels, "store_id = ?", k.storeID)
	if err != nil {
		errMessage := "failed to get all keys"
		k.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var keys []*entities.Key
	for _, key := range keyModels {
		keys = append(keys, key.ToEntity())
	}

	return keys, nil
}

func (k *Keys) GetAllDeleted(ctx context.Context) ([]*entities.Key, error) {
	var keyModels []*models.Key

	err := k.client.SelectDeletedWhere(ctx, &keyModels, "store_id = ?", k.storeID)
	if err != nil {
		errMessage := "failed to get all deleted keys"
		k.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var keys []*entities.Key
	for _, key := range keyModels {
		keys = append(keys, key.ToEntity())
	}

	return keys, nil
}

func (s *Keys) ListIDs(ctx context.Context, isDeleted bool, limit, offset int) ([]string, error) {
	var ids = []string{}
	var err error
	var query string
	args := []interface{}{s.storeID}

	switch {
	case limit != 0 || offset != 0:
		query = fmt.Sprintf("SELECT (array_agg(id ORDER BY created_at ASC))[%d:%d] FROM keys WHERE store_id = ?", offset+1, offset+limit)
	default:
		query = "SELECT array_agg(id ORDER BY created_at ASC) FROM keys WHERE store_id = ?"
	}

	if isDeleted {
		query = fmt.Sprintf("%s AND deleted_at is NOT NULL", query)
	} else {
		query = fmt.Sprintf("%s AND deleted_at is NULL", query)
	}

	err = s.client.Query(ctx, &ids, query, args...)
	if err != nil {
		errMessage := "failed to list keys ids"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return ids, nil
}

func (k *Keys) Add(ctx context.Context, key *entities.Key) (*entities.Key, error) {
	keyModel := models.NewKey(key)
	keyModel.StoreID = k.storeID

	err := k.client.Insert(ctx, keyModel)
	if err != nil {
		errMessage := "failed to add key"
		k.logger.With("id", key.ID).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return keyModel.ToEntity(), nil
}

func (k *Keys) Update(ctx context.Context, key *entities.Key) (*entities.Key, error) {
	keyModel := models.NewKey(key)
	keyModel.StoreID = k.storeID

	err := k.client.UpdatePK(ctx, keyModel)
	if err != nil {
		errMessage := "failed to update key"
		k.logger.With("id", key.ID).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return keyModel.ToEntity(), nil
}

func (k *Keys) Delete(ctx context.Context, id string) error {
	err := k.client.DeletePK(ctx, &models.Key{ID: id, StoreID: k.storeID})
	if err != nil {
		errMessage := "failed to delete key"
		k.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (k *Keys) Restore(ctx context.Context, id string) error {
	err := k.client.UndeletePK(ctx, &models.Key{ID: id, StoreID: k.storeID})
	if err != nil {
		errMessage := "failed to restore key"
		k.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (k *Keys) Purge(ctx context.Context, id string) error {
	err := k.client.ForceDeletePK(ctx, &models.Key{ID: id, StoreID: k.storeID})
	if err != nil {
		errMessage := "failed to permanently delete key"
		k.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}
