package local

import (
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/stores/manager/aws"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys"

	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
	"github.com/consensys/quorum-key-manager/src/stores/manager/akv"
	"github.com/consensys/quorum-key-manager/src/stores/manager/hashicorp"
	eth1 "github.com/consensys/quorum-key-manager/src/stores/store/eth1/local"
	"github.com/consensys/quorum-key-manager/src/stores/types"
)

type Eth1Specs struct {
	Keystore manifest.Kind
	Specs    interface{}
}

func NewEth1(specs *Eth1Specs, db database.Database, logger log.Logger) (*eth1.Store, error) {
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
	case types.AWSKeys:
		spec := &aws.KeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		keyStore, err = aws.NewKeyStore(spec, logger)
	default:
		errMessage := "invalid keystore kind"
		logger.Error(errMessage, "kind", specs.Keystore)
		return nil, errors.InvalidFormatError(errMessage)
	}
	if err != nil {
		return nil, err
	}

	return eth1.New(keyStore, db, logger), nil
}
