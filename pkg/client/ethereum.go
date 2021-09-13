package client

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/src/stores/api/types"
)

const ethPath = "ethereum"

func (c *HTTPClient) CreateEthAccount(ctx context.Context, storeName string, req *types.CreateEthAccountRequest) (*types.EthAccountResponse, error) {
	ethAcc := &types.EthAccountResponse{}
	reqURL := fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), ethPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, ethAcc)
	if err != nil {
		return nil, err
	}

	return ethAcc, nil
}

func (c *HTTPClient) ImportEthAccount(ctx context.Context, storeName string, req *types.ImportEthAccountRequest) (*types.EthAccountResponse, error) {
	ethAcc := &types.EthAccountResponse{}
	reqURL := fmt.Sprintf("%s/%s/import", withURLStore(c.config.URL, storeName), ethPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, ethAcc)
	if err != nil {
		return nil, err
	}

	return ethAcc, nil
}

func (c *HTTPClient) UpdateEthAccount(ctx context.Context, storeName, address string, req *types.UpdateEthAccountRequest) (*types.EthAccountResponse, error) {
	ethAcc := &types.EthAccountResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := patchRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, ethAcc)
	if err != nil {
		return nil, err
	}

	return ethAcc, nil
}

func (c *HTTPClient) SignMessage(ctx context.Context, storeName, address string, req *types.SignMessageRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-message", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) SignTypedData(ctx context.Context, storeName, address string, req *types.SignTypedDataRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-typed-data", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) SignTransaction(ctx context.Context, storeName, address string, req *types.SignETHTransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-transaction", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) SignQuorumPrivateTransaction(ctx context.Context, storeName, address string, req *types.SignQuorumPrivateTransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-quorum-private-transaction", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) SignEEATransaction(ctx context.Context, storeName, address string, req *types.SignEEATransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-eea-transaction", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) GetEthAccount(ctx context.Context, storeName, address string) (*types.EthAccountResponse, error) {
	acc := &types.EthAccountResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), ethPath, address)

	response, err := getRequest(ctx, c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, acc)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (c *HTTPClient) ListEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return listRequest(ctx, c.client, fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), ethPath), false, limit, page)
}

func (c *HTTPClient) ListDeletedEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return listRequest(ctx, c.client, fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), ethPath), true, limit, page)
}

func (c *HTTPClient) DeleteEthAccount(ctx context.Context, storeName, address string) error {
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) DestroyEthAccount(ctx context.Context, storeName, address string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/destroy", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) RestoreEthAccount(ctx context.Context, storeName, address string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/restore", withURLStore(c.config.URL, storeName), ethPath, address)
	response, err := putRequest(ctx, c.client, reqURL, nil)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) ECRecover(ctx context.Context, storeName string, req *types.ECRecoverRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/ec-recover", withURLStore(c.config.URL, storeName), ethPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) VerifyMessage(ctx context.Context, storeName string, req *types.VerifyRequest) error {
	reqURL := fmt.Sprintf("%s/%s/verify-message", withURLStore(c.config.URL, storeName), ethPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) VerifyTypedData(ctx context.Context, storeName string, req *types.VerifyTypedDataRequest) error {
	reqURL := fmt.Sprintf("%s/%s/verify-typed-data", withURLStore(c.config.URL, storeName), ethPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}
