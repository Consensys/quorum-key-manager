package policymanager

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/auth"
	"github.com/ConsenSysQuorum/quorum-key-manager/manifest"
)

// Manager allows to manage policies
type Manager interface {
	// Load policies from manifest messages
	Load(ctx context.Context, mnfsts ...*manifest.Message) error

	// Get auth for given client id, policies and metadata
	Get(ctx context.Context, id string, policies []string, metadata map[string]string) (*auth.Auth, error)
}
