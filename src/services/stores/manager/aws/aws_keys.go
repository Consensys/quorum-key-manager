package aws

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/aws/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/aws"
)

// KeysSpecs is the specs format for an aws secrets manager (aws secretsmanager service)
type KeysSpecs struct {
	Region    string `json:"region"`
	AccessID  string `json:"accessID"`
	SecretKey string `json:"secretKey"`
}

func NewKeyStore(specs *KeysSpecs, logger *log.Logger) (*aws.KeyStore, error) {
	cfg := client.NewBaseConfig(specs.Region, specs.AccessID, specs.SecretKey)
	cli, err := client.NewKmsClient(cfg)
	if err != nil {
		return nil, errors.AWSConnectionError(err.Error())
	}

	store := aws.New(cli, logger)
	return store, nil
}
