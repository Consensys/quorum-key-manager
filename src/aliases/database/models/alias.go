package models

import (
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

// Alias allows the user to associates a RegistryName + a Key to 1 or more
// public keys stored in Value.
type Alias struct {
	tableName struct{} `pg:"aliases"` // nolint:unused,structcheck // reason

	// Key is the unique alias key.
	Key string `pg:",pk"`

	// RegistryName is the unique registry name.
	RegistryName string `pg:",pk"`

	// Value is a slice containing Tessera/Orion keys base64 encoded in strings.
	Value []string
}

// AliasFromEntitiy transforms an alias entity into an alias model.
func AliasFromEntity(ent aliasent.Alias) (alias Alias) {
	return Alias{
		Key:          ent.Key,
		RegistryName: ent.RegistryName,
		Value:        ent.Value,
	}
}

// ToEntity transforms an alias model into an alias entity.
func (a *Alias) ToEntity() *aliasent.Alias {
	return &aliasent.Alias{
		Key:          a.Key,
		RegistryName: a.RegistryName,
		Value:        a.Value,
	}
}

// AliasesToEntity transforms an alias model slice into an alias entity slice.
func AliasesToEntity(aliases []Alias) []aliasent.Alias {
	var ents []aliasent.Alias
	for _, v := range aliases {
		ent := v.ToEntity()
		ents = append(ents, *ent)
	}
	return ents
}
