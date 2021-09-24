package aliaspg

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	aliasmodels "github.com/consensys/quorum-key-manager/src/aliases/store/models"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

var _ aliasent.AliasBackend = &AliasStore{}

// AliasStore stores the alias data in a postgres DB.
type AliasStore struct {
	pgClient postgres.Client
	logger   log.Logger
}

func NewAlias(pgClient postgres.Client, logger log.Logger) *AliasStore {
	return &AliasStore{
		pgClient: pgClient,
		logger:   logger,
	}
}

func (s *AliasStore) CreateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	a := aliasmodels.AliasFromEntity(alias)
	a.RegistryName = aliasmodels.RegistryName(registry)

	err = s.pgClient.Insert(ctx, &a)
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
	if err != nil {
		return nil, err
	}
	return a.ToEntity(), nil
}

func (s *AliasStore) UpdateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	a := aliasmodels.AliasFromEntity(alias)
	a.RegistryName = aliasmodels.RegistryName(registry)

	err := s.pgClient.UpdatePK(ctx, &a)
	if err != nil {
		return nil, err
	}
	return a.ToEntity(), nil
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

	var als []aliasmodels.Alias
	err := s.pgClient.SelectWhere(ctx, &als, "alias.registry_name = ?", reg)
	if err != nil {
		return nil, err
	}

	return aliasmodels.AliasesToEntity(als), nil
}

func (s *AliasStore) DeleteRegistry(ctx context.Context, registry aliasent.RegistryName) error {
	return errors.NotImplementedError("DeleteRegistry not implemented")
}
