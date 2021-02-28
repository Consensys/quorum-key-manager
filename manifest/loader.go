package manifest

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/auth"
)

// Message wraps a manifest with information related to the Loader that loaded it
type Message struct {
	// Name of the loader that loaded the manifest
	Loader string

	// Manifest loaded
	Manifest *Manifest

	// Auth attach to the manifest when loading
	Auth *auth.Auth

	// Action to perform (e.g. create, update, delete...)
	Action string

	// Err while loading manifest
	Err error
}

// Loader loads and broadcast manifests
type Loader interface {
	// Subscribe creates a subscription that will write all Manifest matching the
	// given kinds to the given mnfsts channel

	// If no kind is passed then all manifest are written
	Subscribe(mnfsts chan<- []*Message) (Subscription, error)
}

// Subscription
type Subscription interface {
	Unsubscribe() error
	Error() <-chan error
}
