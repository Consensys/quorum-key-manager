package client

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
)

const (
	secretsPath = "secrets"
)

func (c *HTTPClient) Set(ctx context.Context, req *types.SetSecretRequest) (*types.SecretResponse, error) {
	secret := &types.SecretResponse{}
	reqURL := fmt.Sprintf("%s/%s", c.config.URL, secretsPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(ctx, response, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
