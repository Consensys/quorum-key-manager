package keys

import (
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"

	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type Connector struct {
	store  stores.KeyStore
	db     database.Keys
	logger log.Logger
}

var _ stores.KeyStore = Connector{}

func NewConnector(store stores.KeyStore, db database.Keys, logger log.Logger) *Connector {
	return &Connector{
		store:  store,
		db:     db,
		logger: logger,
	}
}
