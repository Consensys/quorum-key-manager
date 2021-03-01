package basemanager

import (
	"context"

	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/core/manifest/loader"
	auditedaccounts "github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts/audit"
	defaultaccounts "github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts/default"
	akvkeys "github.com/ConsenSysQuorum/quorum-key-manager/core/store/keys/azure-key-vault"
)

// AKVKeysSpecs is the specs format for an Azure Key Vault key store
type AKVKeysSpecs struct {
	AKV     *akvkeys.Config `json:"akv"`
	Audited bool            `json:"audited"`
}

// loadAKVKeys creates and indexes an AKV Key Store
func (mngr *Manager) loadAKVKeys(ctx context.Context, msg *manifestloader.Message) {
	// Unmarshal manifest specs
	specs := new(AKVKeysSpecs)
	msg.UnmarshalSpecs(specs)
	if msg.Err != nil {
		return
	}

	// Creates AKV keys store from specs config
	keysStore, err := akvkeys.New(specs.AKV)
	if err != nil {
		msg.Err = err
		return
	}

	// Mount key store into an account store
	accountsStore := defaultaccounts.NewStore(keysStore)

	// Wraps account store with auditing capabilities
	if specs.Audited {
		accountsStore = auditedaccounts.Wrap(accountsStore)
	}

	// TODO: if the store is common.Runnable, it should be started now

	// setStores on manager for later access
	mngr.setStores(msg.Manifest.Name, nil, keysStore, accountsStore)
}
