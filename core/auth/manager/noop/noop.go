package noop

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/auth"
	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/core/manifest/loader"
)

// Manager is a no operation auth manager
type Manager struct{}

func New() *Manager {
	return &Manager{}
}

// Load policies from manifest messages
func (mngr *Manager) Load(_ context.Context, _ ...*manifestloader.Message) error {
	return nil
}

// Get auth for given client id, policies and metadata
func (mngr *Manager) Get(_ context.Context, _ string, _ []string, _ map[string]string) (*auth.Auth, error) {
	return auth.NotAuthenticatedAuth, nil
}
