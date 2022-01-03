package client

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/src/utils/api/types"
)

const utilsPath = "utilities"

func (c *HTTPClient) VerifyKeySignature(ctx context.Context, req *types.VerifyKeySignatureRequest) error {
	reqURL := fmt.Sprintf("%s/%s/keys/verify-signature", c.config.URL, utilsPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) ECRecover(ctx context.Context, req *types.ECRecoverRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/ethereum/ec-recover", c.config.URL, utilsPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) VerifyMessage(ctx context.Context, req *types.VerifyRequest) error {
	reqURL := fmt.Sprintf("%s/%s/ethereum/verify-message", c.config.URL, utilsPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) VerifyTypedData(ctx context.Context, req *types.VerifyTypedDataRequest) error {
	reqURL := fmt.Sprintf("%s/%s/ethereum/verify-typed-data", c.config.URL, utilsPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}
