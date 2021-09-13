package aliasmanager

import (
	"context"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
)

var _ aliasent.AliasBackend = &BaseManager{}

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
	return m.db.Alias().CreateAlias(ctx, registry, alias)
}

func (m *BaseManager) GetAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) (*aliasent.Alias, error) {
	return m.db.Alias().GetAlias(ctx, registry, aliasKey)
}

func (m *BaseManager) UpdateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	return m.db.Alias().UpdateAlias(ctx, registry, alias)
}

func (m *BaseManager) DeleteAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) error {
	return m.db.Alias().DeleteAlias(ctx, registry, aliasKey)
}

func (m *BaseManager) ListAliases(ctx context.Context, registry aliasent.RegistryName) ([]aliasent.Alias, error) {
	return m.db.Alias().ListAliases(ctx, registry)
}

func (m *BaseManager) DeleteRegistry(ctx context.Context, registry aliasent.RegistryName) error {
	return m.db.Alias().DeleteRegistry(ctx, registry)
}
