package client

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/sdk/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/sdk/utils"
)

const (
	secretsPath = "secrets"
)

func (c *HTTPClient) CreateSecret(ctx context.Context, req *types.CreateSecretRequest) (*types.Secret, error) {
	reqURL := fmt.Sprintf("%v/%s", c.config.URL, secretsPath)
	resp := &types.Secret{}

	response, err := utils.PostRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer utils.CloseResponse(response)
	if err := utils.ParseResponse(response, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
