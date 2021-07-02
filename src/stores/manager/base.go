package storemanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensysquorum/quorum-key-manager/pkg/log"

	"github.com/consensysquorum/quorum-key-manager/src/stores/store/database/memory"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	manifestsmanager "github.com/consensysquorum/quorum-key-manager/src/manifests/manager"
	manifest "github.com/consensysquorum/quorum-key-manager/src/manifests/types"
	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/akv"
	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/aws"
	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/eth1"
	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/hashicorp"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	eth1store "github.com/consensysquorum/quorum-key-manager/src/stores/store/eth1"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/secrets"
	"github.com/consensysquorum/quorum-key-manager/src/stores/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const ID = "StoreManager"

type BaseManager struct {
	manifests manifestsmanager.Manager

	mux          sync.RWMutex
	secrets      map[string]*storeBundle
	keys         map[string]*storeBundle
	eth1Accounts map[string]*storeBundle

	sub    manifestsmanager.Subscription
	mnfsts chan []manifestsmanager.Message

	isLive bool

	logger log.Logger
}

type storeBundle struct {
	manifest *manifest.Manifest
	store    interface{}
}

func New(manifests manifestsmanager.Manager, logger log.Logger) *BaseManager {
	return &BaseManager{
		manifests:    manifests,
		mux:          sync.RWMutex{},
		secrets:      make(map[string]*storeBundle),
		keys:         make(map[string]*storeBundle),
		eth1Accounts: make(map[string]*storeBundle),
		mnfsts:       make(chan []manifestsmanager.Message),
		logger:       logger,
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

func (m *BaseManager) GetSecretStore(_ context.Context, name string) (secrets.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if storeBundle, ok := m.secrets[name]; ok {
		if store, ok := storeBundle.store.(secrets.Store); ok {
			return store, nil
		}
	}

	errMessage := "secret store was not found"
	m.logger.Error(errMessage, "store_name", name)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetKeyStore(_ context.Context, name string) (keys.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if storeBundle, ok := m.keys[name]; ok {
		if store, ok := storeBundle.store.(keys.Store); ok {
			return store, nil
		}
	}

	errMessage := "key store was not found"
	m.logger.Error(errMessage, "store_name", name)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetEth1Store(ctx context.Context, name string) (eth1store.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.getEth1Store(ctx, name)
}

func (m *BaseManager) getEth1Store(_ context.Context, name string) (eth1store.Store, error) {
	if storeBundle, ok := m.eth1Accounts[name]; ok {
		if store, ok := storeBundle.store.(eth1store.Store); ok {
			return store, nil
		}
	}

	errMessage := "account store was not found"
	m.logger.Error(errMessage, "store_name", name)
	return nil, errors.NotFoundError(errMessage)
}

func (m *BaseManager) GetEth1StoreByAddr(ctx context.Context, addr ethcommon.Address) (eth1store.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	for _, storeName := range m.list(ctx, "") {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			acc, err := m.getEth1Store(ctx, storeName)
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

func (m *BaseManager) List(ctx context.Context, kind manifest.Kind) ([]string, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

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

func (m *BaseManager) ListAllAccounts(ctx context.Context) ([]*entities.ETH1Account, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	accs := []*entities.ETH1Account{}
	for _, storeName := range m.list(ctx, "") {
		store, err := m.getEth1Store(ctx, storeName)
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
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
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
		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store}
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
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
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
		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store}
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
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
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
		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types.Eth1Account:
		spec := &eth1.Specs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Eth1 store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := eth1.NewEth1(ctx, spec, memory.New(logger), logger)
		if err != nil {
			return err
		}
		m.eth1Accounts[mnf.Name] = &storeBundle{manifest: mnf, store: store}
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
