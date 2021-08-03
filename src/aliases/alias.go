package aliases

import (
	"context"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

//go:generate mockgen -source=alias.go -destination=mock/alias.go -package=mock

// Alias handles the aliases.
type Alias interface {
	// CreateAlias creates an alias in the registry.
	CreateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error)
	// GetAlias gets an alias from the registry.
	GetAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) (*aliasent.Alias, error)
	// UpdateAlias updates an alias in the registry.
	UpdateAlias(ctx context.Context, registry aliasent.RegistryName, alias aliasent.Alias) (*aliasent.Alias, error)
	// GetAlias deletes an alias from the registry.
	DeleteAlias(ctx context.Context, registry aliasent.RegistryName, aliasKey aliasent.AliasKey) error

	// ListAlias lists all aliases from a registry.
	ListAliases(ctx context.Context, registry aliasent.RegistryName) ([]aliasent.Alias, error)

	// DeleteRegistry deletes a registry, with all the aliases it contained.
	DeleteRegistry(ctx context.Context, registry aliasent.RegistryName) error
}
