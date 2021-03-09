package manifestloader

import (
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/auth"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
)

// Message wraps a manifest with information related to the Loader that loaded it
type Message struct {
	// Name of the loader that loaded the manifest
	Loader string

	// Manifest loaded
	Manifest *manifest.Manifest

	// Auth attach to the manifest when loading
	Auth *auth.Auth

	// Action to perform (e.g. create, update, delete...)
	Action string

	// Err while loading manifest
	Err error
}

func (msg *Message) UnmarshalSpecs(specs interface{}) {
	err := msg.Manifest.UnmarshalSpecs(specs)
	if err != nil {
		msg.Err = fmt.Errorf("invalid specs format: %v", err)
	}
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
