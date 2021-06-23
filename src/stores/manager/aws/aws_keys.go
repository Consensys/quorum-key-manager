package aws

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/aws/client"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/keys/aws"
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
		return nil, err
	}

	store := aws.New(cli, *logger)
	return store, nil
}
