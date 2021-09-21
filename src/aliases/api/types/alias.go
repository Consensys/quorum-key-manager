package types

import aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"

type Alias struct {
	Key   AliasKey   `json:"key"`
	Value AliasValue `json:"value"`

	registryName RegistryName
}

func FormatEntityAlias(ent aliasent.Alias) Alias {
	return Alias{
		registryName: RegistryName(ent.RegistryName),
		Key:          AliasKey(ent.Key),
		Value:        AliasValue(ent.Value),
	}
}

func FormatAlias(registry RegistryName, key string, value AliasValue) aliasent.Alias {
	return aliasent.Alias{
		RegistryName: aliasent.RegistryName(registry),
		Key:          aliasent.AliasKey(key),
		Value:        aliasent.AliasValue(value),
	}
}

func FormatEntityAliases(ents []aliasent.Alias) []Alias {
	var als = []Alias{}
	for _, v := range ents {
		als = append(als, FormatEntityAlias(v))
	}

	return als
}

type AliasValue string

type AliasKey string

type RegistryName string

type AliasRequest struct {
	Value AliasValue `json:"value"`
}

type AliasResponse struct {
	Value AliasValue `json:"value"`
}
