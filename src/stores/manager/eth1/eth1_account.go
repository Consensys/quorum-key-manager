package eth1

import (
	"context"

	"github.com/consensysquorum/quorum-key-manager/pkg/log"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/database"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys"

	manifest "github.com/consensysquorum/quorum-key-manager/src/manifests/types"
	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/akv"
	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/hashicorp"
	eth1 "github.com/consensysquorum/quorum-key-manager/src/stores/store/eth1/local"
	"github.com/consensysquorum/quorum-key-manager/src/stores/types"
)

type Specs struct {
	Keystore manifest.Kind
	Specs    interface{}
}

func NewEth1(ctx context.Context, specs *Specs, eth1Accounts database.ETH1Accounts, logger log.Logger) (*eth1.Store, error) {
	var keyStore keys.Store
	var err error

	switch specs.Keystore {
	case types.HashicorpKeys:
		spec := &hashicorp.KeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		keyStore, err = hashicorp.NewKeyStore(spec, logger)
	case types.AKVKeys:
		spec := &akv.KeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		keyStore, err = akv.NewKeyStore(spec, logger)
	default:
		errMessage := "invalid keystore kind"
		logger.Error(errMessage, "kind", specs.Keystore)
		return nil, errors.InvalidFormatError(errMessage)
	}
	if err != nil {
		return nil, err
	}

	err = InitDB(ctx, keyStore, eth1Accounts)
	if err != nil {
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
