package models

import (
	"github.com/consensys/quorum-key-manager/src/entities"
	"time"
)

type Alias struct {
	tableName struct{} `pg:"aliases"` // nolint:unused,structcheck // reason

	RegistryName string
	Key          string `pg:",pk"`
	Value        AliasValue
	CreatedAt    time.Time `pg:"default:now()"`
	UpdatedAt    time.Time `pg:"default:now()"`
}

type AliasValue struct {
	Kind  string
	Value interface{}
}

func NewAlias(alias *entities.Alias) *Alias {
	return &Alias{
		RegistryName: alias.RegistryName,
		Key:          alias.Key,
		Value: AliasValue{
			Kind:  alias.Kind,
			Value: alias.Value,
		},
		CreatedAt: alias.CreatedAt,
		UpdatedAt: alias.UpdatedAt,
	}
}

func (a *Alias) ToEntity() *entities.Alias {
	return &entities.Alias{
		Key:          a.Key,
		RegistryName: a.RegistryName,
		Kind:         a.Value.Kind,
		Value:        a.Value.Value,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
