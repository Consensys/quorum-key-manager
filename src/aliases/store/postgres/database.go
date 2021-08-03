package aliaspg

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

type Database struct {
	alias aliases.Alias
}

func NewDatabase(pgClient postgres.Client) *Database {
	return &Database{
		alias: NewAlias(pgClient),
	}
}

func (db *Database) Alias() aliases.Alias {
	return db.alias
}
