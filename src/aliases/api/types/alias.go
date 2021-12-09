package types

import (
	"github.com/consensys/quorum-key-manager/src/entities"
)

type Alias struct {
	Key   string              `json:"key"`
	Value entities.AliasValue `json:"value"`

	registryName string
}

// FormatAlias format an alias API type to an alias entity.
func FormatAlias(registry, key string, kind entities.Kind, value interface{}) entities.Alias {
	av := entities.AliasValue{
		Kind:  kind,
		Value: value,
	}
	return entities.Alias{
		RegistryName: registry,
		Key:          key,
		Value:        av,
	}
}

// FormatEntityAliases formats a slice of alias entities into a slice of alias API type.
func FormatEntityAliases(ents []entities.Alias) []Alias {
	var als = []Alias{}
	for _, v := range ents {
		als = append(als, Alias{
			registryName: v.RegistryName,
			Key:          v.Key,
			Value:        v.Value,
		})
	}

	return als
}

// AliasRequest creates or modifies an alias value.
type AliasRequest struct {
	Kind  entities.Kind `json:"type" validate:"required" example:"string"`
	Value interface{}   `json:"value" validate:"required" example:"a2V5MQo=" swaggertype:"string"`
}

// AliasResponse returns the alias value.
type AliasResponse struct {
	Kind  entities.Kind `json:"type"`
	Value interface{}   `json:"value"`
}
