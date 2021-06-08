package storemanager

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
	manifest2 "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/types"
	accounts2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/accounts"
	akv2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/akv"
	aws2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/aws"
	hashicorp2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/hashicorp"
	memory2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/database/memory"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	eth12 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1"
	keys2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys"
	secrets2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/secrets"
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/types"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type BaseManager struct {
	manifests manager.Manager

	mux          sync.RWMutex
	secrets      map[string]*storeBundle
	keys         map[string]*storeBundle
	eth1Accounts map[string]*storeBundle

	sub    manager.Subscription
	mnfsts chan []manager.Message
}

type storeBundle struct {
	manifest *manifest2.Manifest

	store interface{}
}

func New(manifests manager.Manager) *BaseManager {
	return &BaseManager{
		manifests:    manifests,
		mux:          sync.RWMutex{},
		secrets:      make(map[string]*storeBundle),
		keys:         make(map[string]*storeBundle),
		eth1Accounts: make(map[string]*storeBundle),
		mnfsts:       make(chan []manager.Message),
	}
}

var storeKinds = []manifest2.Kind{
	types2.HashicorpSecrets,
	types2.AKVSecrets,
	types2.KMSSecrets,
	types2.AKVKeys,
	types2.HashicorpKeys,
	types2.KMSKeys,
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

func (m *BaseManager) GetSecretStore(_ context.Context, name string) (secrets2.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if storeBundle, ok := m.secrets[name]; ok {
		if store, ok := storeBundle.store.(secrets2.Store); ok {
			return store, nil
		}
	}

	return nil, errors.NotFoundError("secret store %s was not found", name)
}

func (m *BaseManager) GetKeyStore(_ context.Context, name string) (keys2.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if storeBundle, ok := m.keys[name]; ok {
		if store, ok := storeBundle.store.(keys2.Store); ok {
			return store, nil
		}
	}

	return nil, errors.NotFoundError("key store %s was not found", name)
}

func (m *BaseManager) GetEth1Store(ctx context.Context, name string) (eth12.Store, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.getEth1Store(ctx, name)
}

func (m *BaseManager) getEth1Store(_ context.Context, name string) (eth12.Store, error) {
	if storeBundle, ok := m.eth1Accounts[name]; ok {
		if store, ok := storeBundle.store.(eth12.Store); ok {
			return store, nil
		}
	}

	return nil, errors.NotFoundError("account store %s was not found", name)
}

func (m *BaseManager) GetEth1StoreByAddr(ctx context.Context, addr ethcommon.Address) (eth12.Store, error) {
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

func (m *BaseManager) List(ctx context.Context, kind manifest2.Kind) ([]string, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return m.list(ctx, kind)
}

func (m *BaseManager) list(_ context.Context, kind manifest2.Kind) ([]string, error) {
	storeNames := []string{}
	switch kind {
	case "":
		storeNames = append(
			append(m.storeNames(m.secrets, kind), m.storeNames(m.keys, kind)...), m.storeNames(m.eth1Accounts, kind)...)
	case types2.HashicorpSecrets, types2.AKVSecrets, types2.KMSSecrets:
		storeNames = m.storeNames(m.secrets, kind)
	case types2.AKVKeys, types2.HashicorpKeys, types2.KMSKeys:
		storeNames = m.storeNames(m.keys, kind)
	}

	return storeNames, nil
}

func (m *BaseManager) ListAllAccounts(ctx context.Context) ([]*entities2.ETH1Account, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	accs := []*entities2.ETH1Account{}
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

func (m *BaseManager) load(ctx context.Context, mnf *manifest2.Manifest) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	logger := log.FromContext(ctx).
		WithField("kind", mnf.Kind).
		WithField("name", mnf.Name)

	logger.Debug("loading store manifest")
	errMsg := "error creating new store store"

	switch mnf.Kind {
	case types2.HashicorpSecrets:
		spec := &hashicorp2.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := hashicorp2.NewSecretStore(spec, logger)
		if err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types2.HashicorpKeys:
		spec := &hashicorp2.KeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := hashicorp2.NewKeyStore(spec, logger)
		if err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}
		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types2.AKVSecrets:
		spec := &akv2.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := akv2.NewSecretStore(spec, logger)
		if err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types2.AKVKeys:
		spec := &akv2.KeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := akv2.NewKeyStore(spec, logger)
		if err != nil {
			return err
		}
		m.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types2.AWSSecrets:
		spec := &aws2.SecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return errors.InvalidFormatError(err.Error())
		}
		store, err := aws2.NewSecretStore(spec, logger)
		if err != nil {
			return err
		}
		m.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store}
	case types2.Eth1Account:
		spec := &accounts2.Eth1Specs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			logger.WithError(err).Error(errMsg)
			return err
		}

		memdb := memory2.New(logger)
		store, err := accounts2.NewEth1(spec, memdb, logger)
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

func (m *BaseManager) storeNames(list map[string]*storeBundle, kind manifest2.Kind) []string {
	var storeNames []string
	for k, store := range list {
		if kind == "" || store.manifest.Kind == kind {
			storeNames = append(storeNames, k)
		}
	}

	return storeNames
}
