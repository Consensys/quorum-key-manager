package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/consensys/quorum-key-manager/src/aliases/database"
	"github.com/consensys/quorum-key-manager/src/aliases/database/models"
	"github.com/consensys/quorum-key-manager/src/entities"
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
	aliasModel := &models.Alias{Key: key, RegistryName: registry}

	query := "key = ?"
	if tenant != "" {
		query = fmt.Sprintf("%s AND '%s' = ANY(registry.allowed_tenants)", query, tenant)
	}

	err := r.pgClient.SelectWhere(ctx, aliasModel, query, []string{"Registry"}, key)
	if err != nil {
		return nil, err
	}

	return aliasModel.ToEntity(), nil
}

func (r *Alias) Update(ctx context.Context, alias *entities.Alias) (*entities.Alias, error) {
	aliasModel := models.NewAlias(alias)
	aliasModel.UpdatedAt = time.Now()

	err := r.pgClient.UpdatePK(ctx, aliasModel)
	if err != nil {
		return nil, err
	}

	// Update does not update the model, we must update and then get
	return r.FindOne(ctx, alias.RegistryName, alias.Key, "")
}

func (r *Alias) Delete(ctx context.Context, registry, key string) error {
	err := r.pgClient.DeletePK(ctx, &models.Alias{Key: key, RegistryName: registry})
	if err != nil {
		return err
	}

	return nil
}
