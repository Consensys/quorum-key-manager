// nolint
package core

import (
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager"
)

// Backend holds internal Key Manager components and
// makes it available for API components
type Backend interface {
	// StoreManager returns the Store Manager
	StoreManager() storemanager.Manager
}

type BaseBackend struct {
	storeMnger storemanager.Manager
}

func New() *BaseBackend {

	// Create store manager
	storeMngr := storemanager.New()

	return &BaseBackend{
		storeMnger: storeMngr,
	}
}

func (bckend *BaseBackend) StoreManager() storemanager.Manager {
	return bckend.storeMnger
}

