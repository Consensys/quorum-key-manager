package entities

import (
	manifest "github.com/consensys/quorum-key-manager/src/entities"
)

type StoreInfo struct {
	AllowedTenants []string
	Store     interface{}
	StoreType manifest.StoreType
	VaultType manifest.VaultType
}
