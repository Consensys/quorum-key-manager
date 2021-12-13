package testutils

import (
	"github.com/consensys/quorum-key-manager/src/entities"
)

func FakeAlias(registry, key, kind string, value interface{}) *entities.Alias {
	return &entities.Alias{
		Key:          key,
		RegistryName: registry,
		Kind:         kind,
		Value:        value,
	}
}
