package basemanager

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	manifestloader "github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest/loader"
	hashicorp "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
)

// HashicorpSecretSpecs is the specs format for an Hashicorp Vault secret store
type HashicorpSecretSpecs struct {
	Hashicorp     *hashicorp.Config `json:"hashicorp"`
	Audited       bool              `json:"audited"`
	Authenticated bool              `json:"authenticated"`
}

func (mngr *manager) BuildHashicorpSecretStores(specs *HashicorpSecretSpecs) (secrets.Store, keys.Store, accounts.Store, error) {
	return nil, nil, nil, errors.NotImplementedError
	// Creates Hasicorp secrets store from specs config
	// secretsStore, err := secrets.New(specs.Hashicorp)
	// if err != nil {
	// 	return nil, nil, nil, err
	// }
	//
	// // Mount secret store into Key Store
	// keysStore := localkeys.New(secretsStore)
	//
	// // Mount key store into an account store
	// accountsStore := baseaccounts.NewStore(keysStore)
	//
	// // Instrument account store with authentication capabilities
	// if specs.Authenticated {
	// 	accountsStore = authenticatedaccounts.NewInstrument().Apply(accountsStore)
	// }
	//
	// // Instrument account store with auditing capabilities
	// if specs.Audited {
	// 	accountsStore = auditedaccounts.NewInstrument(mngr.auditor).Apply(accountsStore)
	// }
	//
	// return secretsStore, keysStore, accountsStore, nil
}

// loadHashicorpSecrets creates and indexes an Hashicorp secrets store
func (mngr *manager) loadHashicorpSecrets(ctx context.Context, msg *manifestloader.Message) {
	// Unmarshal manifest specs
	specs := new(HashicorpSecretSpecs)
	msg.UnmarshalSpecs(specs)
	if msg.Err != nil {
		return
	}

	secretsStore, keysStore, accountsStore, err := mngr.BuildHashicorpSecretStores(specs)
	if err != nil {
		msg.Err = nil
		return
	}

	// TODO: if the store is common.Runnable, it should be started now

	// setStores on manager for later access
	mngr.setStores(msg, secretsStore, keysStore, accountsStore)
}
