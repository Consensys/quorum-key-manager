package eth1

import (
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	mkeys "github.com/consensys/quorum-key-manager/src/stores/manager/keys"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
)

type LocalEth1Specs struct {
	Keystore manifest.Kind
	Specs    interface{}
}

func NewLocalEth1(specs *LocalEth1Specs, db database.Secrets, logger log.Logger) (stores.KeyStore, error) {
	var keyStore stores.KeyStore
	var err error

	switch specs.Keystore {
	case stores.HashicorpKeys:
		spec := &mkeys.HashicorpKeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewHashicorpKeyStore(spec, logger)
	case stores.AKVKeys:
		spec := &mkeys.AkvKeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewAkvKeyStore(spec, logger)
	case stores.AWSKeys:
		spec := &mkeys.AwsKeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewAwsKeyStore(spec, logger)
	case stores.LocalKeys:
		spec := &mkeys.LocalKeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal local keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewLocalKeyStore(spec, db, logger)
	default:
		errMessage := "invalid keystore kind"
		logger.Error(errMessage, "kind", specs.Keystore)
		return nil, errors.InvalidFormatError(errMessage)
	}
	if err != nil {
		return nil, err
	}

	return keyStore, nil
}
