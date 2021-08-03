package aliasmanager

import (
	"context"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	aliasmodels "github.com/consensys/quorum-key-manager/src/aliases/store/models"
)

type BaseManager struct {
	db aliasstore.Database
}

func New(db aliasstore.Database) *BaseManager {
	return &BaseManager{
		db: db,
	}
}

// Start does nothing as the DB client is already connected.
func (m *BaseManager) Start(_ context.Context) error { return nil }

// Stop does nothing as the DB client should be stopped outside of it.
func (m *BaseManager) Stop(_ context.Context) error { return nil }
func (m *BaseManager) Error() error                 { return nil }
func (m *BaseManager) Close() error                 { return nil }

func (m *BaseManager) CreateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	reg, a := aliasmodels.RegistryName(registry), aliasmodels.AliasFromEntity(alias)
	ent, err := m.db.Alias().CreateAlias(ctx, reg, a)
	return ent.ToEntity(), err
}

func (m *BaseManager) GetAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) (*aliasent.Alias, error) {
	reg, k := aliasmodels.RegistryName(registry), aliasmodels.AliasKey(aliasKey)
	ent, err := m.db.Alias().GetAlias(ctx, reg, k)
	return ent.ToEntity(), err
}

func (m *BaseManager) UpdateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	reg, a := aliasmodels.RegistryName(registry), aliasmodels.AliasFromEntity(alias)
	ent, err := m.db.Alias().UpdateAlias(ctx, reg, a)
	return ent.ToEntity(), err
}
func (m *BaseManager) DeleteAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) error {
	reg, k := aliasmodels.RegistryName(registry), aliasmodels.AliasKey(aliasKey)
	return m.db.Alias().DeleteAlias(ctx, reg, k)
}

func (m *BaseManager) ListAliases(ctx context.Context, registry aliasent.RegistryName) ([]aliasent.Alias, error) {
	reg := aliasmodels.RegistryName(registry)
	als, err := m.db.Alias().ListAliases(ctx, reg)
	if err != nil {
		return nil, err
	}
	ents := aliasmodels.AliasesToEntity(als)
	return ents, nil
}

func (m *BaseManager) DeleteRegistry(ctx context.Context, registry aliasent.RegistryName) error {
	reg := aliasmodels.RegistryName(registry)
	return m.db.Alias().DeleteRegistry(ctx, reg)
}
