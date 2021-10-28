package auth

import (
	"context"
	"crypto/tls"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock

// Authenticator retrieves user info given an authentication method
type Authenticator interface {
	AuthenticateJWT(ctx context.Context, token string) (*entities.UserInfo, error)
	AuthenticateAPIKey(ctx context.Context, apiKey string) (*entities.UserInfo, error)
	AuthenticateTLS(ctx context.Context, connState *tls.ConnectionState) (*entities.UserInfo, error)
}

// Authorizator allows managing authorizations given a set of permissions
type Authorizator interface {
	CheckPermission(ops ...*entities.Operation) error
	CheckAccess(allowedTenants []string) error
}

// Manager allows managing policies and roles
type Manager interface {
	// Role returns role
	Role(name string) (*entities.Role, error)

	// Roles returns roles
	Roles() ([]string, error)

	// UserPermissions Extract User Permissions from UserInfo
	UserPermissions(info *entities.UserInfo) []entities.Permission
}
