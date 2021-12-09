package database

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	Aliases() Aliases
}

type Aliases interface {
	// Create creates an alias in the registry.
	Create(ctx context.Context, registry string, alias *entities.Alias) (*entities.Alias, error)
	// Get gets an alias from the registry.
	Get(ctx context.Context, registry string, aliasKey string) (*entities.Alias, error)
	// Update updates an alias in the registry.
	Update(ctx context.Context, registry string, alias *entities.Alias) (*entities.Alias, error)
	// Delete deletes an alias from the registry.
	Delete(ctx context.Context, registry string, aliasKey string) error
	// List lists all aliases from a registry.
	List(ctx context.Context, registry string) ([]entities.Alias, error)
	// DeleteRegistry deletes a registry, with all the aliases it contained.
	DeleteRegistry(ctx context.Context, registry string) error
}
