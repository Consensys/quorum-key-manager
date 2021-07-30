package keys

import (
	"github.com/consensys/quorum-key-manager/src/stores/connectors"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
)

type Connector struct {
	store  keys.Store
	db     database.Database
	logger log.Logger
}

var _ connectors.KeysConnector = Connector{}

func NewConnector(store keys.Store, db database.Database, logger log.Logger) *Connector {
	return &Connector{
		store:  store,
		db:     db,
		logger: logger,
	}
}
