package postgres

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/store/models"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"

	"github.com/consensys/quorum-key-manager/src/stores/store/database"
)

type Keys struct {
	logger log.Logger
	db     postgres.Client
}

var _ database.Keys = &Keys{}

func NewKeys(logger log.Logger, db postgres.Client) *Keys {
	return &Keys{
		logger: logger,
		db:     db,
	}
}

func (d *Keys) Get(ctx context.Context, id string) (*models.Key, error) {
	key := &models.Key{ID: id}

	err := d.db.SelectPK(ctx, key)
	if err != nil {
		errMessage := "failed to get key"
		d.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return key, nil
}

func (d *Keys) GetDeleted(ctx context.Context, id string) (*models.Key, error) {
	key := &models.Key{ID: id}

	err := d.db.SelectDeletedPK(ctx, key)
	if err != nil {
		errMessage := "failed to get deleted key"
		d.logger.With("id", id).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return key, nil
}

func (d *Keys) GetAll(ctx context.Context) ([]*models.Key, error) {
	var keys []*models.Key

	err := d.db.Select(ctx, &keys)
	if err != nil {
		errMessage := "failed to get all keys"
		d.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return keys, nil
}

func (d *Keys) GetAllDeleted(ctx context.Context) ([]*models.Key, error) {
	var keys []*models.Key

	err := d.db.SelectDeleted(ctx, keys)
	if err != nil {
		errMessage := "failed to get all deleted keys"
		d.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return keys, nil
}

func (d *Keys) Add(ctx context.Context, key *models.Key) error {
	err := d.db.Insert(ctx, key)
	if err != nil {
		errMessage := "failed to add key"
		d.logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (d *Keys) Update(ctx context.Context, key *models.Key) error {
	err := d.db.UpdatePK(ctx, key)
	if err != nil {
		errMessage := "failed to update key"
		d.logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (d *Keys) Delete(ctx context.Context, id string) error {
	key := &models.Key{ID: id}

	err := d.db.DeletePK(ctx, key)
	if err != nil {
		errMessage := "failed to delete key"
		d.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (d *Keys) Restore(ctx context.Context, key *models.Key) error {
	key.DeletedAt = time.Time{}
	err := d.Update(ctx, key)
	if err != nil {
		errMessage := "failed to restore key"
		d.logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (d *Keys) Purge(ctx context.Context, id string) error {
	key := &models.Key{ID: id}

	err := d.db.ForceDeletePK(ctx, key)
	if err != nil {
		errMessage := "failed to permanently delete key"
		d.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}
