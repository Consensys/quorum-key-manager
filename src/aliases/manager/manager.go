package aliasmanager

import (
	"context"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

type BaseManager struct {
	Aliases aliasent.AliasBackend
}

func New(aliasSrv aliasent.AliasBackend) *BaseManager {
	return &BaseManager{
		Aliases: aliasSrv,
	}
}

// Start does nothing as the DB client is already connected.
func (m *BaseManager) Start(_ context.Context) error { return nil }

// Stop does nothing as the DB client should be stopped outside of it.
func (m *BaseManager) Stop(_ context.Context) error { return nil }
func (m *BaseManager) Error() error                 { return nil }
func (m *BaseManager) Close() error                 { return nil }