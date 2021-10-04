package aliaspg

import (
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

type Database struct {
	alias *AliasStore
}

func NewDatabase(pgClient postgres.Client, logger log.Logger) *Database {
	return &Database{
		alias: NewAlias(pgClient, logger),
	}
}

func (db *Database) Alias() aliasent.AliasBackend {
	return db.alias
}
