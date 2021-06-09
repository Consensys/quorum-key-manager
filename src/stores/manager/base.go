package storemanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/database/memory"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	manifestsmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/akv"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/aws"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/eth1"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	eth1store "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/secrets"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type BaseManager struct {
	manifests manifestsmanager.Manager

	mux          sync.RWMutex
	secrets      map[string]*storeBundle
	keys         map[string]*storeBundle
	eth1Accounts map[string]*storeBundle

	sub    manifestsmanager.Subscription
	mnfsts chan []manifestsmanager.Message
}

type storeBundle struct {
	manifest *manifest.Manifest

	store interface{}
}

func New(manifests manifestsmanager.Manager) *BaseManager {
	return &BaseManager{
		manifests:    manifests,
		mux:          sync.RWMutex{},
		secrets:      make(map[string]*storeBundle),
		keys:         make(map[string]*storeBundle),
		eth1Accounts: make(map[string]*storeBundle),
		mnfsts:       make(chan []manifestsmanager.Message),
	}
}

var storeKinds = []manifest.Kind{
	types.HashicorpSecrets,
	types.AKVSecrets,
	types.KMSSecrets,
	types.AKVKeys,
	types.HashicorpKeys,
	types.KMSKeys,
	types.Eth1Account,
}

func (m *BaseManager) Start(ctx context.Context) error {
	m.mux.Lock()
	// Subscribe to manifest of Kind node
	sub, err := m.manifests.Subscribe(storeKinds, m.mnfsts)
	if err != nil {
		return err
	}
	m.sub = sub
	m.mux.Unlock()

	// Start loading manifest
	go m.loadAll(ctx)

	return nil
}

func (m *BaseManager) Stop(context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()
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

	return nil, errors.NotFoundError("secret store %s was not found", name)
}

func (m *BaseManager) GetKeyStore(_ context.Context, name string) (keys.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if storeBundle, ok := m.keys[name]; ok {
		if store, ok := storeBundle.store.(keys.Store); ok {
			return store, nil
		}
	}

	return nil, errors.NotFoundError("key store %s was not found", name)
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

	return nil, errors.NotFoundError("account store %s was not found", name)
}

func (m *BaseManager) GetEth1StoreByAddr(ctx context.Context, addr ethcommon.Address) (eth1store.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	storeNames, err := m.list(ctx, "")
	if err != nil {
		return nil, err
	}

	for _, storeName := range storeNames {
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

	return nil, errors.InvalidParameterError("account %s was not found", addr.String())
}

func (m *BaseManager) List(ctx context.Context, kind manifest.Kind) ([]string, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.list(ctx, kind)
}

func (m *BaseManager) list(_ context.Context, kind manifest.Kind) ([]string, error) {
	storeNames := []string{}
	switch kind {
	case "":
		storeNames = append(
			append(m.storeNames(m.secrets, kind), m.storeNames(m.keys, kind)...), m.storeNames(m.eth1Accounts, kind)...)
	case types.HashicorpSecrets, types.AKVSecrets, types.KMSSecrets:
		storeNames = m.storeNames(m.secrets, kind)
	case types.AKVKeys, types.HashicorpKeys, types.KMSKeys:
		storeNames = m.storeNames(m.keys, kind)
	}

	return storeNames, nil
}

func (m *BaseManager) ListAllAccounts(ctx context.Context) ([]*entities.ETH1Account, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	accs := []*entities.ETH1Account{}
	storeNames, err := m.list(ctx, "")
	if err != nil {
		return accs, err
	}

	for _, storeName := range storeNames {
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

	logger := log.FromContext(ctx).
		WithField("kind", mnf.Kind).
		WithField("name", mnf.Name)

	logger.Debug("loading store manifest")
	errMsg := "error creating new store store"

	switch mnf.Kind {
	case types.HashicorpSecrets:
		spec := &hashicorp.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := hashicorp.NewSecretStore(spec, logger)
		if err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types.HashicorpKeys:
		spec := &hashicorp.KeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := hashicorp.NewKeyStore(spec, logger)
		if err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}
		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types.AKVSecrets:
		spec := &akv.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := akv.NewSecretStore(spec, logger)
		if err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types.AKVKeys:
		spec := &akv.KeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := akv.NewKeyStore(spec, logger)
		if err != nil {
			return err
		}
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types.AWSSecrets:
		spec := &aws.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := aws.NewSecretStore(spec, logger)
		if err != nil {
			return err
		}
		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types.Eth1Account:
		spec := &eth1.Specs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}

		store, err := eth1.NewEth1(ctx, spec, memory.New(logger), logger)
		if err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}
		m.eth1Accounts[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	default:
		err := fmt.Errorf("invalid manifest kind %s", mnf.Kind)
		logger.WithError(err).Error()
		return err
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
