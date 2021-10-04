package aliasmodels

import (
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

// Alias allows the user to associates a RegistryName + a Key to 1 or more public keys stored
// in Value. The Value has 2 formats:
// - a JSON string if AliasKind is an AliasKindString.
// - a JSON array of strings if AliasKind is an AliasKindArray.
type Alias struct {
	tableName struct{} `pg:"aliases"` // nolint:unused,structcheck // reason

	Key          string `pg:",pk"`
	RegistryName string `pg:",pk"`
	// Value is a JSON array containing Tessera/Orion keys base64 encoded in strings.
	Value []string
}

func AliasFromEntity(ent aliasent.Alias) (alias Alias) {
	return Alias{
		Key:          ent.Key,
		RegistryName: ent.RegistryName,
		Value:        ent.Value,
	}
}

func (a *Alias) ToEntity() *aliasent.Alias {
	return &aliasent.Alias{
		Key:          a.Key,
		RegistryName: a.RegistryName,
		Value:        a.Value,
	}
}

func AliasesToEntity(aliases []Alias) []aliasent.Alias {
	var ents []aliasent.Alias
	for _, v := range aliases {
		ent := v.ToEntity()
		ents = append(ents, *ent)
	}
	return ents
}
