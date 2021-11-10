package testutils

import (
	"github.com/consensys/quorum-key-manager/src/aliases/entities"
)

func NewEntAlias(registry, key string, value entities.AliasValue) entities.Alias {
	return entities.Alias{
		RegistryName: registry,
		Key:          key,
		Value:        value,
	}
}
