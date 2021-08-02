package aliaspg

import (
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

type Database struct {
	alias aliasstore.Alias
}

func NewDatabase(pgClient postgres.Client) *Database {
	return &Database{
		alias: NewAlias(pgClient),
	}
}

func (db *Database) Alias() aliasstore.Alias {
	return db.alias
}
