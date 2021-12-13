package types

import (
	"time"

	"github.com/consensys/quorum-key-manager/src/entities"
)

type CreateRegistryRequest struct {
	AllowedTenants []string `json:"allowedTenants,omitempty" example:"tenant1,tenant2"`
}

type RegistryResponse struct {
	Name           string          `json:"name" example:"my-alias-registry"`
	Aliases        []AliasResponse `json:"aliases"`
	AllowedTenants []string        `json:"allowedTenants" example:"tenant1,tenant2"`
	CreatedAt      time.Time       `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt      time.Time       `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
}

func NewRegistryResponse(registry *entities.AliasRegistry) *RegistryResponse {
	aliases := []AliasResponse{}
	for _, alias := range registry.Aliases {
		aliases = append(aliases, *NewAliasResponse(&alias))
	}

	return &RegistryResponse{
		Name:           registry.Name,
		Aliases:        aliases,
		AllowedTenants: registry.AllowedTenants,
		CreatedAt:      registry.CreatedAt,
		UpdatedAt:      registry.UpdatedAt,
	}
}
