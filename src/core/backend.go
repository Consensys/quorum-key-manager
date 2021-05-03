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
	stores := storemanager.New()
	nodes := nodemanager.New(stores)

	return &BaseBackend{
		storeMngr: stores,
		nodeMngr:  nodes,
	}
}
