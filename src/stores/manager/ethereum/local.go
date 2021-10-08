package eth

import (
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	mkeys "github.com/consensys/quorum-key-manager/src/stores/manager/keys"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

type LocalEthSpecs struct {
	Keystore manifest.Kind
	Specs    interface{}
}

func NewLocalEth(localETHSpecs *LocalEthSpecs, db database.Secrets, logger log.Logger) (stores.KeyStore, error) {
	var keyStore stores.KeyStore
	var err error

	switch localETHSpecs.Keystore {
	case manifest.HashicorpKeys:
		spec := &entities.HashicorpSpecs{}
		if err = json.UnmarshalJSON(localETHSpecs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp keystore localETHSpecs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewHashicorpKeyStore(spec, logger)
	case manifest.AKVKeys:
		spec := &mkeys.AkvKeySpecs{}
		if err = json.UnmarshalJSON(localETHSpecs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV keystore localETHSpecs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewAkvKeyStore(spec, logger)
	case manifest.AWSKeys:
		spec := &mkeys.AwsKeySpecs{}
		if err = json.UnmarshalJSON(localETHSpecs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS keystore localETHSpecs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewAwsKeyStore(spec, logger)
	case manifest.LocalKeys:
		spec := &mkeys.LocalKeySpecs{}
		if err = json.UnmarshalJSON(localETHSpecs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal local keystore localETHSpecs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		keyStore, err = mkeys.NewLocalKeyStore(spec, db, logger)
	default:
		errMessage := "invalid keystore kind"
		logger.Error(errMessage, "kind", localETHSpecs.Keystore)
		return nil, errors.InvalidFormatError(errMessage)
	}
	if err != nil {
		return nil, err
	}

	return keyStore, nil
}
