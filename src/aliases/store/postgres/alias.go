package aliaspg

import (
	"context"
	goerrors "errors"

	models "github.com/consensys/quorum-key-manager/src/aliases/models"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
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

func (s *Alias) CreateAlias(ctx context.Context, registryName models.RegistryName, alias models.Alias) error {
	alias.RegistryName = registryName
	err := s.pgClient.Insert(ctx, &alias)
	if err != nil {
		return err
	}
	return nil
}

func (s *Alias) GetAlias(ctx context.Context, registryName models.RegistryName, aliasKey models.AliasKey) (*models.Alias, error) {
	a := models.Alias{Key: aliasKey, RegistryName: registryName}
	err := s.pgClient.SelectPK(ctx, &a)
	return &a, err
}

func (s *Alias) UpdateAlias(ctx context.Context, registryName models.RegistryName, alias models.Alias) error {
	alias.RegistryName = registryName
	return s.pgClient.UpdatePK(ctx, &alias)
}

func (s *Alias) DeleteAlias(ctx context.Context, registryName models.RegistryName, aliasKey models.AliasKey) error {
	a := models.Alias{Key: aliasKey, RegistryName: registryName}
	return s.pgClient.DeletePK(ctx, &a)
}

func (s *Alias) ListAliases(ctx context.Context, registry models.RegistryName) ([]models.Alias, error) {
	als := []models.Alias{}
	err := s.pgClient.SelectWhere(ctx, &als, "alias.registry_name = ?", registry)
	return als, err
}

func (s *Alias) DeleteRegistry(ctx context.Context, registryName models.RegistryName) error {
	return goerrors.New("not implemented")
}
