package aliaspg

import (
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

type Database struct {
	alias *AliasStore
}

func NewDatabase(pgClient postgres.Client) *Database {
	return &Database{
		alias: NewAlias(pgClient),
	}
}

func (db *Database) Alias() aliasent.AliasBackend {
	return db.alias
}
