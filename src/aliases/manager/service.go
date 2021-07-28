package aliasservice

import (
	"context"

	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
)

type BaseManager struct {
	aliasstore.Store
}

func NewManager(store aliasstore.Store) *BaseManager {
	return &BaseManager{
		Store: store,
	}
}

// Start does nothing as the DB is already initialized.
func (s *BaseManager) Start(ctx context.Context) error { return nil }

// Stop does nothing as the DB is already initialized.
func (s *BaseManager) Stop(context.Context) error { return nil }
func (s *BaseManager) Close() error               { return nil }
func (s *BaseManager) Error() error               { return nil }
