package eth1

import (
	"context"
	"fmt"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/database"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/akv"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/hashicorp"
	eth1 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1/local"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/types"
)

type Specs struct {
	Keystore manifest.Kind
	Specs    interface{}
}

func NewEth1(ctx context.Context, specs *Specs, eth1Accounts database.ETH1Accounts, logger *log.Logger) (*eth1.Store, error) {
	var keyStore keys.Store
	var err error

	switch specs.Keystore {
	case types.HashicorpKeys:
		spec := &hashicorp.KeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			logger.WithError(err).Error("failed to unmarshal Hashicorp keystore specs")
			return nil, err
		}
		keyStore, err = hashicorp.NewKeyStore(spec, logger)

		// We sleep to give some time for the token to be set
		// TODO: this needs to be improved to not rely on time
		time.Sleep(time.Second)
	case types.AKVKeys:
		spec := &akv.KeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			logger.WithError(err).Error("failed to unmarshal AKV keystore specs")
			return nil, err
		}
		keyStore, err = akv.NewKeyStore(spec, logger)
	default:
		err = fmt.Errorf("invalid keystore kind %s", specs.Keystore)
		logger.WithError(err).Error()
		return nil, err
	}
	if err != nil {
		logger.WithError(err).Error("failed to create Keystore")
		return nil, err
	}

	err = InitDB(ctx, keyStore, eth1Accounts)
	if err != nil {
		logger.WithError(err).Error("failed to initialize Eth1 store database")
		return nil, err
	}

	return eth1.New(keyStore, eth1Accounts, logger), nil
}

func InitDB(ctx context.Context, keyStore keys.Store, db database.ETH1Accounts) error {
	ids, err := keyStore.List(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {
		key, err := keyStore.Get(ctx, id)
		if err != nil {
			return err
		}

		if key.IsETH1Account() {
			err = db.Add(ctx, eth1.ParseKey(key))
			if err != nil && errors.IsAlreadyExistsError(err) {
				continue
			}
			if err != nil {
				return err
			}
		}
	}

	/* TODO: Uncomment when implemented in all stores
	deletedIDs, err := keyStore.ListDeleted(ctx)
	if err != nil {
		return err
	}

	for _, id := range deletedIDs {
		key, err := keyStore.GetDeleted(ctx, id)
		if err != nil {
			return err
		}

		if key.IsETH1Account() {
			err = db.AddDeleted(ctx, eth1.ParseKey(key))
			if err != nil {
				return err
			}
		}
	}
	*/

	return nil
}
