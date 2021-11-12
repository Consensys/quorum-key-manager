package roles

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"sync"

	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type Interactor struct {
	mux    sync.RWMutex
	roles  map[string]*entities.Role
	logger log.Logger
}

var _ auth.Roles = &Interactor{}

func New(logger log.Logger) *Interactor {
	return &Interactor{
		roles:  make(map[string]*entities.Role),
		logger: logger,
	}
}

// TODO: Move to data layer
func (i *Interactor) createRole(_ context.Context, name string, permissions []entities.Permission) {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.roles[name] = &entities.Role{
		Name:        name,
		Permissions: permissions,
	}
}

// TODO: Move to data layer
func (i *Interactor) getRole(_ context.Context, name string) (*entities.Role, error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	if role, ok := i.roles[name]; ok {
		return role, nil
	}

	errMessage := "role was not found"
	i.logger.Error(errMessage, "name", name)
	return nil, errors.NotFoundError(errMessage)
}
