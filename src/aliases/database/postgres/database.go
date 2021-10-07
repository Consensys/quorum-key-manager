package postgres

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasdb "github.com/consensys/quorum-key-manager/src/aliases/database"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

var _ aliasdb.Database = &Database{}

type Database struct {
	alias *Alias
}

func NewDatabase(pgClient postgres.Client, logger log.Logger) *Database {
	return &Database{
		alias: NewAlias(pgClient, logger),
	}
}

func (db *Database) Alias() aliases.Interactor {
	return db.alias
}
