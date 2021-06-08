package aws

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	client2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/aws/client"
	aws2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/secrets/aws"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
)

// SecretSpecs is the specs format for an aws secrets manager (aws secretsmanager service)
type SecretSpecs struct {
	Region    string `json:"region"`
	AccessID  string `json:"accessID"`
	SecretKey string `json:"secretKey"`
}

func NewSecretStore(specs *SecretSpecs, logger *log.Logger) (*aws2.SecretStore, error) {
	cfg := client2.NewBaseConfig(specs.Region, specs.AccessID, specs.SecretKey)
	cli, err := client2.NewSecretsClient(cfg)
	if err != nil {
		return nil, errors.AWSConnectionError(err.Error())
	}

	store := aws2.New(cli, logger)
	return store, nil
}
