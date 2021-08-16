package aliaspg

import (
	"context"
	goerrors "errors"

	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	aliasmodels "github.com/consensys/quorum-key-manager/src/aliases/store/models"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

var _ aliases.Alias = &AliasStore{}

// AliasStore stores the alias data in a postgres DB.
type AliasStore struct {
	pgClient postgres.Client
}

func NewAlias(pgClient postgres.Client) *AliasStore {
	return &AliasStore{
		pgClient: pgClient,
	}
}

func (s *AliasStore) CreateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	a := aliasmodels.AliasFromEntity(alias)
	a.RegistryName = aliasmodels.RegistryName(registry)

	err := s.pgClient.Insert(ctx, &a)
	if err != nil {
		return nil, err
	}
	return &alias, nil
}

func (s *AliasStore) GetAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) (*aliasent.Alias, error) {
	a := aliasmodels.Alias{
		Key:          aliasmodels.AliasKey(aliasKey),
		RegistryName: aliasmodels.RegistryName(registry),
	}
	err := s.pgClient.SelectPK(ctx, &a)
	return a.ToEntity(), err
}

func (s *AliasStore) UpdateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	a := aliasmodels.AliasFromEntity(alias)
	a.RegistryName = aliasmodels.RegistryName(registry)

	err := s.pgClient.UpdatePK(ctx, &a)
	return a.ToEntity(), err
}

func (s *AliasStore) DeleteAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) error {
	a := aliasmodels.Alias{
		Key:          aliasmodels.AliasKey(aliasKey),
		RegistryName: aliasmodels.RegistryName(registry),
	}
	return s.pgClient.DeletePK(ctx, &a)
}

func (s *AliasStore) ListAliases(ctx context.Context, registry aliasent.RegistryName) ([]aliasent.Alias, error) {
	reg := aliasmodels.RegistryName(registry)

	als := []aliasmodels.Alias{}
	err := s.pgClient.SelectWhere(ctx, &als, "alias.registry_name = ?", reg)
	if err != nil {
		return nil, err
	}

	ents := aliasmodels.AliasesToEntity(als)
	return ents, err
}

func (s *AliasStore) DeleteRegistry(ctx context.Context, registry aliasent.RegistryName) error {
	return goerrors.New("not implemented")
}
