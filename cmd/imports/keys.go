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
	eth "github.com/consensys/quorum-key-manager/src/stores/manager/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/manager/keys"
)

func ImportKeys(ctx context.Context, db database.Keys, mnf *manifest.Manifest, logger log.Logger) error {
	logger.Info("importing keys...", "store", mnf.Kind, "store_name", mnf.Name)

	store, err := getKeyStore(mnf, logger)
	if err != nil {
		return err
	}

	storeIDs, err := store.List(ctx, 0, 0)
	if err != nil {
		return err
	}

	dbIDs, err := db.SearchIDs(ctx, false, 0, 0)
	if err != nil {
		return err
	}

	var n uint
	for _, id := range difference(storeIDs, dbIDs) {
		secret, err := store.Get(ctx, id)
		if err != nil {
			return err
		}

		_, err = db.Add(ctx, secret)
		if err != nil {
			return err
		}

		n++
	}

	logger.Info("keys imported successfully", "n", n)
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

		return keys.NewHashicorpKeyStore(spec, logger)
	case manifest.AKVKeys:
		spec := &entities.AkvSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "invalid AKV key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		return keys.NewAkvKeyStore(spec, logger)
	case manifest.AWSKeys:
		spec := &entities.AwsSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "invalid AWS key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		return keys.NewAwsKeyStore(spec, logger)
	case manifest.Ethereum:
		spec := &entities.LocalEthSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "invalid ethereum store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		return eth.NewLocalEth(spec, nil, logger) // DB here is nil and not the DB we instantiate for the import
	}

	errMessage := "invalid manifest kind for key store"
	logger.Error(errMessage)
	return nil, errors.InvalidFormatError(errMessage)
}
