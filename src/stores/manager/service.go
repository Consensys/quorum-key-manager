package storemanager

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/src/auth"
	storesconnector "github.com/consensys/quorum-key-manager/src/stores/connectors/stores"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/utils"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

const ID = "StoreManager"

type BaseManager struct {
	manifests manifestsmanager.Manager

	sub    manifestsmanager.Subscription
	mnfsts chan []manifestsmanager.Message

	isLive bool
	err    error

	db     database.Database
	logger log.Logger

	utils  stores.Utilities
	stores stores.Stores
}

var _ stores.Manager = &BaseManager{}

func New(manifests manifestsmanager.Manager, authManager auth.Manager, db database.Database, logger log.Logger) *BaseManager {
	return &BaseManager{
		manifests: manifests,
		mnfsts:    make(chan []manifestsmanager.Message),
		logger:    logger,
		db:        db,
		utils:     utils.NewConnector(logger),
		stores:    storesconnector.NewConnector(authManager, db, logger),
	}
}

func (m *BaseManager) Start(ctx context.Context) error {
	defer func() {
		m.isLive = true
	}()

	// Subscribe to manifest of Kind node
	m.sub = m.manifests.Subscribe(manifest.StoreKinds, m.mnfsts)

	// Start loading manifest
	go m.loadAll(ctx)

	return nil
}

func (m *BaseManager) Stop(context.Context) error {
	m.isLive = false

	if m.sub != nil {
		_ = m.sub.Unsubscribe()
	}
	close(m.mnfsts)
	return nil
}

func (m *BaseManager) Error() error {
	return m.err
}

func (m *BaseManager) Close() error {
	return nil
}

func (m *BaseManager) Utilities() stores.Utilities {
	return m.utils
}

func (m *BaseManager) Stores() stores.Stores {
	return m.stores
}

func (m *BaseManager) loadAll(ctx context.Context) {
	for mnfsts := range m.mnfsts {
		for _, mnf := range mnfsts {
			if err := m.stores.Create(ctx, mnf.Manifest); err != nil {
				m.err = errors.CombineErrors(m.err, err)
			}
		}
	}
}

func (m *BaseManager) ID() string { return ID }

func (m *BaseManager) CheckLiveness(_ context.Context) error {
	if m.isLive {
		return nil
	}

	errMessage := fmt.Sprintf("service %s is not live", m.ID())
	m.logger.Error(errMessage, "id", m.ID())
	return errors.HealthcheckError(errMessage)
}

func (m *BaseManager) CheckReadiness(ctx context.Context) error {
	err := m.Error()
	if err != nil {
		return err
	}

	err = m.db.Ping(ctx)
	if err != nil {
		return errors.DependencyFailureError("database connection error: %s", err.Error())
	}

	return nil
}
