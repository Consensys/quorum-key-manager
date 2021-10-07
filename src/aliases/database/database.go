package database

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	Alias() aliases.Interactor
}

type AliasRepository interface {
	// CreateAlias creates an alias in the registry.
	CreateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error)
	// GetAlias gets an alias from the registry.
	GetAlias(ctx context.Context, registry string, aliasKey string) (*aliasent.Alias, error)
	// UpdateAlias updates an alias in the registry.
	UpdateAlias(ctx context.Context, registry string, alias aliasent.Alias) (*aliasent.Alias, error)
	// GetAlias deletes an alias from the registry.
	DeleteAlias(ctx context.Context, registry string, aliasKey string) error

	// ListAlias lists all aliases from a registry.
	ListAliases(ctx context.Context, registry string) ([]aliasent.Alias, error)

	// DeleteRegistry deletes a registry, with all the aliases it contained.
	DeleteRegistry(ctx context.Context, registry string) error
}
