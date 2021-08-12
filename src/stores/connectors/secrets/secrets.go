package secrets

import (
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

type Connector struct {
	store  stores.SecretStore
	logger log.Logger
	db     database.Secrets
}

var _ stores.SecretStore = &Connector{}

func NewSecretConnector(store stores.SecretStore, db database.Secrets, logger log.Logger) *Connector {
	return &Connector{
		store:  store,
		logger: logger,
		db:     db,
	}
}
