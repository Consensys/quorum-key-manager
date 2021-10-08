package imports

import (
	"context"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/consensys/quorum-key-manager/src/stores/manager/secrets"
)

func ImportSecrets(ctx context.Context, db database.Secrets, mnf *manifest.Manifest, logger log.Logger) error {
	startTime := time.Now()

	store, err := getStore(mnf, logger)
	if err != nil {
		return err
	}

	ids, err := store.List(ctx, 0, 0)
	if err != nil {
		return err
	}

	for _, id := range ids {
		secret, err := store.Get(ctx, id, "")
		if err != nil {
			return err
		}

		// If secret already exists in DB, we skip. This allows idempotency of the import script (run it multiple times)
		// No need to treat the error here
		dbSecret, _ := db.Get(ctx, id, secret.Metadata.Version)
		if dbSecret != nil {
			logger.Debug("secret already exists, skipping", "id", id)
			continue
		}

		_, err = db.Add(ctx, secret)
		if err != nil {
			return err
		}
	}

	logger.Info("secrets successfully imported", "n", len(ids), "start_time", startTime, "end_time", time.Now())
	return nil
}

func getStore(mnf *manifest.Manifest, logger log.Logger) (stores.SecretStore, error) {
	switch mnf.Kind {
	case manifest.HashicorpSecrets:
		spec := &entities.HashicorpSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "invalid Hashicorp secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		logger.Info("importing secrets from hashicorp vault...")
		return secrets.NewHashicorpSecretStore(spec, nil, logger) // DB here is nil and not the DB we instantiate for the import
	case manifest.AKVSecrets:
		spec := &secrets.AkvSecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "invalid AKV secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		logger.Info("importing secrets from AKV...")
		return secrets.NewAkvSecretStore(spec, logger)
	case manifest.AWSSecrets:
		spec := &secrets.AwsSecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "invalid AWS secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		logger.Info("importing secrets from AWS...")
		return secrets.NewAwsSecretStore(spec, logger)
	}

	errMessage := "invalid manifest kind for secret store"
	logger.Error(errMessage)
	return nil, errors.InvalidFormatError(errMessage)
}
