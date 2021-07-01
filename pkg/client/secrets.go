package client

import (
	"context"
	"fmt"

	"github.com/consensysquorum/quorum-key-manager/src/stores/api/types"
)

const secretsPath = "secrets"

func (c *HTTPClient) SetSecret(ctx context.Context, storeName, id string, req *types.SetSecretRequest) (*types.SecretResponse, error) {
	secret := &types.SecretResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), secretsPath, id)
	response, err := postRequest(ctx, c.client, reqURL, req)
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
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), secretsPath, id)

	if version != "" {
		reqURL = fmt.Sprintf("%s?version=%s", reqURL, version)
	}

	response, err := getRequest(ctx, c.client, reqURL)
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
	reqURL := fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), secretsPath)
	response, err := getRequest(ctx, c.client, reqURL)
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
