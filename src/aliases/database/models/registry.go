package models

import (
	"github.com/consensys/quorum-key-manager/src/entities"
	"time"
)

type Registry struct {
	tableName struct{} `pg:"registries"` // nolint:unused,structcheck // reason

	Name           string  `pg:",pk"`
	Aliases        []Alias `pg:"rel:has-many"`
	AllowedTenants []string
	CreatedAt      time.Time `pg:"default:now()"`
	UpdatedAt      time.Time `pg:"default:now()"`
}

func NewRegistry(registry *entities.AliasRegistry) *Registry {
	return &Registry{
		Name:           registry.Name,
		AllowedTenants: registry.AllowedTenants,
		CreatedAt:      registry.CreatedAt,
		UpdatedAt:      registry.UpdatedAt,
	}
}

func (r *Registry) ToEntity() *entities.AliasRegistry {
	var aliases []entities.Alias
	for _, aliasModel := range r.Aliases {
		aliases = append(aliases, *aliasModel.ToEntity())
	}

	return &entities.AliasRegistry{
		Name:           r.Name,
		AllowedTenants: r.AllowedTenants,
		Aliases:        aliases,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}
