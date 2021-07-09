package postgres

import (
	"context"
	"time"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"

	"github.com/consensys/quorum-key-manager/src/stores/store/database"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
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

func (d *Keys) Get(_ context.Context, id string) (*entities.Key, error) {
	key := &entities.Key{ID: id}

	err := d.db.SelectPK(key)
	if err != nil {
		d.logger.WithError(err).Error("failed to get key")
		return nil, err
	}

	return key, nil
}

func (d *Keys) GetDeleted(_ context.Context, id string) (*entities.Key, error) {
	key := &entities.Key{ID: id}

	err := d.db.SelectDeletedPK(key)
	if err != nil {
		d.logger.WithError(err).Error("failed to get key")
		return nil, err
	}

	return key, nil
}

func (d *Keys) GetAll(_ context.Context) ([]*entities.Key, error) {
	var keys []*entities.Key

	err := d.db.Select(&keys)
	if err != nil {
		d.logger.WithError(err).Error("failed to list keys")
		return nil, err
	}

	return keys, nil
}

func (d *Keys) GetAllDeleted(_ context.Context) ([]*entities.Key, error) {
	var keys []*entities.Key

	err := d.db.SelectDeleted(&keys)
	if err != nil {
		d.logger.WithError(err).Error("failed to get key")
		return nil, err
	}

	return keys, nil
}

func (d *Keys) Add(_ context.Context, key *entities.Key) error {
	err := d.db.Insert(key)
	if err != nil {
		d.logger.WithError(err).Error("failed to insert key")
		return err
	}

	return nil
}

func (d *Keys) Update(_ context.Context, key *entities.Key) error {
	err := d.db.UpdatePK(key)
	if err != nil {
		d.logger.WithError(err).Error("failed to update key")
		return err
	}

	return nil
}

func (d *Keys) Remove(_ context.Context, id string) error {
	key := &entities.Key{ID: id}
	err := d.db.UpdatePK(key)
	if err != nil {
		d.logger.WithError(err).Error("failed to update key")
		return err
	}

	return nil
}

func (d *Keys) Restore(ctx context.Context, key *entities.Key) error {
	key.Metadata.DeletedAt = time.Time{}
	return d.Update(ctx, key)
}

func (d *Keys) Purge(_ context.Context, id string) error {
	key := &entities.Key{ID: id}
	err := d.db.ForceDeletePK(key)
	if err != nil {
		d.logger.WithError(err).Error("failed to update key")
		return err
	}

	return nil
}
