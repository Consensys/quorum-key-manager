package client

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
)

const registryPathf = "%s/registries/%s"

// CreateRegistry creates an alias in the registry.
func (c *HTTPClient) CreateRegistry(ctx context.Context, registryName string, req *types.CreateRegistryRequest) (*types.RegistryResponse, error) {
	requestURL := fmt.Sprintf(registryPathf, c.config.URL, registryName)
	resp, err := postRequest(ctx, c.client, requestURL, req)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	var registry types.RegistryResponse
	err = parseResponse(resp, &registry)
	if err != nil {
		return nil, err
	}

	return &registry, nil
}

// GetRegistry lists all aliases from a registry.
func (c *HTTPClient) GetRegistry(ctx context.Context, registryName string) (*types.RegistryResponse, error) {
	requestURL := fmt.Sprintf(registryPathf, c.config.URL, registryName)
	resp, err := getRequest(ctx, c.client, requestURL)
	if err != nil {
		return nil, err
	}
	defer closeResponse(resp)

	var registry types.RegistryResponse
	err = parseResponse(resp, &registry)
	if err != nil {
		return nil, err
	}

	return &registry, nil
}

// DeleteRegistry deletes a registry, with all the aliases it contained.
func (c *HTTPClient) DeleteRegistry(ctx context.Context, registryName string) error {
	requestURL := fmt.Sprintf(registryPathf, c.config.URL, registryName)
	resp, err := deleteRequest(ctx, c.client, requestURL)
	if err != nil {
		return err
	}
	defer closeResponse(resp)

	return parseEmptyBodyResponse(resp)
}
