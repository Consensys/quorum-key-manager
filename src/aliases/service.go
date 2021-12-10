package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock

// Registries handles the aliases registries.
type Registries interface {
	// Create creates an alias registry
	Create(ctx context.Context, name string) (*entities.AliasRegistry, error)
	// Get gets an alias registry
	Get(ctx context.Context, name string) (*entities.AliasRegistry, error)
	// Delete deletes an alias registry, with all the aliases it contains
	Delete(ctx context.Context, name string) error
}

// Aliases handles the aliases.
type Aliases interface {
	// Create creates an alias in the registry
	Create(ctx context.Context, registry, key, kind string, value interface{}) (*entities.Alias, error)
	// Get gets an alias from the registry
	Get(ctx context.Context, registry string, key string) (*entities.Alias, error)
	// Update updates an alias in the registry
	Update(ctx context.Context, registry, key, kind string, value interface{}) (*entities.Alias, error)
	// Delete deletes an alias from the registry
	Delete(ctx context.Context, registry string, key string) error
	// Parse parses an alias string and returns the registryName and the aliasKey
	Parse(alias string) (regName string, key string, isAlias bool)
	// Replace replaces a slice of potential aliases with a slice having all the aliases replaced by their value
	Replace(ctx context.Context, addrs []string) ([]string, error)
	// ReplaceSimple replaces a potential alias with its first and only value
	ReplaceSimple(ctx context.Context, addr string) (string, error)
}
