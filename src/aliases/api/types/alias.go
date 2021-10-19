package types

import "github.com/consensys/quorum-key-manager/src/aliases/entities"

type Alias struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`

	registryName string
}

// FormatEntityAlias format an alias entity to an alias API type.
func FormatEntityAlias(ent entities.Alias) Alias {
	return Alias{
		registryName: ent.RegistryName,
		Key:          ent.Key,
		Value:        ent.Value,
	}
}

// FormatAlias format an alias API type to an alias entity.
func FormatAlias(registry, key string, value []string) entities.Alias {
	return entities.Alias{
		RegistryName: registry,
		Key:          key,
		Value:        value,
	}
}

// FormatEntityAliases formats a slice of alias entities into a slice of alias API type.
func FormatEntityAliases(ents []entities.Alias) []Alias {
	var als = []Alias{}
	for _, v := range ents {
		als = append(als, FormatEntityAlias(v))
	}

	return als
}

// AliasRequest creates or modifies an alias value.
type AliasRequest struct {
	Value []string `json:"value" validate:"min=1,unique,dive,base64,required" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}

// AliasResponse returns the alias value.
type AliasResponse struct {
	Value []string `json:"value" validate:"min=1,unique,dive,base64,required" example:"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=,B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="`
}
