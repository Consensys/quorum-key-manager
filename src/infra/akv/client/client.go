package client

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault/keyvaultapi"
)

type AzureClient struct {
	client keyvaultapi.BaseClientAPI
	cfg    *Config
}

func NewClient(cfg *Config) (*AzureClient, error) {
	client := keyvault.New()

	authorizer, err := cfg.ToAzureAuthConfig()
	if err != nil {
		return nil, err
	}
	client.Authorizer = authorizer

	return &AzureClient{client: client, cfg: cfg}, nil
}

func (c *AzureClient) SetSecret(ctx context.Context, secretName string, parameters keyvault.SecretSetParameters) (result keyvault.SecretBundle, err error) {
	return c.client.SetSecret(ctx, c.cfg.Endpoint, secretName, parameters)
}

func (c *AzureClient) GetSecret(ctx context.Context, secretName, secretVersion string) (result keyvault.SecretBundle, err error) {
	return c.client.GetSecret(ctx, c.cfg.Endpoint, secretName, secretVersion)
}

func (c *AzureClient) GetSecrets(ctx context.Context, maxResults *int32) (result keyvault.SecretListResultPage, err error) {
	return c.client.GetSecrets(ctx, c.cfg.Endpoint, maxResults)
}

func (c *AzureClient) UpdateSecret(ctx context.Context, secretName, secretVersion string, parameters keyvault.SecretUpdateParameters) (result keyvault.SecretBundle, err error) {
	return c.client.UpdateSecret(ctx, c.cfg.Endpoint, secretName, secretVersion, parameters)
}

func (c *AzureClient) DeleteSecret(ctx context.Context, secretName string) (result keyvault.DeletedSecretBundle, err error) {
	return c.client.DeleteSecret(ctx, c.cfg.Endpoint, secretName)
}
