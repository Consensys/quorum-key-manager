package basemanager

import (
	"context"

	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/core/manifest/loader"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts"
	auditedaccounts "github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts/audit"
	authenticatedaccounts "github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts/auth"
	baseaccounts "github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts/base"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/keys"
	akvkeys "github.com/ConsenSysQuorum/quorum-key-manager/core/store/keys/azure-key-vault"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/secrets"
)

// AKVKeysSpecs is the specs format for an Azure Key Vault key store
type AKVKeysSpecs struct {
	AKV           *akvkeys.Config `json:"akv"`
	Audited       bool            `json:"audited"`
	Authenticated bool            `json:"authenticated"`
}

func (mngr *Manager) BuildAKVKeyStores(specs *AKVKeysSpecs) (secrets.Store, keys.Store, accounts.Store, error) {
	// Creates AKV keys store from specs config
	keysStore, err := akvkeys.New(specs.AKV)
	if err != nil {
		return nil, nil, nil, err
	}

	// Mount key store into an account store
	accountsStore := baseaccounts.NewStore(keysStore)

	// Instrument account store with authentication capabilities
	if specs.Authenticated {
		accountsStore = authenticatedaccounts.NewInstrument().Apply(accountsStore)
	}

	// Instrument account store with auditing capabilities
	if specs.Audited {
		accountsStore = auditedaccounts.NewInstrument(mngr.auditor).Apply(accountsStore)
	}

	// TODO: returning nil there is concerning, probably
	// we should probably return a NotCompatibleSecretStore (that always return NotCompatibleError)

	return nil, keysStore, accountsStore
}

// loadAKVKeys creates and indexes an AKV Key Store
func (mngr *Manager) loadAKVKeys(ctx context.Context, msg *manifestloader.Message) {
	// Unmarshal manifest specs
	specs := new(AKVKeysSpecs)
	msg.UnmarshalSpecs(specs)
	if msg.Err != nil {
		return
	}

	secretsStore, keysStore, accountsStore, err := mngr.BuildAKVKeyStores(specs)
	if err != nil {
		msg.Err = nil
		return
	}

	// TODO: if the store is common.Runnable, it should be started now

	// setStores on manager for later access
	mngr.setStores(msg, secretsStore, keysStore, accountsStore)
}
