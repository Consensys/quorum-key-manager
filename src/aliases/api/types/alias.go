package types

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/aliases/entities"
)

type Alias struct {
	Key   string     `json:"key"`
	Value AliasValue `json:"value"`

	registryName string
}

func FormatAliasValue(aliasValue AliasValue) entities.AliasValue {
	return entities.AliasValue{
		Kind:  aliasValue.RawKind,
		Value: aliasValue.RawValue,
	}
}

type AliasValue struct {
	RawKind  entities.Kind `json:"type"`
	RawValue interface{}   `json:"value"`
}

func (av AliasValue) MarshalJSON() ([]byte, error) {
	fmt.Fprintf(os.Stderr, "marshal1: %+v\n", av)
	switch av.RawKind {
	case entities.KindArray, entities.KindString:
		fmt.Fprintf(os.Stderr, "marshal: %+v\n", av)
	default:
		return nil, errors.InvalidFormatError(`bad alias value type: "%v"`, av.RawKind)
	}
	type loc AliasValue
	l := loc(av)
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "%+v", av)

	return b, nil

}

func (av *AliasValue) UnmarshalJSON(b []byte) error {
	//fmt.Fprintf(os.Stderr, "in unmarshal: %+v, %s\n", av, b)
	//type loc AliasValue
	type loc struct {
		RawKind  entities.Kind   `json:"type"`
		RawValue json.RawMessage `json:"value"`
	}
	var vv loc

	err := json.Unmarshal(b, &vv)
	//fmt.Fprintf(os.Stderr, "unmarshal2: %+v, %v\n", vv, err)
	if err != nil {
		return err
	}
	*av = AliasValue{
		RawKind: vv.RawKind,
	}

	switch av.RawKind {
	case entities.KindArray:
		var array []string
		err = json.Unmarshal(vv.RawValue, &array)
		if err != nil {
			return errors.InvalidFormatError(`bad alias array value: "%+v"`, string(vv.RawValue))
		}

		av.RawValue = array
	case entities.KindString:
		var s string
		err = json.Unmarshal(vv.RawValue, &s)
		if err != nil {
			return errors.InvalidFormatError(`bad alias string value: "%+v"`, string(vv.RawValue))
		}

		av.RawValue = s
		//fmt.Fprintf(os.Stderr, "unmarshal4: %+v, %v", vv, err)
	default:
		//fmt.Fprintf(os.Stderr, "unmarshal44: %+v, %v", vv, err)
		return errors.InvalidFormatError(`bad alias value type: "%v"`, av.RawKind)
	}

	//fmt.Fprintf(os.Stderr, "unmarshal3: %+v, %v\n", vv, err)
	//*av = AliasValue(vv)
	return nil
}

// FormatEntityAlias format an alias entity to an alias API type.
func FormatEntityAlias(ent entities.Alias) Alias {
	av := AliasValue{
		RawKind:  ent.Value.Kind,
		RawValue: ent.Value.Value,
	}

	return Alias{
		registryName: ent.RegistryName,
		Key:          ent.Key,
		Value:        av,
	}
}

// FormatAlias format an alias API type to an alias entity.
func FormatAlias(registry, key string, value AliasValue) entities.Alias {
	av := entities.AliasValue{
		Kind:  value.RawKind,
		Value: value.RawValue,
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
		als = append(als, FormatEntityAlias(v))
	}

	return als
}

// AliasRequest creates or modifies an alias value.
type AliasRequest struct {
	AliasValue
}

// AliasResponse returns the alias value.
type AliasResponse struct {
	AliasValue
}
