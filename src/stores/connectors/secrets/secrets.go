package secrets

import (
	"github.com/consensys/quorum-key-manager/src/auth/manager"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

type Connector struct {
	store    stores.SecretStore
	logger   log.Logger
	db       database.Secrets
	resolver *manager.Resolver
}

var _ stores.SecretStore = &Connector{}

func NewConnector(store stores.SecretStore, db database.Secrets, resolvr *manager.Resolver, logger log.Logger) *Connector {
	return &Connector{
		store:    store,
		logger:   logger,
		db:       db,
		resolver: resolvr,
	}
}
