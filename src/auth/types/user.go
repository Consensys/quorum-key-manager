package types

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
)

// UserInfo are extracted from request credentials by authentication middleware
type UserInfo struct {
	// AuthMode records the mode that succeeded to Authenticate the request ('tls', 'api-key', 'oidc' or '')
	AuthMode string

	// Tenant belonged by the user
	Tenant string

	// Subject identifies the user
	Username string

	// Roles indicates the user's membership
	Roles []string

	// Permissions specify
	Permissions []Permission
}

func (ui *UserInfo) CheckAccess(mnf *manifest.Manifest) error {
	if len(mnf.AllowedTenants) == 0 {
		return nil
	}

	if ui.Tenant == "" {
		return errors.UnauthorizedError("missing credentials")
	}

	for _, t := range mnf.AllowedTenants {
		if t == ui.Tenant {
			return nil
		}
	}

	return errors.NotFoundError("tenant %s does not have access to %s", ui.Tenant, mnf.Name)
}

var AnonymousUser = &UserInfo{
	Username:    "anonymous",
	Roles:       []string{AnonymousRole},
	Permissions: []Permission{},
}

var AuthenticatedUser = &UserInfo{
	Roles: []string{AnonymousRole},
}
