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
	AuthenticateAPIKey(ctx context.Context, apiKey []byte) (*entities.UserInfo, error)
	AuthenticateTLS(ctx context.Context, connState *tls.ConnectionState) (*entities.UserInfo, error)
}

// Authorizator allows managing authorizations given a set of permissions
type Authorizator interface {
	CheckPermission(ops ...*entities.Operation) error
	CheckAccess(allowedTenants []string) error
}

// Roles allows managing permissions and roles
type Roles interface {
	Create(ctx context.Context, name string, permissions []entities.Permission, userInfo *entities.UserInfo) error
	Get(ctx context.Context, name string, userInfo *entities.UserInfo) (*entities.Role, error)
	List(ctx context.Context, userInfo *entities.UserInfo) ([]string, error)
	UserPermissions(ctx context.Context, userInfo *entities.UserInfo) []entities.Permission
}
