package client

import (
	"context"
	"fmt"
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
)

const eth1Path = "eth1"

func (c *HTTPClient) CreateEth1Account(ctx context.Context, storeName string, req *types2.CreateEth1AccountRequest) (*types2.Eth1AccountResponse, error) {
	eth1Acc := &types2.Eth1AccountResponse{}
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
