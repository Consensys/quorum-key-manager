package client

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault/keyvaultapi"
)

type AKVClient struct {
	client keyvaultapi.BaseClientAPI
	cfg    *Config
}

func NewClient(cfg *Config) (*AKVClient, error) {
	client := keyvault.New()

	authorizer, err := cfg.ToAzureAuthConfig()
	if err != nil {
		return nil, err
	}
	client.Authorizer = authorizer

	return &AKVClient{client: client, cfg: cfg}, nil
}
