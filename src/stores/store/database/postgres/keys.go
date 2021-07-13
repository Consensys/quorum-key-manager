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

func (d *Keys) Get(ctx context.Context, id string) (*entities.Key, error) {
	key := &entities.Key{ID: id}

	err := d.db.SelectPK(ctx, key)
	if err != nil {
		d.logger.WithError(err).Error("failed to get key")
		return nil, err
	}

	return key, nil
}

func (d *Keys) GetDeleted(ctx context.Context, id string) (*entities.Key, error) {
	key := &entities.Key{ID: id}

	err := d.db.SelectDeletedPK(ctx, key)
	if err != nil {
		d.logger.WithError(err).Error("failed to get key")
		return nil, err
	}

	return key, nil
}

func (d *Keys) GetAll(ctx context.Context) ([]*entities.Key, error) {
	var keys []*entities.Key

	err := d.db.Select(ctx, &keys)
	if err != nil {
		d.logger.WithError(err).Error("failed to list keys")
		return nil, err
	}

	return keys, nil
}

func (d *Keys) GetAllDeleted(ctx context.Context) ([]*entities.Key, error) {
	var keys []*entities.Key

	err := d.db.SelectDeleted(ctx, keys)
	if err != nil {
		d.logger.WithError(err).Error("failed to get key")
		return nil, err
	}

	return keys, nil
}

func (d *Keys) Add(ctx context.Context, key *entities.Key) error {
	err := d.db.Insert(ctx, key)
	if err != nil {
		d.logger.WithError(err).Error("failed to insert key")
		return err
	}

	return nil
}

func (d *Keys) Update(ctx context.Context, key *entities.Key) error {
	err := d.db.UpdatePK(ctx, key)
	if err != nil {
		d.logger.WithError(err).Error("failed to update key")
		return err
	}

	return nil
}

func (d *Keys) Remove(ctx context.Context, id string) error {
	key := &entities.Key{ID: id}
	err := d.db.UpdatePK(ctx, key)
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

func (d *Keys) Purge(ctx context.Context, id string) error {
	key := &entities.Key{ID: id}
	err := d.db.ForceDeletePK(ctx, key)
	if err != nil {
		d.logger.WithError(err).Error("failed to update key")
		return err
	}

	return nil
}
