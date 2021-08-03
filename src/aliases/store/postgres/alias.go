package aliaspg

import (
	"context"
	goerrors "errors"

	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	aliasmodels "github.com/consensys/quorum-key-manager/src/aliases/store/models"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

var _ aliasstore.Alias = &Alias{}

// Alias stores the alias data in a postgres DB.
type Alias struct {
	pgClient postgres.Client
}

func NewAlias(pgClient postgres.Client) *Alias {
	return &Alias{
		pgClient: pgClient,
	}
}

func (s *Alias) CreateAlias(ctx context.Context, registryName aliasmodels.RegistryName, alias aliasmodels.Alias) (*aliasmodels.Alias, error) {
	alias.RegistryName = registryName
	err := s.pgClient.Insert(ctx, &alias)
	if err != nil {
		return nil, err
	}
	return &alias, nil
}

func (s *Alias) GetAlias(ctx context.Context, registryName aliasmodels.RegistryName, aliasKey aliasmodels.AliasKey) (*aliasmodels.Alias, error) {
	alias := aliasmodels.Alias{Key: aliasKey, RegistryName: registryName}
	err := s.pgClient.SelectPK(ctx, &alias)
	return &alias, err
}

func (s *Alias) UpdateAlias(ctx context.Context, registryName aliasmodels.RegistryName, alias aliasmodels.Alias) (*aliasmodels.Alias, error) {
	alias.RegistryName = registryName
	err := s.pgClient.UpdatePK(ctx, &alias)
	return &alias, err
}

func (s *Alias) DeleteAlias(ctx context.Context, registryName aliasmodels.RegistryName, aliasKey aliasmodels.AliasKey) error {
	alias := aliasmodels.Alias{Key: aliasKey, RegistryName: registryName}
	return s.pgClient.DeletePK(ctx, &alias)
}

func (s *Alias) ListAliases(ctx context.Context, registry aliasmodels.RegistryName) ([]aliasmodels.Alias, error) {
	als := []aliasmodels.Alias{}
	err := s.pgClient.SelectWhere(ctx, &als, "alias.registry_name = ?", registry)
	return als, err
}

func (s *Alias) DeleteRegistry(ctx context.Context, registryName aliasmodels.RegistryName) error {
	return goerrors.New("not implemented")
}
