//nolint
package akv

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	akvkeys "github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/azure-key-vault"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"
)

// AKVKeysSpecs is the specs format for an Azure Key Vault key store
type akvKeysSpecs struct {
	AKV           *akvkeys.Config `json:"akv"`
}

func BuildAKVKeyStores(specs *akvKeysSpecs) (secrets.Store, keys.Store, accounts.Store, error) {
	return nil, nil, nil, errors.NotImplementedError
	// // Creates AKV keys store from specs config
	// keysStore, err := akvkeys.New(specs.AKV)
	// if err != nil {
	// 	return nil, nil, nil, err
	// }
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
	// // TODO: returning nil there is concerning, probably
	// // we should probably return a NotCompatibleSecretStore (that always return NotCompatibleError)
	//
	// return nil, keysStore, accountsStore
}

// loadAKVKeys creates and indexes an AKV Key Store
// func (mngr *manager) loadAKVKeys(ctx context.Context, msg *manifestloader.Message) {
// 	// Unmarshal manifest specs
// 	specs := new(akvKeysSpecs)
// 	msg.UnmarshalSpecs(specs)
// 	if msg.Err != nil {
// 		return
// 	}
//
// 	secretsStore, keysStore, accountsStore, err := mngr.BuildAKVKeyStores(specs)
// 	if err != nil {
// 		msg.Err = nil
// 		return
// 	}
//
// 	// TODO: if the store is common.Runnable, it should be started now
//
// 	// setStores on manager for later access
// 	mngr.setStores(msg, secretsStore, keysStore, accountsStore)
// }
