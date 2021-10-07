package testutils

import (
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

func NewEntAlias(registry, key string, value []string) aliasent.Alias {
	return aliasent.Alias{
		RegistryName: registry,
		Key:          key,
		Value:        value,
	}
}
