package database

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Registry interface {
	// Insert inserts a new alias registry
	Insert(ctx context.Context, registry *entities.AliasRegistry) (*entities.AliasRegistry, error)
	// FindOne gets an alias registry
	FindOne(ctx context.Context, name, tenant string) (*entities.AliasRegistry, error)
	// Delete deletes an alias registry
	Delete(ctx context.Context, name, tenant string) error
}

type Alias interface {
	// Insert inserts an alias in the registry
	Insert(ctx context.Context, alias *entities.Alias) (*entities.Alias, error)
	// FindOne gets an alias from the registry
	FindOne(ctx context.Context, registry, key, tenant string) (*entities.Alias, error)
	// Update updates an alias in the registry
	Update(ctx context.Context, alias *entities.Alias) (*entities.Alias, error)
	// Delete deletes an alias from the registry
	Delete(ctx context.Context, registry, key string) error
}
