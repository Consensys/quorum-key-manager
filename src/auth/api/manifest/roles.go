package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/api/types"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

type RolesHandler struct {
	roles    auth.Roles
	userInfo *entities.UserInfo
}

func NewRolesHandler(roles auth.Roles) *RolesHandler {
	return &RolesHandler{
		roles:    roles,
		userInfo: entities.NewWildcardUser(), // This handler always use the wildcard user because it's a manifest handler
	}
}

func (h *RolesHandler) Create(ctx context.Context, name string, specs interface{}) error {
	createReq := &types.CreateRoleRequest{}
	err := json.UnmarshalYAML(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.roles.Create(ctx, name, createReq.Permissions, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
