package aliasstore

import (
	"context"
	"errors"

	aliases "github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/go-pg/pg/v10"
)

var _ aliases.API = &Store{}

type Store struct {
	db *pg.DB
}

func New(db *pg.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateAlias(ctx context.Context, registry aliases.RegistryID, alias aliases.AliasID) error {
	a := aliases.Alias{ID: alias, RegistryID: registry}
	q := s.db.ModelContext(ctx, &a)
	_, err := q.Insert()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAlias(ctx context.Context, registry aliases.RegistryID, alias aliases.AliasID) (*aliases.Alias, error) {
	a := aliases.Alias{ID: alias, RegistryID: registry}
	q := s.db.ModelContext(ctx, &a)
	ret := aliases.Alias{}
	err := q.WherePK().Select(&ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (s *Store) UpdateAlias(ctx context.Context, registry aliases.RegistryID, alias aliases.AliasID) error {
	return errors.New("not implemented")
}

func (s *Store) DeleteAlias(ctx context.Context, registry aliases.RegistryID, alias aliases.AliasID) error {
	return errors.New("not implemented")
}

func (s *Store) GetRegistry(ctx context.Context, registry aliases.RegistryID) (*aliases.Registry, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) DeleteRegistry(ctx context.Context, registry aliases.RegistryID) error {
	return errors.New("not implemented")
}
