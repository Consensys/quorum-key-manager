package aliasstore

import (
	"context"

	aliasmodels "github.com/consensys/quorum-key-manager/src/aliases/models"
)

//go:generate mockgen -source=store.go -destination=mock/store.go -package=mock

type Database interface {
	Alias() Alias
}

// Alias handles the alias storing.
type Alias interface {
	// CreateAlias creates an alias in the registry.
	CreateAlias(ctx context.Context, registry aliasmodels.RegistryName, alias aliasmodels.Alias) error
	// GetAlias gets an alias from the registry.
	GetAlias(ctx context.Context, registry aliasmodels.RegistryName, aliasKey aliasmodels.AliasKey) (*aliasmodels.Alias, error)
	// UpdateAlias updates an alias in the registry.
	UpdateAlias(ctx context.Context, registry aliasmodels.RegistryName, alias aliasmodels.Alias) error
	// GetAlias deletes an alias from the registry.
	DeleteAlias(ctx context.Context, registry aliasmodels.RegistryName, aliasKey aliasmodels.AliasKey) error

	// ListAlias lists all aliases from a registry.
	ListAliases(ctx context.Context, registry aliasmodels.RegistryName) ([]aliasmodels.Alias, error)

	// DeleteRegistry deletes a registry, with all the aliases it contained.
	DeleteRegistry(ctx context.Context, registry aliasmodels.RegistryName) error
}
