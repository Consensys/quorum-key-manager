package aws

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/aws/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/aws"
)

// SecretSpecs is the specs format for an aws secrets manager (aws secretsmanager service)
type SecretSpecs struct {
	Region    string `json:"region"`
	AccessID  string `json:"accessID"`
	SecretKey string `json:"secretKey"`
}

func NewSecretStore(specs *SecretSpecs) (*aws.SecretStore, error) {
	cfg := client.NewBaseConfig(specs.Region, specs.AccessID, specs.SecretKey)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, errors.AWSConnectionError(err.Error())
	}

	store := aws.New(cli)
	return store, nil
}
