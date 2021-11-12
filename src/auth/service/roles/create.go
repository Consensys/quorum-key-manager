package roles

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (i *Interactor) Create(ctx context.Context, name string, permissions []entities.Permission, _ *entities.UserInfo) error {
	logger := i.logger.With("name", name, "permissions", permissions)
	logger.Debug("creating role")

	// TODO: Implement {Resource/Role}BAC for roles

	if _, ok := i.roles[name]; ok {
		errMessage := fmt.Sprintf("role %s already exist", name)
		logger.Error(errMessage)
		return errors.AlreadyExistsError(errMessage)
	}

	i.createRole(ctx, name, permissions)

	logger.Info("role created successfully")
	return nil
}
