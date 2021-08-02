package aliaspg

import (
	"context"
	goerrors "errors"

	models "github.com/consensys/quorum-key-manager/src/aliases/models"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

var _ aliasstore.Store = &Store{}

// Store stores the alias data in a postgres DB.
type Store struct {
	pgClient postgres.Client
}

func New(pgClient postgres.Client) *Store {
	return &Store{
		pgClient: pgClient,
	}
}

func (s *Store) CreateAlias(ctx context.Context, registryName models.RegistryName, alias models.Alias) error {
	alias.RegistryName = registryName
	err := s.pgClient.Insert(ctx, &alias)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAlias(ctx context.Context, registryName models.RegistryName, aliasKey models.AliasKey) (*models.Alias, error) {
	a := models.Alias{Key: aliasKey, RegistryName: registryName}
	err := s.pgClient.SelectPK(ctx, &a)
	return &a, err
}

func (s *Store) UpdateAlias(ctx context.Context, registryName models.RegistryName, alias models.Alias) error {
	alias.RegistryName = registryName
	return s.pgClient.UpdatePK(ctx, &alias)
}

func (s *Store) DeleteAlias(ctx context.Context, registryName models.RegistryName, aliasKey models.AliasKey) error {
	a := models.Alias{Key: aliasKey, RegistryName: registryName}
	return s.pgClient.DeletePK(ctx, &a)
}

func (s *Store) ListAliases(ctx context.Context, registry models.RegistryName) ([]models.Alias, error) {
	als := []models.Alias{}
	err := s.pgClient.SelectWhere(ctx, &als, "alias.registry_name = ?", registry)
	return als, err
}

func (s *Store) DeleteRegistry(ctx context.Context, registryName models.RegistryName) error {
	return goerrors.New("not implemented")
}
