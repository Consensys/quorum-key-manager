package keys

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/aws/client"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/aws"
)

func NewAwsKeyStore(specs *entities.AwsSpecs, logger log.Logger) (*aws.Store, error) {
	cfg := client.NewConfig(specs.Region, specs.AccessID, specs.SecretKey, specs.Debug)
	cli, err := client.NewKmsClient(cfg)
	if err != nil {
		errMessage := "failed to instantiate AWS client (keys)"
		logger.WithError(err).Error(errMessage, "specs", specs)
		return nil, errors.ConfigError(errMessage)
	}

	store := aws.New(cli, logger)
	return store, nil
}
