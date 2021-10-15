package storemanager

import (
	"context"
	"fmt"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/manifests"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

const ID = "StoreManager"

type BaseManager struct {
	manifestReader manifests.Reader

	isLive bool
	err    error

	db     database.Database
	logger log.Logger

	stores stores.Stores
}

var _ stores.Manager = &BaseManager{}

func New(storesConnector stores.Stores, manifestReader manifests.Reader, db database.Database, logger log.Logger) *BaseManager {
	return &BaseManager{
		manifestReader: manifestReader,
		logger:         logger,
		db:             db,
		stores:         storesConnector,
	}
}

func (m *BaseManager) Start(ctx context.Context) error {
	mnfs, err := m.manifestReader.Load()
	if err != nil {
		errMessage := "failed to load manifest file"
		m.logger.WithError(err).Error(errMessage)
		return errors.ConfigError(errMessage)
	}

	for _, mnf := range mnfs {
		// TODO: Filter on Load() function from reader when Kind Store implemented
		if mnf.Kind == manifest.Role || mnf.Kind == manifest.Node {
			continue
		}

		storeType := manifest.StoreType(mnf.Kind)
		switch storeType {
		case manifest.HashicorpSecrets, manifest.AKVSecrets, manifest.AWSSecrets:
			err = m.stores.CreateSecret(ctx, mnf.Name, storeType, mnf.Specs, mnf.AllowedTenants)
		case manifest.HashicorpKeys, manifest.AKVKeys, manifest.AWSKeys, manifest.LocalKeys:
			err = m.stores.CreateKey(ctx, mnf.Name, storeType, mnf.Specs, mnf.AllowedTenants)
		case manifest.Ethereum:
			err = m.stores.CreateEthereum(ctx, mnf.Name, storeType, mnf.Specs, mnf.AllowedTenants)
		}
		if err != nil {
			return err
		}
	}

	m.isLive = true

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
