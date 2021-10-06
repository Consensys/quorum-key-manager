package storemanager

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

const ID = "StoreManager"

type BaseManager struct {
	manifests manifestsmanager.Manager

	isLive bool
	err    error

	db     database.Database
	logger log.Logger

	stores stores.Stores
}

var _ stores.Manager = &BaseManager{}

func New(storesConnector stores.Stores, manifests manifestsmanager.Manager, db database.Database, logger log.Logger) *BaseManager {
	return &BaseManager{
		manifests: manifests,
		logger:    logger,
		db:        db,
		stores:    storesConnector,
	}
}

func (m *BaseManager) Start(ctx context.Context) error {
	messages, err := m.manifests.Load()
	if err != nil {
		return err
	}

	for _, message := range messages {
		err = m.stores.Create(ctx, message.Manifest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *BaseManager) Stop(context.Context) error {
	m.isLive = false
	return nil
}

func (m *BaseManager) Error() error {
	return m.err
}

func (m *BaseManager) Close() error {
	return nil
}

func (m *BaseManager) Stores() stores.Stores {
	return m.stores
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
		return err
	}

	return nil
}
