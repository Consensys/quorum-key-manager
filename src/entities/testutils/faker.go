package testutils

import (
	"time"

	"github.com/consensys/quorum-key-manager/pkg/common"
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

func FakeAliasRegistry() *entities.AliasRegistry {
	return &entities.AliasRegistry{
		Name:           common.RandString(10),
		AllowedTenants: []string{"tenant_1"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
