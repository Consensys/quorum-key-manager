package client

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
)

const keysPath = "keys"

func (c *HTTPClient) CreateKey(ctx context.Context, storeName string, req *types.CreateKeyRequest) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s", c.config.URL, keysPath)
	response, err := postRequest(withStore(ctx, storeName), c.client, reqURL, req)
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
	reqURL := fmt.Sprintf("%s/%s/import", c.config.URL, keysPath)
	response, err := postRequest(withStore(ctx, storeName), c.client, reqURL, req)
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

func (c *HTTPClient) Sign(ctx context.Context, storeName, id string, req *types.SignPayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign", c.config.URL, keysPath, id)
	response, err := postRequest(withStore(ctx, storeName), c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) GetKey(ctx context.Context, storeName, id, version string) (*types.KeyResponse, error) {
	key := &types.KeyResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.URL, keysPath, id)

	if version != "" {
		reqURL = fmt.Sprintf("%s?version=%s", reqURL, version)
	}

	response, err := getRequest(withStore(ctx, storeName), c.client, reqURL)
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
	reqURL := fmt.Sprintf("%s/%s", c.config.URL, keysPath)
	response, err := getRequest(withStore(ctx, storeName), c.client, reqURL)
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
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.URL, keysPath, id)
	response, err := deleteRequest(withStore(ctx, storeName), c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}
