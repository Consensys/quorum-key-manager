package testutils

import (
	"github.com/consensys/quorum-key-manager/src/entities"
)

func AliasFaker(registry, key string, kind entities.AliasKind, value interface{}) entities.Alias {
	return entities.Alias{
		RegistryName: registry,
		Key:          key,
		Kind:         kind,
		Value:        value,
	}
}
