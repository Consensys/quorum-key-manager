package postgres

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/entities"

	"github.com/consensys/quorum-key-manager/src/aliases/database"
	"github.com/consensys/quorum-key-manager/src/aliases/database/models"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

type Registry struct {
	pgClient postgres.Client
}

var _ database.Registry = &Registry{}

func NewRegistry(pgClient postgres.Client) *Registry {
	return &Registry{pgClient: pgClient}
}

func (r *Registry) Insert(ctx context.Context, registry *entities.AliasRegistry) (*entities.AliasRegistry, error) {
	registryModel := models.NewRegistry(registry)

	err := r.pgClient.Insert(ctx, registryModel)
	if err != nil {
		return nil, err
	}

	return registryModel.ToEntity(), nil
}

func (r *Registry) FindOne(ctx context.Context, name, tenant string) (*entities.AliasRegistry, error) {
	registryModel := &models.Registry{Name: name}

	err := r.pgClient.SelectWhere(ctx, registryModel, r.whereTenant(tenant), name)
	if err != nil {
		return nil, err
	}

	return registryModel.ToEntity(), nil
}

func (r *Registry) Delete(ctx context.Context, name, tenant string) error {
	err := r.pgClient.DeleteWhere(ctx, &models.Registry{Name: name}, r.whereTenant(tenant), name)
	if err != nil {
		return err
	}

	return nil
}

func (r *Registry) whereTenant(tenant string) string {
	query := "name = ?"
	if tenant != "" {
		return fmt.Sprintf("%s AND '%s' = ANY(allowed_tenants)", query, tenant)
	}

	return query
}
