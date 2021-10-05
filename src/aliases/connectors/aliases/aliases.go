package aliases

import (
	"context"

	aliasdb "github.com/consensys/quorum-key-manager/src/aliases/database"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

// We make sure Connector implements aliasent.AliasBackend
var _ aliasent.AliasBackend = &Connector{}

// Connector is the service layer for other service to query.
type Connector struct {
	db aliasdb.Database

	logger log.Logger
}

func NewConnector(db aliasdb.Database, logger log.Logger) *Connector {
	return &Connector{
		db:     db,
		logger: logger,
	}
}

func (m *Connector) CreateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error) {
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

func (m *Connector) GetAlias(ctx context.Context, registry, aliasKey string) (*aliasent.Alias, error) {
	return m.db.Alias().GetAlias(ctx, registry, aliasKey)
}

func (m *Connector) UpdateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error) {
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

func (m *Connector) DeleteAlias(ctx context.Context, registry, aliasKey string) error {
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

func (m *Connector) ListAliases(ctx context.Context, registry string) ([]aliasent.Alias, error) {
	return m.db.Alias().ListAliases(ctx, registry)
}

func (m *Connector) DeleteRegistry(ctx context.Context, registry string) error {
	logger := m.logger.With(
		"registry_name", registry,
	)
	err := m.db.Alias().DeleteRegistry(ctx, registry)
	if err != nil {
		return err
	}
	logger.Info("registry deleted successfully")
	return nil
}
