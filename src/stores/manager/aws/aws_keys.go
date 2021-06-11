package aws

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	client2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/aws/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/aws"
)

// KeysSpecs is the specs format for an aws secrets manager (aws secretsmanager service)
type KeysSpecs struct {
	Region    string `json:"region"`
	AccessID  string `json:"accessID"`
	SecretKey string `json:"secretKey"`
}

func NewKeyStore(specs *KeysSpecs, logger *log.Logger) (*aws.KeyStore, error) {
	cfg := client2.NewBaseConfig(specs.Region, specs.AccessID, specs.SecretKey)
	cli, err := client2.NewKmsClient(cfg)
	if err != nil {
		return nil, err
	}

	store := aws.New(cli, logger)
	return store, nil
}
