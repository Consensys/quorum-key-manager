package manager

import (
	"fmt"

	manifest "github.com/consensysquorum/quorum-key-manager/src/manifests/types"
)

type Action string

const CreateAction = "Create"
const UpdateAction = "Update"
const DeleteAction = "Update"

// Message wraps a manifest with information related to the Loader that loaded it
type Message struct {
	// Name of the loader that loaded the manifest
	Loader string

	// Manifest loaded
	Manifest *manifest.Manifest

	// Action to perform (e.g. create, update, delete...)
	Action Action

	// Err while loading manifest
	Err error
}

func (msg *Message) UnmarshalSpecs(specs interface{}) {
	err := msg.Manifest.UnmarshalSpecs(specs)
	if err != nil {
		msg.Err = fmt.Errorf("invalid specs format: %v", err)
	}
}

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager manages manifests
type Manager interface {
	// Subscribe creates a subscription that writes all
	// Manifests matching kinds on the given channel

	// If kinds is nil then all manifest are written
	Subscribe(kinds []manifest.Kind, messages chan<- []Message) (Subscription, error)
}

// Subscription
type Subscription interface {
	Unsubscribe() error
	Error() <-chan error
}
