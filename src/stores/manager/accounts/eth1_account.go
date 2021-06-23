package accounts

import (
	"fmt"
	manifest2 "github.com/consensysquorum/quorum-key-manager/src/manifests/types"
	akv2 "github.com/consensysquorum/quorum-key-manager/src/stores/manager/akv"
	hashicorp2 "github.com/consensysquorum/quorum-key-manager/src/stores/manager/hashicorp"
	database2 "github.com/consensysquorum/quorum-key-manager/src/stores/store/database"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/eth1/local"
	types2 "github.com/consensysquorum/quorum-key-manager/src/stores/types"

	"github.com/consensysquorum/quorum-key-manager/pkg/log"
)

type Eth1Specs struct {
	Keystore manifest2.Kind
	Specs    interface{}
}

func NewEth1(specs *Eth1Specs, eth1Accounts database2.ETH1Accounts, logger log.Logger) (*local.Store, error) {
	switch specs.Keystore {
	case types2.HashicorpKeys:
		spec := &hashicorp2.KeySpecs{}
		if err := manifest2.UnmarshalSpecs(specs.Specs, spec); err != nil {
			logger.WithError(err).Error("failed to unmarshal Hashicorp keystore specs")
			return nil, err
		}
		store, err := hashicorp2.NewKeyStore(spec, logger)
		if err != nil {
			logger.WithError(err).Error("failed to create new Hashicorp Keystore")
			return nil, err
		}
		return local.New(store, eth1Accounts, logger), nil
	case types2.AKVKeys:
		spec := &akv2.KeySpecs{}
		if err := manifest2.UnmarshalSpecs(specs.Specs, spec); err != nil {
			logger.WithError(err).Error("failed to unmarshal AKV keystore specs")
			return nil, err
		}
		store, err := akv2.NewKeyStore(spec, logger)
		if err != nil {
			logger.WithError(err).Error("failed to create new AKV Keystore")
			return nil, err
		}
		return local.New(store, eth1Accounts, logger), nil
	default:
		err := fmt.Errorf("invalid keystore kind %s", specs.Keystore)
		logger.WithError(err).Error("invalid keystore kind")
		return nil, err
	}
}
