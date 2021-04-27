package client

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
)

const (
	secretsPath          = "secrets"
	secretStoreHeaderKey = "X-Secret-Store"
)

func (c *HTTPClient) SetSecret(ctx context.Context, storeName string, req *types.SetSecretRequest) (*types.SecretResponse, error) {
	secret := &types.SecretResponse{}
	reqURL := fmt.Sprintf("%s/%s", c.config.URL, secretsPath)
	response, err := postRequest(withSecretStore(ctx, storeName), c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (c *HTTPClient) GetSecret(ctx context.Context, storeName, id, version string) (*types.SecretResponse, error) {
	secret := &types.SecretResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s/%s", c.config.URL, secretsPath, id, version)
	response, err := getRequest(withSecretStore(ctx, storeName), c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (c *HTTPClient) ListSecrets(ctx context.Context, storeName string) ([]string, error) {
	var ids []string
	reqURL := fmt.Sprintf("%s/%s", c.config.URL, secretsPath)
	response, err := getRequest(withSecretStore(ctx, storeName), c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func withSecretStore(ctx context.Context, storeName string) context.Context {
	return context.WithValue(ctx, RequestHeaderKey, map[string]string{
		secretStoreHeaderKey: storeName,
	})
}
