package aliases

import (
	"context"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

//go:generate mockgen -destination=mock/service.go -package=mock . Service,Interactor,Parser

type Service interface {
	Interactor
	Parser
}

// Interactor handles the aliases storage.
type Interactor interface {
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

// Parser parses and replace aliases.
type Parser interface {
	ParseAlias(alias string) (regName string, aliasKey string, isAlias bool)
	ReplaceAliases(ctx context.Context, addrs []string) ([]string, error)
}
