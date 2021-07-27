package aliasservice

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

// RegisterService creates and register the alias service.
func RegisterService(a *app.App, pgClient postgres.Client) error {
	store := aliasstore.New(pgClient)
	m := NewManager(store)
	err := a.RegisterService(m)
	if err != nil {
		return err
	}

	return nil
}

type BaseManager struct {
	aliases.Backend
}

func NewManager(backend aliases.Backend) *BaseManager {
	return &BaseManager{
		Backend: backend,
	}
}

// Start does nothing as the DB is already initialized.
func (s *BaseManager) Start(ctx context.Context) error { return nil }

// Stop does nothing as the DB is already initialized.
func (s *BaseManager) Stop(context.Context) error { return nil }
func (s *BaseManager) Close() error               { return nil }
func (s *BaseManager) Error() error               { return nil }
