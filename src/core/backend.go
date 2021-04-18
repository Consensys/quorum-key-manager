// nolint
package core

import (
	nodemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/node-manager"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager"
)

// Backend holds internal Key Manager components and
// makes it available for API components
type Backend interface {
	// StoreManager returns the Store Manager
	StoreManager() storemanager.Manager

	// NodeManager returns the Node Manager
	NodeManager() nodemanager.Manager
}

type BaseBackend struct {
	storeMngr storemanager.Manager
	nodeMngr  nodemanager.Manager
}

func New() *BaseBackend {
	return &BaseBackend{
		storeMngr: storemanager.New(),
		nodeMngr:  nodemanager.New(),
	}
}

func (bckend *BaseBackend) StoreManager() storemanager.Manager {
	return bckend.storeMngr
}

func (bckend *BaseBackend) NodeManager() nodemanager.Manager {
	return bckend.nodeMngr
}
