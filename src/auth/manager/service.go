package manager

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensys/quorum-key-manager/src/infra/manifests"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

const ID = "AuthManager"

type BaseManager struct {
	manifestReader manifests.Reader

	mux   sync.RWMutex
	roles map[string]*types.Role

	logger log.Logger
}

func New(manifestReader manifests.Reader, logger log.Logger) *BaseManager {
	return &BaseManager{
		manifestReader: manifestReader,
		roles:          make(map[string]*types.Role),
		logger:         logger,
	}
}

func (mngr *BaseManager) Start(_ context.Context) error {
	mnfs, err := mngr.manifestReader.Load()
	if err != nil {
		errMessage := "failed to load manifest file"
		mngr.logger.WithError(err).Error(errMessage)
		return errors.ConfigError(errMessage)
	}

	for _, mnf := range mnfs {
		_ = mngr.load(mnf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mngr *BaseManager) Stop(context.Context) error { return nil }
func (mngr *BaseManager) Error() error               { return nil }
func (mngr *BaseManager) Close() error               { return nil }

func (mngr *BaseManager) UserPermissions(user *types.UserInfo) []types.Permission {
	var permissions []types.Permission
	if user == nil {
		return permissions
	}

	permissions = append(permissions, user.Permissions...)

	for _, roleName := range user.Roles {
		role, err := mngr.Role(roleName)
		if err != nil {
			mngr.logger.WithError(err).With("role", roleName).Debug("could not load role")
			continue
		}

		permissions = append(permissions, role.Permissions...)
		for _, p := range role.Permissions {
			permissions = append(permissions, types.ListWildcardPermission(string(p))...)
		}
	}

	return permissions
}

func (mngr *BaseManager) Role(name string) (*types.Role, error) {
	if group, ok := mngr.roles[name]; ok {
		return group, nil
	}

	errMessage := "role not found"
	mngr.logger.With("name", name).Error(errMessage)
	return nil, errors.NotFoundError(errMessage)
}

func (mngr *BaseManager) Roles() ([]string, error) {
	roles := make([]string, 0, len(mngr.roles))
	for role := range mngr.roles {
		roles = append(roles, role)
	}

	return roles, nil
}

func (mngr *BaseManager) load(mnf *manifest.Manifest) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()

	logger := mngr.logger.With("name", mnf.Name)

	if mnf.Kind == manifest.Role {
		if _, ok := mngr.roles[mnf.Name]; ok {
			errMessage := fmt.Sprintf("role %s already exist", mnf.Name)
			logger.Error(errMessage)
			return errors.AlreadyExistsError(errMessage)
		}

		specs := new(RoleSpecs)
		if err := mnf.UnmarshalSpecs(specs); err != nil {
			errMessage := fmt.Sprintf("invalid Role specs for role %s", mnf.Name)
			logger.WithError(err).Error(errMessage)
			return errors.InvalidParameterError(errMessage)
		}

		mngr.roles[mnf.Name] = &types.Role{
			Name:        mnf.Name,
			Permissions: specs.Permissions,
		}

		logger.Info("Role created successfully")
		return nil
	}

	return nil
}

func (mngr *BaseManager) ID() string                             { return ID }
func (mngr *BaseManager) CheckLiveness(_ context.Context) error  { return nil }
func (mngr *BaseManager) CheckReadiness(_ context.Context) error { return mngr.Error() }
