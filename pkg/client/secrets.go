package client

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/src/stores/api/types"
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

func (c *HTTPClient) DeleteSecret(ctx context.Context, storeName, id string) error {
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), secretsPath, id)

	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	err = parseResponse(response, new(string))
	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPClient) RestoreSecret(ctx context.Context, storeName, id string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/restore", withURLStore(c.config.URL, storeName), secretsPath, id)

	response, err := putRequest(ctx, c.client, reqURL, nil)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	err = parseResponse(response, new(string))
	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPClient) DestroySecret(ctx context.Context, storeName, id string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/destroy", withURLStore(c.config.URL, storeName), secretsPath, id)

	response, err := putRequest(ctx, c.client, reqURL, nil)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	err = parseResponse(response, new(string))
	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPClient) GetDeletedSecret(ctx context.Context, storeName, id string) (*types.SecretResponse, error) {
	secret := &types.SecretResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s/deleted", withURLStore(c.config.URL, storeName), secretsPath, id)

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

func (c *HTTPClient) ListDeletedSecrets(ctx context.Context, storeName string) ([]string, error) {
	var ids []string
	reqURL := fmt.Sprintf("%s/%s/deleted", withURLStore(c.config.URL, storeName), secretsPath)
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
