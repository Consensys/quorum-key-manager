package types

import (
	"github.com/consensys/quorum-key-manager/src/entities"
	"time"
)

type RegistryResponse struct {
	Name      string           `json:"name" example:"my-alias-registry"`
	Aliases   []entities.Alias `json:"aliases"`
	CreatedAt time.Time        `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt time.Time        `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
}

func NewRegistryResponse(registry *entities.AliasRegistry) *RegistryResponse {
	return &RegistryResponse{
		Name:      registry.Name,
		Aliases:   registry.Aliases,
		CreatedAt: registry.CreatedAt,
		UpdatedAt: registry.UpdatedAt,
	}
}
