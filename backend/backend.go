package backend

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/store/manager"
)

// Backend synchronizes internal Key Manager components
// and provide access to it
type Backend interface {
	// StoreManager returns the Store Manager
	StoreManager() manager.Manager

	// ManifestLoader returns the Manifest Loader
	ManifestLoader() manifest.Loader
}
