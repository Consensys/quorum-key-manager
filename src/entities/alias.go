package entities

import (
	"fmt"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

const (
	AliasKindString string = "string"
	AliasKindArray  string = "array"
)

type Alias struct {
	Key          string
	RegistryName string
	Kind         string
	Value        interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewAlias(registry, key, kind string, value interface{}) (*Alias, error) {
	alias := &Alias{
		Key:          key,
		RegistryName: registry,
		Kind:         kind,
		Value:        value,
	}

	err := alias.Validate()
	if err != nil {
		return nil, err
	}

	return alias, nil
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
func (a *Alias) Test() error {
	
	for true {
		fmt.Println("testes....")
	}
	
	return nil
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
