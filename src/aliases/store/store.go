package aliasstore

import (
	"context"
	goerrors "errors"

	aliases "github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

var _ aliases.API = &Store{}

type Store struct {
	pgClient postgres.Client
}

func New(pgClient postgres.Client) *Store {
	return &Store{
		pgClient: pgClient,
	}
}

func (s *Store) CreateAlias(ctx context.Context, alias aliases.Alias) error {
	err := s.pgClient.Insert(ctx, &alias)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAlias(ctx context.Context, registry aliases.RegistryID, alias aliases.AliasID) (*aliases.Alias, error) {
	a := aliases.Alias{ID: alias, RegistryID: registry}
	err := s.pgClient.Select(ctx, &a)
	return &a, err
}

func (s *Store) UpdateAlias(ctx context.Context, alias aliases.Alias) error {
	return s.pgClient.UpdatePK(ctx, &alias)
}

func (s *Store) DeleteAlias(ctx context.Context, registry aliases.RegistryID, alias aliases.AliasID) error {
	a := aliases.Alias{ID: alias, RegistryID: registry}
	return s.pgClient.DeletePK(ctx, &a)
}

func (s *Store) ListAliases(ctx context.Context, registry aliases.RegistryID) ([]aliases.Alias, error) {
	als := []aliases.Alias{}
	err := s.pgClient.SelectMany(ctx, &aliases.Alias{}, &als, "alias.registry_id = ?", registry)
	return als, err
}

func (s *Store) DeleteRegistry(ctx context.Context, registry aliases.RegistryID) error {
	return goerrors.New("not implemented")
}
