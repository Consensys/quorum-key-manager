package aws

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/aws/client"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/secrets/aws"
)

// SecretSpecs is the specs format for an aws secrets manager (aws secretsmanager service)
type SecretSpecs struct {
	Region    string `json:"region"`
	AccessID  string `json:"accessID"`
	SecretKey string `json:"secretKey"`
	Debug     bool   `json:"debug"`
}

func NewSecretStore(specs *SecretSpecs, logger log.Logger) (*aws.SecretStore, error) {
	cfg := client.NewConfig(specs.Region, specs.AccessID, specs.SecretKey, specs.Debug)
	cli, err := client.NewSecretsClient(cfg)
	if err != nil {
		errMessage := "failed to instantiate AWS client (secrets)"
		logger.WithError(err).Error(errMessage, "specs", specs)
		return nil, errors.ConfigError(errMessage)
	}

	store := aws.New(cli, logger)
	return store, nil
}
