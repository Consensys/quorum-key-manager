package client

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/api/types"
)

const keysPath = "keys"

func (c *HTTPClient) CreateKey(ctx context.Context, storeName string, req *types.CreateKeyRequest) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), keysPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) ImportKey(ctx context.Context, storeName string, req *types.ImportKeyRequest) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s/import", withURLStore(c.config.URL, storeName), keysPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) Sign(ctx context.Context, storeName, id string, req *types.SignBase64PayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign", withURLStore(c.config.URL, storeName), keysPath, id)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), keysPath, id)

	response, err := getRequest(ctx, c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (c *HTTPClient) ListKeys(ctx context.Context, storeName string) ([]string, error) {
	var ids []string
	reqURL := fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), keysPath)
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

func (c *HTTPClient) DestroyKey(ctx context.Context, storeName, id string) error {
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), keysPath, id)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}
