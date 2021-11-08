package entities

import (
	"fmt"
)

// Alias allows the user to associates a RegistryName + a Key to 1 or more
// public keys stored in Value.
type Alias struct {
	// Key is the unique alias key.
	Key string

	// RegistryName is the unique registry name.
	RegistryName string

	// Value is a slice containing Tessera/Orion keys base64 encoded in strings.
	Value AliasValue
}

type Kind string

const (
	KindUnknown Kind = ""
	KindString  Kind = "string"
	KindArray   Kind = "array"
)

type AliasValue struct {
	Kind  Kind
	Value interface{}
}

func (av AliasValue) Array() ([]string, error) {
	keys, ok := av.Value.([]interface{})
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
func (av AliasValue) String() (string, error) {
	value, ok := av.Value.(string)
	if !ok {
		return "", fmt.Errorf("alias value is not a string")
	}
	return value, nil
}
