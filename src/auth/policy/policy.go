package policy

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows to manage policies and groups
type Manager interface {
	// Policy returns policy
	Policy(ctx context.Context, name string) (*types.Policy, error)

	// Policies returns all policies
	Policies(context.Context) ([]string, error)

	// Group returns group
	Group(ctx context.Context, name string) (*types.Group, error)

	// Groups returns groups
	Groups(context.Context) ([]string, error)

	// Extract User Policies from UserInfo
	UserPolicies(ctx context.Context, info *types.UserInfo) []types.Policy
}
