package entities

import (
	"fmt"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"time"
)

type AliasKind string

const (
	AliasKindString AliasKind = "string"
	AliasKindArray  AliasKind = "array"
)

// Alias allows the user to associates a RegistryName + a Key to 1 or more public keys stored in Value.
type Alias struct {
	// Key is the unique alias key.
	Key string

	// RegistryName is the unique registry name.
	RegistryName string

	Kind AliasKind

	// Value is a slice containing Tessera/Orion keys base64 encoded in strings.
	Value interface{}

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Alias) Array() ([]string, error) {
	keys, ok := a.Value.([]interface{})
	if !ok {
		return nil, fmt.Errorf(`alias value is not an array`)
	}

	var array []string
	for _, v := range keys {
		val, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf(`element in array is not a string: %+v`, v)
		}
		array = append(array, val)
	}

	return array, nil
}

func (a *Alias) String() (string, error) {
	value, ok := a.Value.(string)
	if !ok {
		return "", fmt.Errorf("alias value is not a string")
	}

	return value, nil
}

func (a *Alias) Validate() error {
	switch a.Kind {
	case AliasKindArray:
		_, err := a.Array()
		if err != nil {
			return errors.InvalidParameterError("alias value is not an array")
		}
	case AliasKindString:
		_, err := a.String()
		if err != nil {
			return errors.InvalidParameterError("alias value is not a string")
		}
	default:
		return errors.InvalidParameterError("invalid alias type")
	}

	return nil
}
