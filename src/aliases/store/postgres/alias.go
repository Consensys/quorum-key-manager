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

func (s *AliasStore) CreateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error) {
	logger := s.logger.With(
		"registry_name", registry,
		"alias_key", alias.Key,
	)
	a := aliasmodels.AliasFromEntity(alias)
	a.RegistryName = registry

	err := s.pgClient.Insert(ctx, &a)
	if err != nil {
		msg := "failed to create alias"
		logger.WithError(err).Error(msg)
		return nil, err
	}
	return &alias, nil
}

func (s *AliasStore) GetAlias(ctx context.Context, registry string, aliasKey string) (*aliasent.Alias, error) {
	logger := s.logger.With(
		"registry_name", registry,
		"alias_key", aliasKey,
	)
	a := aliasmodels.Alias{
		Key:          aliasKey,
		RegistryName: registry,
	}
	err := s.pgClient.SelectPK(ctx, &a)
	if err != nil {
		msg := "failed to get alias"
		logger.WithError(err).Error(msg)
		return nil, err
	}
	return a.ToEntity(), nil
}

func (s *AliasStore) UpdateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error) {
	logger := s.logger.With(
		"registry_name", registry,
		"alias_key", alias.Key,
	)
	a := aliasmodels.AliasFromEntity(alias)
	a.RegistryName = registry

	err := s.pgClient.UpdatePK(ctx, &a)
	if err != nil {
		msg := "failed to update alias"
		logger.WithError(err).Error(msg)
		return nil, err
	}
	return a.ToEntity(), nil
}

func (s *AliasStore) DeleteAlias(ctx context.Context, registry string, aliasKey string) error {
	logger := s.logger.With(
		"registry_name", registry,
		"alias_key", aliasKey,
	)
	a := aliasmodels.Alias{
		Key:          aliasKey,
		RegistryName: registry,
	}

	err := s.pgClient.DeletePK(ctx, &a)
	if err != nil {
		msg := "failed to delete alias"
		logger.WithError(err).Error(msg)
		return err
	}
	return nil
}

func (s *AliasStore) ListAliases(ctx context.Context, registry string) ([]aliasent.Alias, error) {
	logger := s.logger.With(
		"registry_name", registry,
	)
	reg := registry

	var als []aliasmodels.Alias
	err := s.pgClient.SelectWhere(ctx, &als, "alias.registry_name = ?", reg)
	if err != nil {
		msg := "failed to list aliases"
		logger.WithError(err).Error(msg)
		return nil, err
	}

	return aliasmodels.AliasesToEntity(als), nil
}

func (s *AliasStore) DeleteRegistry(ctx context.Context, registry string) error {
	return errors.NotImplementedError("DeleteRegistry not implemented")
}
