package basemanager

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/audit"
	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/core/manifest/loader"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/secrets"
)

type Manager struct {
	mux sync.RWLock
	secrets  map[string]secrets.Store
	keys     map[string]keys.Store
	accounts map[string]accounts.Store

	auditor audit.Auditor
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

func (mngr *Manager) setStores(name string, secretsStore secrets.Store, keysStore keys.Store, accountsStore accounts.Store) {
	mngr.mux.Lock()
	mngr.secrets[name] = secretsStore
	mngr.keys[name] = keysStore
	mngr.accounts[name] = accountsStore
	mngr.mux.Unlock()
}

func (mngr *Manager) GetSecretStore(ctx context.Context, name string) (secrets.Store, error) {
	mngr.mux.RLock()
	defer mngr.mux.RUnlock()

	s, ok := mngr.secrets[name]
	if !ok {
		return nil, fmt.Errorf("store not found")
	}
	
	return s, nil
}

func (mngr *Manager) GetKeyStore(ctx context.Context, name string) (keys.Store, error) {
	mngr.mux.RLock()
	defer mngr.mux.RUnlock()
	
	s, ok := mngr.keys[name]
	if !ok {
		return nil, fmt.Errorf("store not found")
	}

	return s, nil
}

func (mngr *Manager) GetAccountStore(ctx context.Context, name string) (accounts.Store, error) {
	mngr.mux.RLock()
	defer mngr.mux.RUnlock()
	
	s, ok := mngr.accounts[name]
	if !ok {
		return nil, fmt.Errorf("store not found")
	}

	return s, nil
}
