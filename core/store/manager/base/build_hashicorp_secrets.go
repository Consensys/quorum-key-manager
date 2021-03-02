package basemanager

import (
	"context"

	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/core/manifest/loader"
	auditedaccounts "github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts/audit"
	defaultaccounts "github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts/default"
	localkeys "github.com/ConsenSysQuorum/quorum-key-manager/core/store/keys/local"
	hashicorpsecrets "github.com/ConsenSysQuorum/quorum-key-manager/core/store/secrets/hashicorp"
)

// HashicorpSecretSpecs is the specs format for an Hashicorp Vault secret store
type HashicorpSecretSpecs struct {
	Hashicorp *hashicorpsecrets.Config `json:"hasicorp"`
	Audited   bool                     `json:"audited"`
}

func BuildHashicorpSecretStores(specs *HashicorpSecretSpecs) (secrets.Store, keys.Store, accounts.Store, error) {
	// Creates Hasicorp secrets store from specs config
	secretsStore, err := hashicorpsecrets.New(specs.Hashicorp)
	if err != nil {
		return nil, nil, nil, err
	}

	// Mount secret store into Key Store
	keysStore := localkeys.New(secretsStore)

	// Mount key store into an account store
	accountsStore := defaultaccounts.NewStore(keysStore)

	// Wraps account store with auditing capabilities
	if specs.Audited {
		accountsStore = auditedaccounts.Wrap(accountsStore)
	}
	
	return secretsStore, keysStore, accountsStore, nil
}

// loadHashicorpSecrets creates and indexes an Hashicorp secrets store
func (mngr *Manager) loadHashicorpSecrets(ctx context.Context, msg *manifestloader.Message) {
	// Unmarshal manifest specs
	specs := new(HashicorpSecretSpecs)
	msg.UnmarshalSpecs(specs)
	if msg.Err != nil {
		return
	}

	secretsStore, keysStore, accountsStore, err := BuildHashicorpSecretStores(specs)
	if err != nil {
		msg.Err = nil
		return
	}
	
	// TODO: if the store is common.Runnable, it should be started now

	// setStores on manager for later access
	mngr.setStores(msg, secretsStore, keysStore, accountsStore)
}
