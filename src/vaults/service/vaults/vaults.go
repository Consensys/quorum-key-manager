package vaults

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/vaults"
	"sync"
)

type Vaults struct {
	logger log.Logger
	mux    sync.RWMutex
	vaults map[string]*entities.Vault
	roles  auth.Roles
}

var _ vaults.Vaults = &Vaults{}

func New(roles auth.Roles, logger log.Logger) *Vaults {
	return &Vaults{
		logger: logger,
		mux:    sync.RWMutex{},
		vaults: make(map[string]*entities.Vault),
		roles:  roles,
	}
}

// TODO: Move to in-memory data layer
func (c *Vaults) createVault(name, vaultType string, cli interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.vaults[name] = &entities.Vault{
		Name:      name,
		Client:    cli,
		VaultType: vaultType,
	}
}

// TODO: Move to data layer
func (c *Vaults) getVault(_ context.Context, name string) (*entities.Vault, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if vault, ok := c.vaults[name]; ok {
		return vault, nil
	}

	errMessage := "vault was not found"
	c.logger.Error(errMessage, "name", name)
	return nil, errors.NotFoundError(errMessage)
}
