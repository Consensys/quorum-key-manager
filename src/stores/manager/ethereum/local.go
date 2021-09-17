package eth

import (
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	mkeys "github.com/consensys/quorum-key-manager/src/stores/manager/keys"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
)

type LocalEthSpecs struct {
	Keystore manifest.Kind
	Specs    interface{}
}

func NewLocalEth(specs *LocalEthSpecs, db database.Secrets, logger log.Logger) (stores.KeyStore, error) {
	var keyStore stores.KeyStore
	var err error

	switch specs.Keystore {
	case manifest.HashicorpKeys:
		spec := &mkeys.HashicorpKeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewHashicorpKeyStore(spec, logger)
	case manifest.AKVKeys:
		spec := &mkeys.AkvKeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewAkvKeyStore(spec, logger)
	case manifest.AWSKeys:
		spec := &mkeys.AwsKeySpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS keystore specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewAwsKeyStore(spec, logger)
	case manifest.LocalKeys:
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
