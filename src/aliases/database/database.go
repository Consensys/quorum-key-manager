package aliasdb

import "github.com/consensys/quorum-key-manager/src/aliases"

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	Alias() aliases.AliasBackend
}
