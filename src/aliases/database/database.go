package aliasdb

import (
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	Alias() aliasent.AliasBackend
}
