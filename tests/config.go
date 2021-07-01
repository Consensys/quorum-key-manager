package tests

import (
	"encoding/json"
	"fmt"
	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/aws"
	"os"

	"github.com/consensysquorum/quorum-key-manager/src/stores/manager/akv"
)

const envVar = "TEST_DATA"

type Config struct {
	AkvClient            *akvClient `json:"akv_client"`
	AwsClient            *awsClient `json:"aws_client"`
	KeyManagerURL        string     `json:"key_manager_url"`
	HealthKeyManagerURL  string     `json:"health_key_manager_url"`
	HashicorpSecretStore string     `json:"hashicorp_secret_store"`
	HashicorpKeyStore    string     `json:"hashicorp_key_store"`
	Eth1Store            string     `json:"eth1_store"`
	QuorumNodeID         string     `json:"quorum_node_id"`
	BesuNodeID           string     `json:"besu_node_id"`
}

type akvClient struct {
	VaultName    string `json:"vault_name"`
	TenantID     string `json:"tenant_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type awsClient struct {
	AccessID  string `json:"access_id"`
	Region    string `json:"region"`
	SecretKey string `json:"secret_key"`
}

func NewConfig() (*Config, error) {
	cfgStr := os.Getenv(envVar)
	if cfgStr == "" {
		return nil, fmt.Errorf("expected test data at environment variable '%s'", envVar)
	}

	cfg := &Config{}
	if err := json.Unmarshal([]byte(cfgStr), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) AkvSecretSpecs() *akv.SecretSpecs {
	return &akv.SecretSpecs{
		ClientID:     c.AkvClient.ClientID,
		TenantID:     c.AkvClient.TenantID,
		VaultName:    c.AkvClient.VaultName,
		ClientSecret: c.AkvClient.ClientSecret,
	}
}

func (c *Config) AkvKeySpecs() *akv.KeySpecs {
	return &akv.KeySpecs{
		ClientID:     c.AkvClient.ClientID,
		TenantID:     c.AkvClient.TenantID,
		VaultName:    c.AkvClient.VaultName,
		ClientSecret: c.AkvClient.ClientSecret,
	}
}

func (c *Config) AwsSecretSpecs() *aws.SecretSpecs {
	return &aws.SecretSpecs{
		Region:    c.AwsClient.Region,
		AccessID:  c.AwsClient.AccessID,
		SecretKey: c.AwsClient.SecretKey,
	}
}

func (c *Config) AwsKeySpecs() *aws.KeySpecs {
	return &aws.KeySpecs{
		Region:    c.AwsClient.Region,
		AccessID:  c.AwsClient.AccessID,
		SecretKey: c.AwsClient.SecretKey,
	}
}
