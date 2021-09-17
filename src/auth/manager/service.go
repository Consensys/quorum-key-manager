package manager

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/auth/types"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
)

const ID = "AuthManager"

var authKinds = []manifest.Kind{
	RoleKind,
}

type BaseManager struct {
	manifests manifestsmanager.Manager

	mux   sync.RWMutex
	roles map[string]*types.Role

	sub    manifestsmanager.Subscription
	mnfsts chan []manifestsmanager.Message

	logger log.Logger
	err    error
	isLive bool
}

func New(manifests manifestsmanager.Manager, logger log.Logger) *BaseManager {
	return &BaseManager{
		manifests: manifests,
		roles:     make(map[string]*types.Role),
		mnfsts:    make(chan []manifestsmanager.Message),
		logger:    logger,
	}
}

func (mngr *BaseManager) Start(_ context.Context) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()
	defer func() {
		mngr.isLive = true
	}()

	// Subscribe to manifest of Kind Role and Policy
	mngr.sub = mngr.manifests.Subscribe(authKinds, mngr.mnfsts)

	// Start loading manifest
	go mngr.loadAll()

	return nil
}

func (mngr *BaseManager) Stop(context.Context) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()
	defer close(mngr.mnfsts)

	if mngr.sub != nil {
		_ = mngr.sub.Unsubscribe()
	}
	return nil
}

func (mngr *BaseManager) Error() error {
	return mngr.err
}

func (mngr *BaseManager) Close() error {
	close(mngr.mnfsts)
	return nil
}

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
	return mngr.role(name)
}

func (mngr *BaseManager) Roles() ([]string, error) {
	roles := make([]string, 0, len(mngr.roles))
	for role := range mngr.roles {
		roles = append(roles, role)
	}
	return roles, nil
}

func (mngr *BaseManager) role(name string) (*types.Role, error) {
	if group, ok := mngr.roles[name]; ok {
		return group, nil
	}

	return nil, fmt.Errorf("role %q not found", name)
}

func (mngr *BaseManager) loadAll() {
	for mnfsts := range mngr.mnfsts {
		for _, mnf := range mnfsts {
			_ = mngr.load(mnf.Manifest)
		}
	}
}

func (mngr *BaseManager) load(mnf *manifest.Manifest) error {
	mngr.mux.Lock()
	defer mngr.mux.Unlock()

	logger := mngr.logger.With("kind", mnf.Kind).With("name", mnf.Name)

	switch mnf.Kind {
	case RoleKind:
		err := mngr.loadRole(mnf)
		if err != nil {
			logger.WithError(err).Error("could not load Role")
			return err
		}
		logger.Info("loaded Role")
	default:
		err := fmt.Errorf("invalid manifest kind %s", mnf.Kind)
		logger.WithError(err).Error("error starting node")
		return err
	}

	return nil
}

func (mngr *BaseManager) loadRole(mnf *manifest.Manifest) error {
	if _, ok := mngr.roles[mnf.Name]; ok {
		return fmt.Errorf("role %q already exist", mnf.Name)
	}

	specs := new(RoleSpecs)
	if err := mnf.UnmarshalSpecs(specs); err != nil {
		return fmt.Errorf("invalid Role specs: %v", err)
	}

	mngr.roles[mnf.Name] = &types.Role{
		Name:        mnf.Name,
		Permissions: specs.Permissions,
	}

	return nil
}

func (mngr *BaseManager) ID() string { return ID }
func (mngr *BaseManager) CheckLiveness(_ context.Context) error {
	if mngr.isLive {
		return nil
	}

	errMessage := fmt.Sprintf("service %s is not live", mngr.ID())
	mngr.logger.Error(errMessage, "id", mngr.ID())
	return errors.HealthcheckError(errMessage)
}

func (mngr *BaseManager) CheckReadiness(_ context.Context) error {
	return mngr.Error()
}
