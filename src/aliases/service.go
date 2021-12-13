package aliases

import (
	"context"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock

// Registries handles the aliases registries.
type Registries interface {
	// Create creates an alias registry
	Create(ctx context.Context, name string, allowedTenants []string, userInfo *auth.UserInfo) (*entities.AliasRegistry, error)
	// Get gets an alias registry
	Get(ctx context.Context, name string, userInfo *auth.UserInfo) (*entities.AliasRegistry, error)
	// Delete deletes an alias registry, with all the aliases it contains
	Delete(ctx context.Context, name string, userInfo *auth.UserInfo) error
}

// Aliases handles the aliases.
type Aliases interface {
	// Create creates an alias in the registry
	Create(ctx context.Context, registry, key, kind string, value interface{}, userInfo *auth.UserInfo) (*entities.Alias, error)
	// Get gets an alias from the registry
	Get(ctx context.Context, registry string, key string, userInfo *auth.UserInfo) (*entities.Alias, error)
	// Update updates an alias in the registry
	Update(ctx context.Context, registry, key, kind string, value interface{}, userInfo *auth.UserInfo) (*entities.Alias, error)
	// Delete deletes an alias from the registry
	Delete(ctx context.Context, registry string, key string, userInfo *auth.UserInfo) error
	// Parse parses an alias string and returns the registryName and the aliasKey
	Parse(alias string) (regName string, key string, isAlias bool)
	// Replace replaces a slice of potential aliases with a slice having all the aliases replaced by their value
	Replace(ctx context.Context, addrs []string, userInfo *auth.UserInfo) ([]string, error)
	// ReplaceSimple replaces a potential alias with its first and only value
	ReplaceSimple(ctx context.Context, addr string, userInfo *auth.UserInfo) (string, error)
}
