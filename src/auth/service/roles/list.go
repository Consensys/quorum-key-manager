package roles

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (i *Roles) List(_ context.Context, _ *entities.UserInfo) ([]string, error) {
	// TODO: Implement {Resource/Role}BAC for roles

	roles := make([]string, 0, len(i.roles))
	for role := range i.roles {
		roles = append(roles, role)
	}

	i.logger.Debug("roles listed successfully")
	return roles, nil
}
