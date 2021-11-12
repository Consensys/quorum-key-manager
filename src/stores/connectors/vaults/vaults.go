package vaults

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"sync"
)

type Connector struct {
	logger log.Logger
	mux    sync.RWMutex
	vaults map[string]*entities.Vault
}

var _ stores.Vaults = &Connector{}

func NewConnector(logger log.Logger) *Connector {
	return &Connector{
		logger: logger,
		mux:    sync.RWMutex{},
		vaults: make(map[string]*entities.Vault),
	}
}

// TODO: Move to in-memory data layer
func (c *Connector) createVault(name, vaultType string, cli interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.vaults[name] = &entities.Vault{
		Name:      name,
		Client:    cli,
		VaultType: vaultType,
	}
}

// TODO: Move to data layer
func (c *Connector) getVault(_ context.Context, name string) (*entities.Vault, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if vault, ok := c.vaults[name]; ok {
		return vault, nil
	}

	errMessage := "vault was not found"
	c.logger.Error(errMessage, "name", name)
	return nil, errors.NotFoundError(errMessage)
}
