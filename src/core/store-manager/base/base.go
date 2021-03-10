package basemanager

import (
	"context"
	"fmt"
	"sync"

	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest/loader"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
)

type manager struct {
	mux    sync.RWMutex
	stores map[string]*storeBundle
}

type storeBundle struct {
	msg      *manifestloader.Message
	secrets  secrets.Store
	keys     keys.Store
	accounts accounts.Store
}

func New() storemanager.Manager {
	return &manager{}
}

func (mngr *manager) Load(ctx context.Context, msgs ...*manifestloader.Message) {
	for _, msg := range msgs {
		mngr.loadMessage(ctx, msg)
	}
}

func (mngr *manager) loadMessage(ctx context.Context, msg *manifestloader.Message) {
	switch msg.Manifest.Kind {
	case "HashicorpSecrets":
		mngr.loadHashicorpSecrets(ctx, msg)
	case "AKVKeys":
		mngr.loadAKVKeys(ctx, msg)
	default:
		msg.Err = fmt.Errorf("invalid manifest Kind")
	}
}

func (mngr *manager) setStores(msg *manifestloader.Message, secretsStore secrets.Store, keysStore keys.Store, accountsStore accounts.Store) {
	mngr.mux.Lock()
	// @TODO Set name for new store bundle
	mngr.stores["name"] = &storeBundle{
		msg:      msg,
		secrets:  secretsStore,
		keys:     keysStore,
		accounts: accountsStore,
	}
	mngr.mux.Unlock()
}

func (mngr *manager) GetSecretStore(ctx context.Context, name string) (secrets.Store, error) {
	mngr.mux.RLock()
	defer mngr.mux.RUnlock()

	s, ok := mngr.stores[name]
	if !ok {
		return nil, fmt.Errorf("store not found")
	}

	if s.msg.Err != nil {
		return nil, s.msg.Err
	}

	return s.secrets, nil
}

func (mngr *manager) GetKeyStore(ctx context.Context, name string) (keys.Store, error) {
	mngr.mux.RLock()
	defer mngr.mux.RUnlock()

	s, ok := mngr.stores[name]
	if !ok {
		return nil, fmt.Errorf("store not found")
	}

	if s.msg.Err != nil {
		return nil, s.msg.Err
	}

	return s.keys, nil
}

func (mngr *manager) GetAccountStore(ctx context.Context, name string) (accounts.Store, error) {
	mngr.mux.RLock()
	defer mngr.mux.RUnlock()

	s, ok := mngr.stores[name]
	if !ok {
		return nil, fmt.Errorf("store not found")
	}

	if s.msg.Err != nil {
		return nil, s.msg.Err
	}

	return s.accounts, nil
}

func (mngr *manager) List(ctx context.Context, kind string) ([]string, error) {
	panic("implement me")
}
