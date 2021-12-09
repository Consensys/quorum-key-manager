package aliases

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/database"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type Aliases struct {
	db     database.Aliases
	logger log.Logger
}

var _ aliases.Aliases = &Aliases{}

func New(db database.Aliases, logger log.Logger) *Aliases {
	return &Aliases{db: db, logger: logger}
}
