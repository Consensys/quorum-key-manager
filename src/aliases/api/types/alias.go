package types

import aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"

type Alias struct {
	Key   AliasKey   `json:"key"`
	Value AliasValue `json:"value"`

	registryName RegistryName
}

func FromEntityAlias(ent aliasent.Alias) Alias {
	return Alias{
		registryName: RegistryName(ent.RegistryName),
		Key:          AliasKey(ent.Key),
		Value:        AliasValue(ent.Value),
	}
}

func ToEntityAlias(registry RegistryName, alias Alias) aliasent.Alias {
	return aliasent.Alias{
		RegistryName: aliasent.RegistryName(registry),
		Key:          aliasent.AliasKey(alias.Key),
		Value:        aliasent.AliasValue(alias.Value),
	}
}

func FromEntityAliases(ents []aliasent.Alias) []Alias {
	var als []Alias
	for _, v := range ents {
		als = append(als, FromEntityAlias(v))
	}

	return als
}

type AliasValue string

type AliasKey string

type RegistryName string

type CreateAliasRequest struct {
	Alias
}

type CreateAliasResponse struct {
	Alias
}

type GetAliasResponse struct {
	Alias
}

// TODO the: declare all remaining types

type UpdateAliasRequest struct {
	Value AliasValue `json:"value"`
}

type UpdateAliasResponse struct {
	Value AliasValue `json:"value"`
}
