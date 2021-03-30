package client

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault/keyvaultapi"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/akv"
)

type azureClient struct {
	client       keyvaultapi.BaseClientAPI
	vaultBaseURL string
}

func NewClient(cfg *Config) (akv.Client, error) {
	client := keyvault.New()

	authorizer, err := cfg.ToAzureAuthConfig()
	if err != nil {
		return nil, err
	}
	client.Authorizer = authorizer

	return &azureClient{client: client, vaultBaseURL: cfg.Endpoint}, nil
}

func (c *azureClient) SetSecret(ctx context.Context, secretName string, parameters keyvault.SecretSetParameters) (result keyvault.SecretBundle, err error) {
	return c.client.SetSecret(ctx, c.vaultBaseURL, secretName, parameters)
}

func (c *azureClient) GetSecret(ctx context.Context, secretName, secretVersion string) (result keyvault.SecretBundle, err error) {
	return c.client.GetSecret(ctx, c.vaultBaseURL, secretName, secretVersion)
}

func (c *azureClient) GetSecrets(ctx context.Context, maxResults *int32) (result keyvault.SecretListResultPage, err error) {
	return c.client.GetSecrets(ctx, c.vaultBaseURL, maxResults)
}

func (c *azureClient) UpdateSecret(ctx context.Context, secretName, secretVersion string, parameters keyvault.SecretUpdateParameters) (result keyvault.SecretBundle, err error) {
	return c.client.UpdateSecret(ctx, c.vaultBaseURL, secretName, secretVersion, parameters)
}

func (c *azureClient) DeleteSecret(ctx context.Context, secretName string) (result keyvault.DeletedSecretBundle, err error) {
	return c.client.DeleteSecret(ctx, c.vaultBaseURL, secretName)
}
