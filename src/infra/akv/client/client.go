package client

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault/keyvaultapi"
	"github.com/consensys/quorum-key-manager/src/infra/akv"
)

type AKVClient struct {
	client keyvaultapi.BaseClientAPI
	cfg    *Config
}

var _ akv.Client = &AKVClient{}

func NewClient(cfg *Config) (*AKVClient, error) {
	client := keyvault.New()

	authorizer, err := cfg.ToAzureAuthConfig()
	if err != nil {
		return nil, err
	}
	client.Authorizer = authorizer

	return &AKVClient{client: client, cfg: cfg}, nil
}
