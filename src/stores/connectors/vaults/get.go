package vaults

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c *Connector) Get(_ context.Context, name string) (*entities.Vault, error) {
	logger := c.logger.With("name", name)

	c.mux.RLock()
	defer c.mux.RUnlock()

	if vault, ok := c.vaults[name]; ok {
		logger.Debug("vault found successfully")
		return vault, nil
	}

	errMessage := "vault was not found"
	c.logger.Error(errMessage)
	return nil, errors.NotFoundError(errMessage)
}
