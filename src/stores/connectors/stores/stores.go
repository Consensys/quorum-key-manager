package stores

import (
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"sync"
)

type Connector struct {
	logger      log.Logger
	mux         sync.RWMutex
	authManager auth.Manager

	secrets     map[string]*storeBundle
	keys        map[string]*storeBundle
	ethAccounts map[string]*storeBundle

	db database.Database
}

type storeBundle struct {
	manifest *manifest.Manifest
	logger   log.Logger
	store    interface{}
}

var _ stores.Stores = &Connector{}

func NewConnector(authMngr auth.Manager, db database.Database, logger log.Logger) *Connector {
	return &Connector{
		logger:      logger,
		mux:         sync.RWMutex{},
		authManager: authMngr,
		secrets:     make(map[string]*storeBundle),
		keys:        make(map[string]*storeBundle),
		ethAccounts: make(map[string]*storeBundle),
		db:          db,
	}
}
