package aliasstore

import (
	"context"
	goerrors "errors"

	aliases "github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

var _ aliases.Backend = &Store{}

type Store struct {
	pgClient postgres.Client
}

func New(pgClient postgres.Client) *Store {
	return &Store{
		pgClient: pgClient,
	}
}

func (s *Store) CreateAlias(ctx context.Context, registryName aliases.RegistryName, alias aliases.Alias) error {
	alias.RegistryName = registryName
	err := s.pgClient.Insert(ctx, &alias)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAlias(ctx context.Context, registryName aliases.RegistryName, aliasKey aliases.AliasKey) (*aliases.Alias, error) {
	a := aliases.Alias{Key: aliasKey, RegistryName: registryName}
	err := s.pgClient.SelectPK(ctx, &a)
	return &a, err
}

func (s *Store) UpdateAlias(ctx context.Context, registryName aliases.RegistryName, alias aliases.Alias) error {
	alias.RegistryName = registryName
	return s.pgClient.UpdatePK(ctx, &alias)
}

func (s *Store) DeleteAlias(ctx context.Context, registryName aliases.RegistryName, aliasKey aliases.AliasKey) error {
	a := aliases.Alias{Key: aliasKey, RegistryName: registryName}
	return s.pgClient.DeletePK(ctx, &a)
}

func (s *Store) ListAliases(ctx context.Context, registry aliases.RegistryName) ([]aliases.Alias, error) {
	als := []aliases.Alias{}
	err := s.pgClient.SelectMany(ctx, &aliases.Alias{}, &als, "alias.registry_name = ?", registry)
	return als, err
}

func (s *Store) DeleteRegistry(ctx context.Context, registryName aliases.RegistryName) error {
	return goerrors.New("not implemented")
}
