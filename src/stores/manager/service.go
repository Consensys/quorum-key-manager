package storemanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensys/quorum-key-manager/src/auth/policy"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/connectors"

	"github.com/consensys/quorum-key-manager/src/stores/store/database"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
	"github.com/consensys/quorum-key-manager/src/stores/manager/akv"
	"github.com/consensys/quorum-key-manager/src/stores/manager/aws"
	"github.com/consensys/quorum-key-manager/src/stores/manager/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores/manager/local"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	eth1store "github.com/consensys/quorum-key-manager/src/stores/store/eth1"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"
	"github.com/consensys/quorum-key-manager/src/stores/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const ID = "StoreManager"

type BaseManager struct {
	manifests     manifestsmanager.Manager
	policyManager policy.Manager

	mux          sync.RWMutex
	secrets      map[string]*storeBundle
	keys         map[string]*storeBundle
	eth1Accounts map[string]*storeBundle

	sub    manifestsmanager.Subscription
	mnfsts chan []manifestsmanager.Message

	isLive bool

	logger log.Logger
	db     database.Database
}

type storeBundle struct {
	manifest *manifest.Manifest
	logger   log.Logger
	store    interface{}
}

func New(manifests manifestsmanager.Manager, policyManager policy.Manager, db database.Database, logger log.Logger) *BaseManager {
	return &BaseManager{
		manifests:     manifests,
		policyManager: policyManager,
		mux:           sync.RWMutex{},
		secrets:       make(map[string]*storeBundle),
		keys:          make(map[string]*storeBundle),
		eth1Accounts:  make(map[string]*storeBundle),
		mnfsts:        make(chan []manifestsmanager.Message),
		logger:        logger,
		db:            db,
	}
}

var storeKinds = []manifest.Kind{
	types.HashicorpSecrets,
	types.HashicorpKeys,
	types.AKVSecrets,
	types.AKVKeys,
	types.AWSSecrets,
	types.AWSKeys,
	types.Eth1Account,
}

func (m *BaseManager) Start(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	defer func() {
		m.isLive = true
	}()

	// Subscribe to manifest of Kind node
	m.sub = m.manifests.Subscribe(storeKinds, m.mnfsts)

	// Start loading manifest
	go m.loadAll(ctx)

	return nil
}

func (m *BaseManager) Stop(context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.isLive = false

	if m.sub != nil {
		_ = m.sub.Unsubscribe()
	}
	close(m.mnfsts)
	return nil
}

func (m *BaseManager) Error() error {
	return nil
}

func (m *BaseManager) Close() error {
	return nil
}

func (m *BaseManager) loadAll(ctx context.Context) {
	for mnfsts := range m.mnfsts {
		for _, mnf := range mnfsts {
			_ = m.load(ctx, mnf.Manifest)
		}
	}
}

func (m *BaseManager) GetSecretStore(ctx context.Context, name string, userInfo *authtypes.UserInfo) (secrets.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if storeBundle, ok := m.secrets[name]; ok {
		if store, ok := storeBundle.store.(secrets.Store); ok {
			policies := m.policyManager.UserPolicies(ctx, userInfo)
			resolvr, err := policy.NewRadixResolver(policies...)
			if err != nil {
				return nil, err
			}
			return connectors.NewSecretConnector(store, resolvr, storeBundle.logger), nil
		}
	}

	errMessage := "secret store was not found"
	m.logger.Error(errMessage, "store_name", name)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetKeyStore(ctx context.Context, name string, userInfo *authtypes.UserInfo) (keys.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if storeBundle, ok := m.keys[name]; ok {
		if store, ok := storeBundle.store.(keys.Store); ok {
			policies := m.policyManager.UserPolicies(ctx, userInfo)
			resolvr, err := policy.NewRadixResolver(policies...)
			if err != nil {
				return nil, err
			}
			return connectors.NewKeyConnector(store, resolvr, storeBundle.logger), nil
		}
	}

	errMessage := "key store was not found"
	m.logger.Error(errMessage, "store_name", name)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetEth1Store(ctx context.Context, name string, userInfo *authtypes.UserInfo) (eth1store.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.getEth1Store(ctx, name, userInfo)
}

func (m *BaseManager) getEth1Store(ctx context.Context, name string, userInfo *authtypes.UserInfo) (eth1store.Store, error) {
	if storeBundle, ok := m.eth1Accounts[name]; ok {
		if store, ok := storeBundle.store.(eth1store.Store); ok {
			policies := m.policyManager.UserPolicies(ctx, userInfo)
			resolvr, err := policy.NewRadixResolver(policies...)
			if err != nil {
				return nil, err
			}
			return connectors.NewEth1Connector(store, resolvr, storeBundle.logger), nil
		}
	}

	errMessage := "account store was not found"
	m.logger.Error(errMessage, "store_name", name)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetEth1StoreByAddr(ctx context.Context, addr ethcommon.Address, userInfo *authtypes.UserInfo) (eth1store.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	for _, storeName := range m.list(ctx, "") {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			acc, err := m.getEth1Store(ctx, storeName, userInfo)
			if err == nil {
				// Check if account exists in store and returns it
				_, err := acc.Get(ctx, addr.Hex())
				if err == nil {
					return acc, nil
				}
			}
		}
	}

	errMessage := "account was not found"
	m.logger.Error(errMessage, "account", addr.Hex())
	return nil, errors.InvalidParameterError(errMessage)
}

func (m *BaseManager) List(ctx context.Context, kind manifest.Kind, _ *authtypes.UserInfo) ([]string, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	// TODO Filter available store using userInfo groups
	return m.list(ctx, kind), nil
}

func (m *BaseManager) list(_ context.Context, kind manifest.Kind) []string {
	storeNames := []string{}
	switch kind {
	case "":
		storeNames = append(
			append(m.storeNames(m.secrets, kind), m.storeNames(m.keys, kind)...), m.storeNames(m.eth1Accounts, kind)...)
	case types.HashicorpSecrets, types.AKVSecrets, types.AWSSecrets:
		storeNames = m.storeNames(m.secrets, kind)
	case types.AKVKeys, types.HashicorpKeys, types.AWSKeys:
		storeNames = m.storeNames(m.keys, kind)
	case types.Eth1Account:
		storeNames = m.storeNames(m.eth1Accounts, kind)
	}

	return storeNames
}

func (m *BaseManager) ListAllAccounts(ctx context.Context, userInfo *authtypes.UserInfo) ([]*entities.ETH1Account, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	accs := []*entities.ETH1Account{}
	for _, storeName := range m.list(ctx, "") {
		store, err := m.getEth1Store(ctx, storeName, userInfo)
		if err == nil {
			storeAccs, err := store.GetAll(ctx)
			if err == nil {
				accs = append(accs, storeAccs...)
			}
		}
	}

	return accs, nil
}

func (m *BaseManager) load(ctx context.Context, mnf *manifest.Manifest) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	logger := m.logger.With("kind", mnf.Kind).With("name", mnf.Name)
	logger.Debug("loading store manifest")

	switch mnf.Kind {
	case types.HashicorpSecrets:
		spec := &hashicorp.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := hashicorp.NewSecretStore(spec, logger)
		if err != nil {
			return err
		}

		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case types.HashicorpKeys:
		spec := &hashicorp.KeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := hashicorp.NewKeyStore(spec, logger)
		if err != nil {
			return err
		}

		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case types.AKVSecrets:
		spec := &akv.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AKV secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := akv.NewSecretStore(spec, logger)
		if err != nil {
			return err
		}

		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case types.AKVKeys:
		spec := &akv.KeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AKV key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := akv.NewKeyStore(spec, logger)
		if err != nil {
			return err
		}

		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case types.AWSSecrets:
		spec := &aws.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AWS secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := aws.NewSecretStore(spec, logger)
		if err != nil {
			return err
		}

		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case types.AWSKeys:
		spec := &aws.KeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AWS key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := aws.NewKeyStore(spec, logger)
		if err != nil {
			return err
		}

		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case types.Eth1Account:
		spec := &local.Eth1Specs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Eth1 store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := local.NewEth1(ctx, spec, m.db.ETH1Accounts(), logger)
		if err != nil {
			return err
		}

		m.eth1Accounts[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	default:
		errMessage := "invalid manifest kind"
		logger.Error(errMessage, "kind", mnf.Kind)
		return errors.InvalidFormatError(errMessage)
	}

	logger.Info("store manifest loaded successfully")
	return nil
}

func (m *BaseManager) storeNames(list map[string]*storeBundle, kind manifest.Kind) []string {
	var storeNames []string
	for k, store := range list {
		if kind == "" || store.manifest.Kind == kind {
			storeNames = append(storeNames, k)
		}
	}

	return storeNames
}

func (m *BaseManager) ID() string { return ID }
func (m *BaseManager) CheckLiveness() error {
	if m.isLive {
		return nil
	}

	errMessage := fmt.Sprintf("service %s is not live", m.ID())
	m.logger.Error(errMessage, "id", m.ID())
	return errors.HealthcheckError(errMessage)
}

func (m *BaseManager) CheckReadiness() error {
	return m.Error()
}
