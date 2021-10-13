package tests

import (
	"encoding/json"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"os"

	pgclient "github.com/consensys/quorum-key-manager/src/infra/postgres/client"
)

const envVar = "TEST_DATA"

type Config struct {
	AkvClient           *akvClient       `json:"akv_client"`
	AwsClient           *awsClient       `json:"aws_client"`
	KeyManagerURL       string           `json:"key_manager_url"`
	HealthKeyManagerURL string           `json:"health_key_manager_url"`
	SecretStores        []string         `json:"secret_stores"`
	KeyStores           []string         `json:"key_stores"`
	EthStores           []string         `json:"eth_stores"`
	QuorumNodeID        string           `json:"quorum_node_id"`
	BesuNodeID          string           `json:"besu_node_id"`
	GethNodeID          string           `json:"geth_node_id"`
	Postgres            *pgclient.Config `json:"postgres"`
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

func (c *Config) AkvSpecs() *entities.AkvSpecs {
	return &entities.AkvSpecs{
		ClientID:     c.AkvClient.ClientID,
		TenantID:     c.AkvClient.TenantID,
		VaultName:    c.AkvClient.VaultName,
		ClientSecret: c.AkvClient.ClientSecret,
	}
}

func (c *Config) AwsSpecs() *entities.AwsSpecs {
	return &entities.AwsSpecs{
		Region:    c.AwsClient.Region,
		AccessID:  c.AwsClient.AccessID,
		SecretKey: c.AwsClient.SecretKey,
	}
}
