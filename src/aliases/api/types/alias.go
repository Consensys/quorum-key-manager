package types

import (
	"time"

	"github.com/consensys/quorum-key-manager/src/entities"
)

// AliasRequest creates or modifies an alias value.
type AliasRequest struct {
	Kind  string      `json:"type" validate:"required,isAliasKind" example:"string"`
	Value interface{} `json:"value" validate:"required" example:"my-alias" swaggertype:"string"`
}

// AliasResponse returns the alias value.
type AliasResponse struct {
	Key       string      `json:"key" example:"my-alias"`
	Kind      string      `json:"type" example:"string"`
	Value     interface{} `json:"value" example:"my-alias-value" swaggertype:"string"`
	Registry  string      `json:"registry" example:"my-registry"`
	CreatedAt time.Time   `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt time.Time   `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
}

func NewAliasResponse(alias *entities.Alias) *AliasResponse {
	return &AliasResponse{
		Key:       alias.Key,
		Kind:      alias.Kind,
		Value:     alias.Value,
		Registry:  alias.RegistryName,
		CreatedAt: alias.CreatedAt,
		UpdatedAt: alias.UpdatedAt,
	}
}
