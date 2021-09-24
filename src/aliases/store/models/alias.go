package aliasmodels

import (
	"encoding/base64"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

// Alias allows the user to associates a RegistryName + a Key to 1 or more public keys stored
// in Value. The Value has 2 formats:
// - a JSON string if AliasKind is an AliasKindString.
// - a JSON array of strings if AliasKind is an AliasKindArray.
type Alias struct {
	tableName struct{} `pg:"aliases"` // nolint:unused,structcheck // reason

	Key          AliasKey     `pg:",pk"`
	RegistryName RegistryName `pg:",pk"`
	// Value is a JSON array containing Tessera/Orion keys base64 encoded in strings.
	Value AliasValue
}

func checkAliasValue(value aliasent.AliasValue) error {
	for _, v := range value {
		// we just want the decoding to happen without error
		_, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return errors.InvalidFormatError("alias value is incorrect: %v: %v", err, v)
		}
	}
	return nil
}

func AliasFromEntity(ent aliasent.Alias) (alias Alias) {
	return Alias{
		Key:          AliasKey(ent.Key),
		RegistryName: RegistryName(ent.RegistryName),
		Value:        AliasValue(ent.Value),
	}
}

func (a *Alias) ToEntity() *aliasent.Alias {
	return &aliasent.Alias{
		Key:          aliasent.AliasKey(a.Key),
		RegistryName: aliasent.RegistryName(a.RegistryName),
		Value:        aliasent.AliasValue(a.Value),
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

type AliasKey string

type AliasValue []string

type RegistryName string
