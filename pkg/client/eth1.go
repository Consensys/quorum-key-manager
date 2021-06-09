package client

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
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

func (c *HTTPClient) ImportEth1Account(ctx context.Context, storeName string, req *types.ImportEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	eth1Acc := &types.Eth1AccountResponse{}
	reqURL := fmt.Sprintf("%s/%s/import", withURLStore(c.config.URL, storeName), eth1Path)
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

func (c *HTTPClient) UpdateEth1Account(ctx context.Context, storeName string, req *types.UpdateEth1AccountRequest) (*types.Eth1AccountResponse, error) {
	eth1Acc := &types.Eth1AccountResponse{}
	reqURL := fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), eth1Path)
	response, err := patchRequest(ctx, c.client, reqURL, req)
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

func (c *HTTPClient) SignEth1(ctx context.Context, storeName, account string, req *types.SignHexPayloadRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign", withURLStore(c.config.URL, storeName), eth1Path, account)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) SignTypedData(ctx context.Context, storeName, account string, req *types.SignTypedDataRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-typed-data", withURLStore(c.config.URL, storeName), eth1Path, account)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) SignTransaction(ctx context.Context, storeName, account string, req *types.SignETHTransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-transaction", withURLStore(c.config.URL, storeName), eth1Path, account)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) SignQuorumPrivateTransaction(ctx context.Context, storeName, account string, req *types.SignQuorumPrivateTransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-quorum-private-transaction", withURLStore(c.config.URL, storeName), eth1Path, account)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) SignEEATransaction(ctx context.Context, storeName, account string, req *types.SignEEATransactionRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign-eea-transaction", withURLStore(c.config.URL, storeName), eth1Path, account)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) GetEth1Account(ctx context.Context, storeName, account string) (*types.Eth1AccountResponse, error) {
	acc := &types.Eth1AccountResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), eth1Path, account)

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

func (c *HTTPClient) ListEth1Accounts(ctx context.Context, storeName string) ([]string, error) {
	var accs []string
	reqURL := fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), eth1Path)
	response, err := getRequest(ctx, c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, &accs)
	if err != nil {
		return nil, err
	}

	return accs, nil
}

func (c *HTTPClient) DeleteEth1Account(ctx context.Context, storeName, account string) error {
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), eth1Path, account)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) DestroyEth1Account(ctx context.Context, storeName, account string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/destroy", withURLStore(c.config.URL, storeName), eth1Path, account)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) RestoreEth1Account(ctx context.Context, storeName, account string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/restore", withURLStore(c.config.URL, storeName), eth1Path, account)
	response, err := postRequest(ctx, c.client, reqURL, nil)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) ECRecover(ctx context.Context, storeName string, req *types.ECRecoverRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/ec-recover", withURLStore(c.config.URL, storeName), eth1Path)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) VerifyEth1Signature(ctx context.Context, storeName string, req *types.VerifyEth1SignatureRequest) error {
	reqURL := fmt.Sprintf("%s/%s/verify-signature", withURLStore(c.config.URL, storeName), eth1Path)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) VerifyTypedDataSignature(ctx context.Context, storeName string, req *types.VerifyTypedDataRequest) error {
	reqURL := fmt.Sprintf("%s/%s/verify-typed-data-signature", withURLStore(c.config.URL, storeName), eth1Path)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}
