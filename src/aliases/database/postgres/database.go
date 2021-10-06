package aliaspg

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
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

func (db *Database) Alias() aliases.Repository {
	return db.alias
}
