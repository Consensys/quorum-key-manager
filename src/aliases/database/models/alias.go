package models

import (
	"log"

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

	Kind string

	// StringValue is a slice containing Tessera/Orion keys base64 encoded in strings.
	StringValue string
	ArrayValue  []string
	Value       entities.AliasValue
}

// AliasFromEntitiy transforms an alias entity into an alias model.
func AliasFromEntity(ent entities.Alias) Alias {
	av := Alias{
		Key: ent.Key,
		//TODO check for casting
		RegistryName: ent.RegistryName,
		Kind:         string(ent.Value.Kind),
	}
	switch ent.Value.Kind {
	case entities.KindArray:
		av.ArrayValue = ent.Value.Value.([]string)
	case entities.KindString:
		av.StringValue = ent.Value.Value.(string)
	default:
		// TODO the: add real error handling
		// we should make sure the error handling is when we save/get from the DB
		log.Fatal("AliasFromEntity: bad value type")
	}
	return av
}

// ToEntity transforms an alias model into an alias entity.
func (a *Alias) ToEntity() *entities.Alias {
	av := entities.AliasValue{
		Kind: entities.Kind(a.Kind),
	}
	switch av.Kind {
	case entities.KindArray:
		av.Value = a.ArrayValue
	case entities.KindString:
		av.Value = a.ArrayValue
	default:
		// TODO the: add real error handling
		// we should make sure the error handling is when we save/get from the DB
		log.Fatal("Alias.ToEntity: bad value type")
	}

	return &entities.Alias{
		Key:          a.Key,
		RegistryName: a.RegistryName,
		//TODO check for casting
		Value: av,
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
