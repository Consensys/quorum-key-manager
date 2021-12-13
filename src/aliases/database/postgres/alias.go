package postgres

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/entities"

	"github.com/consensys/quorum-key-manager/src/aliases/database"
	"github.com/consensys/quorum-key-manager/src/aliases/database/models"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

type Alias struct {
	pgClient postgres.Client
}

var _ database.Alias = &Alias{}

func NewAlias(pgClient postgres.Client) *Alias {
	return &Alias{pgClient: pgClient}
}

func (r *Alias) Insert(ctx context.Context, alias *entities.Alias) (*entities.Alias, error) {
	aliasModel := models.NewAlias(alias)

	err := r.pgClient.Insert(ctx, aliasModel)
	if err != nil {
		return nil, err
	}

	return aliasModel.ToEntity(), nil
}

func (r *Alias) FindOne(ctx context.Context, registry, key, tenant string) (*entities.Alias, error) {
	aliasModel := &models.Alias{
		Key:          key,
		RegistryName: registry,
	}
	err := r.pgClient.SelectWhere(ctx, aliasModel, r.whereTenant(tenant), []string{"Registry._"}, key)
	if err != nil {
		return nil, err
	}

	return aliasModel.ToEntity(), nil
}

func (r *Alias) Update(ctx context.Context, alias *entities.Alias, tenant string) (*entities.Alias, error) {
	aliasModel := models.NewAlias(alias)

	err := r.pgClient.UpdateWhere(ctx, aliasModel, r.whereTenant(tenant), alias.Key)
	if err != nil {
		return nil, err
	}

	return aliasModel.ToEntity(), nil
}

func (r *Alias) Delete(ctx context.Context, registry, key, tenant string) error {
	aliasModel := &models.Alias{
		Key:          key,
		RegistryName: registry,
	}

	err := r.pgClient.DeleteWhere(ctx, aliasModel, r.whereTenant(tenant), key)
	if err != nil {
		return err
	}

	return nil
}

func (r *Alias) whereTenant(tenant string) string {
	query := "key = ?"
	if tenant != "" {
		return fmt.Sprintf("%s AND '%s' = ANY(registry.allowed_tenants)", query, tenant)
	}

	return query
}
