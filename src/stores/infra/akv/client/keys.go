package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
)

func (c *AKVClient) CreateKey(ctx context.Context, keyName string, kty keyvault.JSONWebKeyType,
	crv keyvault.JSONWebKeyCurveName, attr *keyvault.KeyAttributes, ops []keyvault.JSONWebKeyOperation,
	tags map[string]string) (keyvault.KeyBundle, error) {

	result, err := c.client.CreateKey(ctx, c.cfg.Endpoint, keyName, keyvault.KeyCreateParameters{
		Kty:           kty,
		Curve:         crv,
		Tags:          common.Tomapstrptr(tags),
		KeyAttributes: attr,
	})
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) ImportKey(ctx context.Context, keyName string, k *keyvault.JSONWebKey, attr *keyvault.KeyAttributes, tags map[string]string) (keyvault.KeyBundle, error) {
	result, err := c.client.ImportKey(ctx, c.cfg.Endpoint, keyName, keyvault.KeyImportParameters{
		Key:           k,
		Tags:          common.Tomapstrptr(tags),
		KeyAttributes: attr,
	})
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetKey(ctx context.Context, keyName, version string) (keyvault.KeyBundle, error) {
	result, err := c.client.GetKey(ctx, c.cfg.Endpoint, keyName, version)
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetKeys(ctx context.Context, maxResults int32) ([]keyvault.KeyItem, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}
	res, err := c.client.GetKeys(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	items := []keyvault.KeyItem{}
	for {
		items = append(items, res.Values()...)
		if !res.NotDone() {
			break
		}

		err := res.NextWithContext(ctx)
		if err != nil {
			return items, err
		}

		if maxResults != 0 && len(items) >= int(maxResults) {
			break
		}
	}

	if maxResults != 0 && len(items) > int(maxResults) {
		return items[0:maxResults], nil
	}

	return items, nil
}

func (c *AKVClient) UpdateKey(ctx context.Context, keyName, version string, attr *keyvault.KeyAttributes,
	ops []keyvault.JSONWebKeyOperation, tags map[string]string) (keyvault.KeyBundle, error) {
	result, err := c.client.UpdateKey(ctx, c.cfg.Endpoint, keyName, version, keyvault.KeyUpdateParameters{
		KeyAttributes: attr,
		Tags:          common.Tomapstrptr(tags),
		KeyOps:        &ops,
	})
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) DeleteKey(ctx context.Context, keyName string) (keyvault.DeletedKeyBundle, error) {
	result, err := c.client.DeleteKey(ctx, c.cfg.Endpoint, keyName)
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetDeletedKey(ctx context.Context, keyName string) (keyvault.DeletedKeyBundle, error) {
	result, err := c.client.GetDeletedKey(ctx, c.cfg.Endpoint, keyName)
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetDeletedKeys(ctx context.Context, maxResults int32) ([]keyvault.DeletedKeyItem, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}
	res, err := c.client.GetDeletedKeys(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, parseErrorResponse(err)
	}
	if len(res.Values()) == 0 {
		return []keyvault.DeletedKeyItem{}, nil
	}

	return res.Values(), nil
}

func (c *AKVClient) PurgeDeletedKey(ctx context.Context, keyName string) (bool, error) {
	res, err := c.client.PurgeDeletedKey(ctx, c.cfg.Endpoint, keyName)
	if err != nil {
		return false, parseErrorResponse(err)
	}
	return res.StatusCode == http.StatusNoContent, nil
}

func (c *AKVClient) RecoverDeletedKey(ctx context.Context, keyName string) (keyvault.KeyBundle, error) {
	result, err := c.client.RecoverDeletedKey(ctx, c.cfg.Endpoint, keyName)
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) Sign(ctx context.Context, keyName, version string, alg keyvault.JSONWebKeySignatureAlgorithm, payload string) (string, error) {
	res, err := c.client.Sign(ctx, c.cfg.Endpoint, keyName, version, keyvault.KeySignParameters{
		Value:     &payload,
		Algorithm: alg,
	})
	if err != nil {
		return "", parseErrorResponse(err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to sign")
	}
	return *res.Result, nil
}

func (c *AKVClient) Encrypt(ctx context.Context, keyName, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, payload string) (string, error) {
	res, err := c.client.Encrypt(ctx, c.cfg.Endpoint, keyName, version, keyvault.KeyOperationsParameters{
		Value:     &payload,
		Algorithm: alg,
	})
	if err != nil {
		return "", parseErrorResponse(err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to encrypt")
	}
	return *res.Result, nil
}

func (c *AKVClient) Decrypt(ctx context.Context, keyName, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, value string) (string, error) {
	res, err := c.client.Decrypt(ctx, c.cfg.Endpoint, keyName, version, keyvault.KeyOperationsParameters{
		Value:     &value,
		Algorithm: alg,
	})
	if err != nil {
		return "", parseErrorResponse(err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to decrypt")
	}

	return *res.Result, nil
}
