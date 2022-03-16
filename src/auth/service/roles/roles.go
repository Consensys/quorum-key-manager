package roles

import (
	"context"
	"sync"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"

	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type Roles struct {
	mux    sync.RWMutex
	roles  map[string]*entities.Role
	logger log.Logger
}

var _ auth.Roles = &Roles{}

func New(logger log.Logger) *Roles {
	return &Roles{
		roles:  make(map[string]*entities.Role),
		logger: logger,
	}
}

// TODO: Move to data layer
func (i *Roles) createRole(_ context.Context, name string, permissions []entities.Permission) {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.roles[name] = &entities.Role{
		Name:        name,
		Permissions: permissions,
	}
}

// TODO: Move to data layer
func (i *Roles) getRole(_ context.Context, name string) (*entities.Role, error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	if role, ok := i.roles[name]; ok {
		return role, nil
	}

	return nil, errors.NotFoundError("role was not found")
}
