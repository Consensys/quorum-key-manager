package manager

import (
	"context"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/auth/types"
)

type Auth interface {
	Policies(ctx context.Context, name string) (*types.Policy, error)

	ListPolicies(context.Context) ([]*types.Policy, error)

	Group(ctx context.Context, name string) (*types.Group, error)

	ListGroups(context.Context) ([]*types.Group, error)

	// Authenticate request
	Authenticate(req *http.Request) error
}
