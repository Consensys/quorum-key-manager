package roles

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (i *Roles) UserPermissions(ctx context.Context, userInfo *entities.UserInfo) []entities.Permission {
	if userInfo == nil {
		return []entities.Permission{}
	}

	var permissions []entities.Permission
	copy(permissions, userInfo.Permissions)

	for _, roleName := range userInfo.Roles {
		role, err := i.Get(ctx, roleName, userInfo)
		if err != nil {
			continue
		}

		permissions = append(permissions, role.Permissions...)
		for _, p := range role.Permissions {
			permissions = append(permissions, entities.ListWildcardPermission(string(p))...)
		}
	}

	i.logger.Debug("permissions extracted successfully", "tenant", userInfo.Tenant, "username", userInfo.Username, "permissions", permissions)
	return permissions
}
