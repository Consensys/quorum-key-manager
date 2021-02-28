package backend

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/audit"
	policymanager "github.com/ConsenSysQuorum/quorum-key-manager/auth/policy/manager"
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

	// PolicyManager returns the PolicyManager
	PolicyManager() policymanager.Manager

	// Auditor returns the PolicyManager
	Auditor() audit.Auditor
}
