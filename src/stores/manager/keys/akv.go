package keys

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/akv/client"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/akv"
)

func NewAkvKeyStore(spec *entities.AkvSpecs, logger log.Logger) (*akv.Store, error) {
	cfg := client.NewConfig(spec.VaultName, spec.TenantID, spec.ClientID, spec.ClientSecret)
	cli, err := client.NewClient(cfg)
	if err != nil {
		errMessage := "failed to instantiate AKV client (keys)"
		logger.WithError(err).Error(errMessage, "specs", spec)
		return nil, errors.ConfigError(errMessage)
	}

	store := akv.New(cli, logger)
	return store, nil
}
