package aliasservice

import (
	"context"

	"github.com/go-pg/pg/v10"

	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/aliases"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
)

func RegisterService(a *app.App, logger log.Logger, db database.Database) error {
	// Create and register the stores service

	// TODO replace by the database.Database abstraction
	pgdb := &pg.DB{}
	// TODO replace by the database.Database abstraction
	store := aliasstore.New(pgdb)
	m := NewManager(store)
	err := a.RegisterService(m)
	if err != nil {
		return err
	}

	return nil
}

type BaseManager struct {
	store aliases.API
}

func NewManager(store *aliasstore.Store) *BaseManager {
	return &BaseManager{
		store: store,
	}
}

// Start does nothing as the DB is already initialized.
func (s *BaseManager) Start(ctx context.Context) error { return nil }

// Stop does nothing as the DB is already initialized.
func (s *BaseManager) Stop(context.Context) error { return nil }
func (s *BaseManager) Close() error               { return nil }
func (s *BaseManager) Error() error               { return nil }
