package auth

import (
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows managing policies and roles
type Manager interface {
	// Role returns role
	Role(name string) (*types.Role, error)

	// Roles returns roles
	Roles() ([]string, error)

	// UserPermissions Extract User Permissions from UserInfo
	UserPermissions(info *types.UserInfo) []types.Permission
}
