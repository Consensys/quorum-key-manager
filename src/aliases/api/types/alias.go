package types

import aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"

type Alias struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`

	registryName string
}

func FormatEntityAlias(ent aliasent.Alias) Alias {
	return Alias{
		registryName: ent.RegistryName,
		Key:          ent.Key,
		Value:        ent.Value,
	}
}

func FormatAlias(registry, key string, value []string) aliasent.Alias {
	return aliasent.Alias{
		RegistryName: registry,
		Key:          key,
		Value:        value,
	}
}

func FormatEntityAliases(ents []aliasent.Alias) []Alias {
	var als = []Alias{}
	for _, v := range ents {
		als = append(als, FormatEntityAlias(v))
	}

	return als
}

type AliasRequest struct {
	Value []string `json:"value"`
}

type AliasResponse struct {
	Value []string `json:"value"`
}
