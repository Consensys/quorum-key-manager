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
	stores      map[string]*entities.Store
	vaults      stores.Vaults
	db          database.Database
}

var _ stores.Stores = &Connector{}

func NewConnector(authMngr auth.Manager, db database.Database, vaults stores.Vaults, logger log.Logger) *Connector {
	return &Connector{
		logger:      logger,
		mux:         sync.RWMutex{},
		authManager: authMngr,
		stores:      make(map[string]*entities.Store),
		vaults:      vaults,
		db:          db,
	}
}

// TODO: Move to data layer
func (c *Connector) createStore(name, storeType string, store interface{}, allowedTenants []string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.stores[name] = &entities.Store{
		Name:           name,
		AllowedTenants: allowedTenants,
		Store:          store,
		StoreType:      storeType,
	}
}

// TODO: Move to data layer
func (c *Connector) getStore(_ context.Context, name string, resolver auth.Authorizator) (*entities.Store, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if bundle, ok := c.stores[name]; ok {
		if err := resolver.CheckAccess(bundle.AllowedTenants); err != nil {
			return nil, err
		}

		return bundle, nil
	}

	errMessage := "store was not found"
	c.logger.Error(errMessage, "name", name)
	return nil, errors.NotFoundError(errMessage)
}
