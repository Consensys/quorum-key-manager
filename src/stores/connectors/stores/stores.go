package stores

import (
	"context"
	"sync"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

type Connector struct {
	logger      log.Logger
	mux         sync.RWMutex
	authManager auth.Manager

	stores map[string]*entities.StoreInfo

	db database.Database
}

var _ stores.Stores = &Connector{}

func NewConnector(authMngr auth.Manager, db database.Database, logger log.Logger) *Connector {
	return &Connector{
		logger:      logger,
		mux:         sync.RWMutex{},
		authManager: authMngr,
		stores:      make(map[string]*entities.StoreInfo),
		db:          db,
	}
}

func (c *Connector) getStore(_ context.Context, storeName string, resolver auth.Authorizator) (*entities.StoreInfo, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if bundle, ok := c.stores[storeName]; ok {
		if err := resolver.CheckAccess(bundle.AllowedTenants); err != nil {
			return nil, err
		}

		return bundle, nil
	}

	errMessage := "store was not found"
	c.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}
