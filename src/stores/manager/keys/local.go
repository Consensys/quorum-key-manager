package keys

import (
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	localkeys "github.com/consensys/quorum-key-manager/src/stores/store/keys/local"

	msecrets "github.com/consensys/quorum-key-manager/src/stores/manager/secrets"
)

func NewLocalKeyStore(localKeyStoreSpecs *entities.LocalKeySpecs, db database.Secrets, logger log.Logger) (*localkeys.Store, error) {
	var secretStore stores.SecretStore
	var err error

	switch localKeyStoreSpecs.SecretStore {
	case manifest.HashicorpSecrets:
		spec := &entities.HashicorpSpecs{}
		if err = json.UnmarshalJSON(localKeyStoreSpecs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp secret store localKeyStoreSpecs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		secretStore, err = msecrets.NewHashicorpSecretStore(spec, db, logger)
	case manifest.AKVSecrets:
		spec := &entities.AkvSpecs{}
		if err = json.UnmarshalJSON(localKeyStoreSpecs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV secret store localKeyStoreSpecs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		secretStore, err = msecrets.NewAkvSecretStore(spec, logger)
	case manifest.AWSSecrets:
		spec := &entities.AwsSpecs{}
		if err = json.UnmarshalJSON(localKeyStoreSpecs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS secret store localKeyStoreSpecs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		secretStore, err = msecrets.NewAwsSecretStore(spec, logger)
	default:
		errMessage := "invalid secret store kind"
		logger.Error(errMessage, "kind", localKeyStoreSpecs.SecretStore)
		return nil, errors.InvalidFormatError(errMessage)
	}
	if err != nil {
		return nil, err
	}

	return localkeys.New(secretStore, db, logger), nil
}
