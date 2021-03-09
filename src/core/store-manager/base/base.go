package basemanager

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/audit"
	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest/loader"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
)

type Manager struct {
	mux    sync.RWLock
	stores map[string]*storeBundle

	auditor audit.Auditor
}

type storeBundle struct {
	msg      *manifestloader.Message
	secrets  secrets.Store
	keys     keys.Store
	accounts accounts.Store
}

func New(auditor audit.Auditor) *Manager {
	return &Manager{
		secrets:  make(map[string]secrets.Store),
		keys:     make(map[string]keys.Store),
		accounts: make(map[string]accounts.Store),
		auditor:  auditor,
	}
}

func (mngr *Manager) Load(ctx context.Context, msgs ...*manifestloader.Message) {
	for _, msg := range msgs {
		mngr.loadMessage(ctx, msg)
	}
}

func (mngr *Manager) loadMessage(ctx context.Context, msg *manifestloader.Message) {
	switch msg.Manifest.Kind {
	case "HashicorpSecrets":
		mngr.loadHashicorpSecrets(ctx, msg)
	case "AKVKeys":
		mngr.loadAKVKeys(ctx, msg)
	default:
		msg.Err = fmt.Errorf("invalid manifest Kind")
	}
}

func (mngr *Manager) setStores(msg *manifestloader.Message, secretsStore secrets.Store, keysStore keys.Store, accountsStore accounts.Store) {
	mngr.mux.Lock()
	mngr.stores[name] = &storeBundle{
		msg:      msg,
		secrets:  secretsStore,
		keys:     keysStore,
		accounts: accountsStore,
	}
	mngr.mux.Unlock()
}

func (mngr *Manager) GetSecretStore(ctx context.Context, name string) (secrets.Store, error) {
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

func (mngr *Manager) GetKeyStore(ctx context.Context, name string) (keys.Store, error) {
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

func (mngr *Manager) GetAccountStore(ctx context.Context, name string) (accounts.Store, error) {
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
