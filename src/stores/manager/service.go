package storemanager

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

const ID = "StoreManager"

type BaseManager struct {
	isLive bool
	err    error
	db     database.Database
	logger log.Logger
}

func New(db database.Database, logger log.Logger) *BaseManager {
	return &BaseManager{
		logger: logger,
		db:     db,
	}
}

func (m *BaseManager) Start(_ context.Context) error {
	m.isLive = true
	return nil
}

func (m *BaseManager) Stop(_ context.Context) error {
	m.isLive = false
	return nil
}

func (m *BaseManager) Error() error { return m.err }
func (m *BaseManager) Close() error { return nil }
func (m *BaseManager) ID() string   { return ID }

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
