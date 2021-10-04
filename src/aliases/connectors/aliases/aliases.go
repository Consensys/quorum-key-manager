package aliases

import (
	"context"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

var _ aliasent.AliasBackend = &Connector{}

type Connector struct {
	db aliasstore.Database

	logger log.Logger
}

func NewConnector(db aliasstore.Database, logger log.Logger) *Connector {
	return &Connector{
		db:     db,
		logger: logger,
	}
}

func (m *Connector) CreateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	logger := m.logger.With(
		"registry_name", registry,
		"alias_key", alias.Key,
	)
	a, err := m.db.Alias().CreateAlias(ctx, registry, alias)
	if err != nil {
		return nil, err
	}
	logger.Info("alias created successfully")
	return a, nil
}

func (m *Connector) GetAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) (*aliasent.Alias, error) {
	return m.db.Alias().GetAlias(ctx, registry, aliasKey)
}

func (m *Connector) UpdateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error) {
	logger := m.logger.With(
		"registry_name", registry,
		"alias_key", alias.Key,
	)
	a, err := m.db.Alias().UpdateAlias(ctx, registry, alias)
	if err != nil {
		return nil, err
	}
	logger.Info("alias updated successfully")
	return a, nil
}

func (m *Connector) DeleteAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) error {
	logger := m.logger.With(
		"registry_name", registry,
		"alias_key", aliasKey,
	)
	err := m.db.Alias().DeleteAlias(ctx, registry, aliasKey)
	if err != nil {
		return err
	}
	logger.Info("alias deleted successfully")
	return nil
}

func (m *Connector) ListAliases(ctx context.Context, registry aliasent.RegistryName) ([]aliasent.Alias, error) {
	return m.db.Alias().ListAliases(ctx, registry)
}

func (m *Connector) DeleteRegistry(ctx context.Context, registry aliasent.RegistryName) error {
	return m.db.Alias().DeleteRegistry(ctx, registry)
}
