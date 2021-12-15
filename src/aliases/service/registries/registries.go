package registries

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/database"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type Registries struct {
	db           database.Registry
	logger       log.Logger
	authorizator auth.Authorizator
}

var _ aliases.Registries = &Registries{}

func New(db database.Registry, authorizator auth.Authorizator, logger log.Logger) *Registries {
	return &Registries{
		db:           db,
		logger:       logger,
		authorizator: authorizator,
	}
}
