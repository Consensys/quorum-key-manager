package aliasservice

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

func RegisterService(a *app.App, logger log.Logger, pgClient postgres.Client) error {
	// Create and register the stores service

	// TODO replace by the database.Database abstraction
	// TODO replace by the database.Database abstraction
	store := aliasstore.New(pgClient)
	m := NewManager(store)
	err := a.RegisterService(m)
	if err != nil {
		return err
	}

	return nil
}

type BaseManager struct {
	aliases.API
}

func NewManager(backend aliases.API) *BaseManager {
	return &BaseManager{
		API: backend,
	}
}

// Start does nothing as the DB is already initialized.
func (s *BaseManager) Start(ctx context.Context) error { return nil }

// Stop does nothing as the DB is already initialized.
func (s *BaseManager) Stop(context.Context) error { return nil }
func (s *BaseManager) Close() error               { return nil }
func (s *BaseManager) Error() error               { return nil }
