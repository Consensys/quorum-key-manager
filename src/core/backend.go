// nolint
package core

import (
	nodemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/node-manager"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager"
)

//go:generate mockgen -source=backend.go -destination=mocks/backend.go -package=mocks

// Backend holds internal Key Manager components and
// makes it available for API components
type Backend interface {
	// StoreManager returns the Store Manager
	StoreManager() storemanager.StoreManager

	// NodeManager returns the Node Manager
	NodeManager() nodemanager.Manager
}

type BaseBackend struct {
	storeMngr storemanager.StoreManager
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

func (bckend *BaseBackend) StoreManager() storemanager.StoreManager {
	return bckend.storeMngr
}

func (bckend *BaseBackend) NodeManager() nodemanager.Manager {
	return bckend.nodeMngr
}
