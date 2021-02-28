package backend

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/audit"
	noopauditor "github.com/ConsenSysQuorum/quorum-key-manager/audit/noop"
	authmanager "github.com/ConsenSysQuorum/quorum-key-manager/auth/manager"
	noopauthmanager "github.com/ConsenSysQuorum/quorum-key-manager/auth/manager/noop"
	"github.com/ConsenSysQuorum/quorum-key-manager/manifest"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/store/manager"
)

// Backend holds internal Key Manager components and
// makes it available for API components
type Backend interface {
	// StoreManager returns the Store Manager
	StoreManager() storemanager.Manager

	// ManifestLoader returns the Manifest Loader
	ManifestLoader() manifest.Loader

	// AuthManager returns the PolicyManager
	AuthManager() authmanager.Manager

	// Auditor returns the PolicyManager
	Auditor() audit.Auditor
}

type BaseBackend struct {
	auditor    audit.Auditor
	authMngr   authmanager.Manager
	storeMnger storemanager.Manager
	loader     manifest.Loader
}

func New() *BaseBackend {
	bckend := &BaseBackend{
		auditor:  noopauditor.New(),
		authMngr: noopauthmanager.New(),
	}

	return bckend
}

func (bckend *BaseBackend) StoreManager() storemanager.Manager {
	return bckend.storeMnger
}

func (bckend *BaseBackend) ManifestLoader() manifest.Loader {
	return bckend.loader
}

func (bckend *BaseBackend) AuthManager() authmanager.Manager {
	return bckend.authMngr
}

func (bckend *BaseBackend) Auditor() audit.Auditor {
	return bckend.auditor
}
