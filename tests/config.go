package tests

import (
	"encoding/json"
	"fmt"
	"os"
)

const TestDataEnv = "TEST_DATA"

type Config struct {
	KeyManagerURL        string   `json:"key_manager_url"`
	HealthKeyManagerURL  string   `json:"health_key_manager_url"`
	SecretStores         []string `json:"secret_stores"`
	KeyStores            []string `json:"key_stores"`
	EthStores            []string `json:"eth_stores"`
	QuorumNodeID         string   `json:"quorum_node_id"`
	BesuNodeID           string   `json:"besu_node_id"`
	GethNodeID           string   `json:"geth_node_id"`
	ApiKey               string   `json:"api_key"`
	AuthTLSKey           string   `json:"tls_key"`
	AuthTLSCert          string   `json:"tls_cert"`
	AuthOIDCTokenURL     string   `json:"oidc_token_url"`
	AuthOIDCClientID     string   `json:"oidc_client_id"`
	AuthOIDCClientSecret string   `json:"oidc_client_secret"`
	AuthOIDCAudience     string   `json:"oidc_audience"`
}

func NewConfig() (*Config, error) {
	cfgStr := os.Getenv(TestDataEnv)
	if cfgStr == "" {
		return nil, fmt.Errorf("expected test data at environment variable '%s'", TestDataEnv)
	}

	cfg := &Config{}
	if err := json.Unmarshal([]byte(cfgStr), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
