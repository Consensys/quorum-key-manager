package aliases

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/database"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type Aliases struct {
	aliasDB    database.Alias
	registryDB database.Registry
	logger     log.Logger
}

var _ aliases.Aliases = &Aliases{}

func New(aliasDB database.Alias, registryDB database.Registry, logger log.Logger) *Aliases {
	return &Aliases{aliasDB: aliasDB, registryDB: registryDB, logger: logger}
}
