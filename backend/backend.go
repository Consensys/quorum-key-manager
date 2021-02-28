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
	auditor     audit.Auditor
	authManager authmanager.Manager
}

func New() *BaseBackend {
	bckend := &BaseBackend{
		auditor:     noopauditor.New(),
		authManager: noopauthmanager.New(),
	}

	return bckend
}
