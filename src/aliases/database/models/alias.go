package models

import (
	"github.com/consensys/quorum-key-manager/src/aliases/entities"
)

// Alias allows the user to associates a RegistryName + a Key to 1 or more
// public keys stored in Value.
type Alias struct {
	tableName struct{} `pg:"aliases"` // nolint:unused,structcheck // reason

	// Key is the unique alias key.
	Key string `pg:",pk"`

	// RegistryName is the unique registry name.
	RegistryName string `pg:",pk"`

	Value entities.AliasValue
}

// AliasFromEntitiy transforms an alias entity into an alias model.
func AliasFromEntity(ent entities.Alias) Alias {
	av := Alias{
		Key:          ent.Key,
		RegistryName: ent.RegistryName,
		Value:        ent.Value,
	}

	return av
}

// ToEntity transforms an alias model into an alias entity.
func (a *Alias) ToEntity() *entities.Alias {
	return &entities.Alias{
		Key:          a.Key,
		RegistryName: a.RegistryName,
		//TODO check for casting
		Value: a.Value,
	}
}

// AliasesToEntity transforms an alias model slice into an alias entity slice.
func AliasesToEntity(aliases []Alias) []entities.Alias {
	var ents []entities.Alias
	for _, v := range aliases {
		ent := v.ToEntity()
		ents = append(ents, *ent)
	}
	return ents
}
