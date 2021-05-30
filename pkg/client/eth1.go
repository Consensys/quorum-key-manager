package client

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/api/types"
)

const eth1Path = "eth1"

func (c *HTTPClient) CreateEth1Account(ctx context.Context, storeName string, req *types.CreateEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	eth1Acc := &types.Eth1AccountResponse{}
	reqURL := fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), eth1Path)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, eth1Acc)
	if err != nil {
		return nil, err
	}

	return eth1Acc, nil
}
