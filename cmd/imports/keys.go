package imports

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/consensys/quorum-key-manager/src/stores/manager/keys"
)

func ImportKeys(ctx context.Context, db database.Keys, mnf *manifest.Manifest, logger log.Logger) error {
	store, err := getKeyStore(mnf, logger)
	if err != nil {
		return err
	}

	ids, err := store.List(ctx, 0, 0)
	if err != nil {
		return err
	}

	for _, id := range ids {
		secret, err := store.Get(ctx, id)
		if err != nil {
			return err
		}

		// If key already exists in DB, we skip. This allows idempotency of the import script (run it multiple times)
		// No need to treat the error here
		dbKey, _ := db.Get(ctx, id)
		if dbKey != nil {
			logger.Debug("key already exists, skipping", "id", id)
			continue
		}

		_, err = db.Add(ctx, secret)
		if err != nil {
			return err
		}
	}

	logger.Info("keys successfully imported", "n", len(ids))
	return nil
}

func getKeyStore(mnf *manifest.Manifest, logger log.Logger) (stores.KeyStore, error) {
	switch mnf.Kind {
	case manifest.HashicorpKeys:
		spec := &entities.HashicorpSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "invalid Hashicorp key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		logger.Info("importing keys from hashicorp vault...")
		return keys.NewHashicorpKeyStore(spec, logger)
	case manifest.AKVKeys:
		spec := &entities.AkvSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "invalid AKV key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		logger.Info("importing keys from AKV...")
		return keys.NewAkvKeyStore(spec, logger)
	case manifest.AWSKeys:
		spec := &entities.AwsSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "invalid AWS key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		logger.Info("importing keys from AWS KMS...")
		return keys.NewAwsKeyStore(spec, logger)
	}

	errMessage := "invalid manifest kind for key store"
	logger.Error(errMessage)
	return nil, errors.InvalidFormatError(errMessage)
}
