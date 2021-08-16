package storemanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensys/quorum-key-manager/src/stores"
	eth1connector "github.com/consensys/quorum-key-manager/src/stores/connectors/eth1"
	keysconnector "github.com/consensys/quorum-key-manager/src/stores/connectors/keys"
	secretsconnector "github.com/consensys/quorum-key-manager/src/stores/connectors/secrets"

	"github.com/consensys/quorum-key-manager/src/auth/policy"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/database"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
	meth1 "github.com/consensys/quorum-key-manager/src/stores/manager/eth1"
	mkeys "github.com/consensys/quorum-key-manager/src/stores/manager/keys"
	msecrets "github.com/consensys/quorum-key-manager/src/stores/manager/secrets"
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

var _ stores.Manager = &BaseManager{}

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
	stores.HashicorpSecrets,
	stores.HashicorpKeys,
	stores.AKVSecrets,
	stores.AKVKeys,
	stores.AWSSecrets,
	stores.AWSKeys,
	stores.LocalKeys,
	stores.Eth1Account,
}

func (m *BaseManager) Start(_ context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	defer func() {
		m.isLive = true
	}()

	// Subscribe to manifest of Kind node
	m.sub = m.manifests.Subscribe(storeKinds, m.mnfsts)

	// Start loading manifest
	go m.loadAll()

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

func (m *BaseManager) loadAll() {
	for mnfsts := range m.mnfsts {
		for _, mnf := range mnfsts {
			_ = m.load(mnf.Manifest)
		}
	}
}

func (m *BaseManager) GetSecretStore(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.SecretStore, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if storeBundle, ok := m.secrets[storeName]; ok {
		if store, ok := storeBundle.store.(stores.SecretStore); ok {
			policies := m.policyManager.UserPolicies(ctx, userInfo)
			_, _ = policy.NewResolver(policies)
			return secretsconnector.NewConnector(store, m.db.Secrets(storeName), storeBundle.logger), nil
		}
	}

	errMessage := "secret store was not found"
	m.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetKeyStore(_ context.Context, storeName string, _ *authtypes.UserInfo) (stores.KeyStore, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if storeBundle, ok := m.keys[storeName]; ok {
		if store, ok := storeBundle.store.(stores.KeyStore); ok {
			return keysconnector.NewConnector(store, m.db.Keys(storeName), storeBundle.logger), nil
		}
	}

	errMessage := "key store was not found"
	m.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetEth1Store(ctx context.Context, name string, userInfo *authtypes.UserInfo) (stores.Eth1Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.getEth1Store(ctx, name, userInfo)
}

func (m *BaseManager) getEth1Store(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.Eth1Store, error) {
	if storeBundle, ok := m.eth1Accounts[storeName]; ok {
		if store, ok := storeBundle.store.(stores.KeyStore); ok {
			policies := m.policyManager.UserPolicies(ctx, userInfo)
			_, _ = policy.NewResolver(policies)
			return eth1connector.NewConnector(store, m.db.ETH1Accounts(storeName), storeBundle.logger), nil
		}
	}

	errMessage := "account store was not found"
	m.logger.Error(errMessage, "store_name", storeName)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetEth1StoreByAddr(ctx context.Context, addr ethcommon.Address, userInfo *authtypes.UserInfo) (stores.Eth1Store, error) {
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
				_, err := acc.Get(ctx, addr)
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
	case stores.HashicorpSecrets, stores.AKVSecrets, stores.AWSSecrets:
		storeNames = m.storeNames(m.secrets, kind)
	case stores.AKVKeys, stores.HashicorpKeys, stores.AWSKeys:
		storeNames = m.storeNames(m.keys, kind)
	case stores.Eth1Account:
		storeNames = m.storeNames(m.eth1Accounts, kind)
	}

	return storeNames
}

func (m *BaseManager) ListAllAccounts(ctx context.Context, userInfo *authtypes.UserInfo) ([]ethcommon.Address, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	accs := []ethcommon.Address{}
	for _, storeName := range m.list(ctx, "") {
		store, err := m.getEth1Store(ctx, storeName, userInfo)
		if err == nil {
			storeAccs, err := store.List(ctx)
			if err == nil {
				accs = append(accs, storeAccs...)
			}
		}
	}

	return accs, nil
}

func (m *BaseManager) load(mnf *manifest.Manifest) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	logger := m.logger.With("kind", mnf.Kind).With("name", mnf.Name)
	logger.Debug("loading store manifest")

	switch mnf.Kind {
	case stores.HashicorpSecrets:
		spec := &msecrets.HashicorpSecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := msecrets.NewHashicorpSecretStore(spec, logger)
		if err != nil {
			return err
		}

		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case stores.HashicorpKeys:
		spec := &mkeys.HashicorpKeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := mkeys.NewHashicorpKeyStore(spec, logger)
		if err != nil {
			return err
		}

		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case stores.AKVSecrets:
		spec := &msecrets.AkvSecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AKV secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := msecrets.NewAkvSecretStore(spec, logger)
		if err != nil {
			return err
		}

		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case stores.AKVKeys:
		spec := &mkeys.AkvKeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AKV key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := mkeys.NewAkvKeyStore(spec, logger)
		if err != nil {
			return err
		}

		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case stores.AWSSecrets:
		spec := &msecrets.AwsSecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AWS secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := msecrets.NewAwsSecretStore(spec, logger)
		if err != nil {
			return err
		}

		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case stores.AWSKeys:
		spec := &mkeys.AwsKeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AWS key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := mkeys.NewAwsKeyStore(spec, logger)
		if err != nil {
			return err
		}

		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case stores.LocalKeys:
		spec := &mkeys.LocalKeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal local key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := mkeys.NewLocalKeyStore(spec, logger)
		if err != nil {
			return err
		}

		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case stores.Eth1Account:
		spec := &meth1.LocalEth1Specs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Eth1 store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := meth1.NewLocalEth1(spec, logger)
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
