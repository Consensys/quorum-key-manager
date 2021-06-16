package manager

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

type Manager interface {
	Policy(ctx context.Context, name string) (*types.Policy, error)

	Policies(context.Context) ([]string, error)

	Group(ctx context.Context, name string) (*types.Group, error)

	Groups(context.Context) ([]string, error)
}
