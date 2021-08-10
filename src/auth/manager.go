package auth

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows to manage policies and roles
type Manager interface {
	// Role returns role
	Role(ctx context.Context, name string) (*types.Role, error)

	// Roles returns roles
	Roles(context.Context) ([]string, error)

	// Extract User Permissions from UserInfo
	UserPermissions(ctx context.Context, info *types.UserInfo) []types.Permission
}
