// nolint
package core

import (
	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest/loader"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager"
	basemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/base"
)

// Backend holds internal Key Manager components and
// makes it available for API components
type Backend interface {
	// StoreManager returns the Store Manager
	StoreManager() storemanager.Manager

	// ManifestLoader returns the Manifest Loader
	ManifestLoader() manifestloader.Loader
}

type BaseBackend struct {
	storeMnger storemanager.Manager
	loader     manifestloader.Loader
}

func New() *BaseBackend {

	// Create store manager
	storeMngr := basemanager.New()

	return &BaseBackend{
		storeMnger: storeMngr,
	}
}

func (bckend *BaseBackend) StoreManager() storemanager.Manager {
	return bckend.storeMnger
}

func (bckend *BaseBackend) ManifestLoader() manifestloader.Loader {
	return bckend.loader
}
