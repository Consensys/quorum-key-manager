package local

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"

	"github.com/consensys/quorum-key-manager/src/stores/manager/aws"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
	localkeys "github.com/consensys/quorum-key-manager/src/stores/store/keys/local"

	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
	"github.com/consensys/quorum-key-manager/src/stores/manager/akv"
	"github.com/consensys/quorum-key-manager/src/stores/manager/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores/types"
)

type KeySpecs struct {
	SecretStore manifest.Kind
	Specs       interface{}
}

func NewLocalKeys(_ context.Context, specs *KeySpecs, db database.Database, logger log.Logger) (*localkeys.Store, error) {
	var secretStore secrets.Store
	var err error

	switch specs.SecretStore {
	case types.HashicorpSecrets:
		spec := &hashicorp.SecretSpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		secretStore, err = hashicorp.NewSecretStore(spec, logger)
	case types.AKVSecrets:
		spec := &akv.SecretSpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		secretStore, err = akv.NewSecretStore(spec, logger)
	case types.AWSSecrets:
		spec := &aws.SecretSpecs{}
		if err = manifest.UnmarshalSpecs(specs.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		secretStore, err = aws.NewSecretStore(spec, logger)
	default:
		errMessage := "invalid secret store kind"
		logger.Error(errMessage, "kind", specs.SecretStore)
		return nil, errors.InvalidFormatError(errMessage)
	}
	if err != nil {
		return nil, err
	}

	return localkeys.New(secretStore, db, logger), nil
}
