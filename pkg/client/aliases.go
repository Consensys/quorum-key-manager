package client

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
)

const aliasPathf = "%s/registries/%s/aliases/%s"

// CreateAlias creates an alias in the registry.
func (c *HTTPClient) CreateAlias(ctx context.Context, registry, aliasKey string, req *types.AliasRequest) (*types.AliasResponse, error) {
	requestURL := fmt.Sprintf(aliasPathf, c.config.URL, registry, aliasKey)
	resp, err := postRequest(ctx, c.client, requestURL, req)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	var a types.AliasResponse
	err = parseResponse(resp, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAlias gets an alias from the registry.
func (c *HTTPClient) GetAlias(ctx context.Context, registry, aliasKey string) (*types.AliasResponse, error) {
	requestURL := fmt.Sprintf(aliasPathf, c.config.URL, registry, aliasKey)
	resp, err := getRequest(ctx, c.client, requestURL)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	var a types.AliasResponse
	err = parseResponse(resp, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateAlias updates an alias in the registry.
func (c *HTTPClient) UpdateAlias(ctx context.Context, registry, aliasKey string, req *types.AliasRequest) (*types.AliasResponse, error) {
	requestURL := fmt.Sprintf(aliasPathf, c.config.URL, registry, aliasKey)
	resp, err := putRequest(ctx, c.client, requestURL, req)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	var a types.AliasResponse
	err = parseResponse(resp, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// DeleteAlias deletes an alias from the registry.
func (c *HTTPClient) DeleteAlias(ctx context.Context, registry, aliasKey string) error {
	requestURL := fmt.Sprintf(aliasPathf, c.config.URL, registry, aliasKey)
	resp, err := deleteRequest(ctx, c.client, requestURL)
	if err != nil {
		return err
	}
	defer closeResponse(resp)

	return parseEmptyBodyResponse(resp)
}
