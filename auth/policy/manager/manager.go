package policymanager

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/auth/policy"
	"github.com/ConsenSysQuorum/quorum-key-manager/manifest"
)

// Manager allows to manage policies
type Manager interface {
	// Load policies from manifest messages
	Load(ctx context.Context, mnfsts ...*manifest.Message) error

	// Get policies for given client id and metadata
	Get(ctx context.Context, id string, metadata map[string]string) ([]*policy.Policy, error)
}
