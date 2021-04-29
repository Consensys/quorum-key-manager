package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
)

func (c *AKVClient) CreateKey(ctx context.Context, keyName string, kty keyvault.JSONWebKeyType,
	crv keyvault.JSONWebKeyCurveName, attr *keyvault.KeyAttributes, ops []keyvault.JSONWebKeyOperation,
	tags map[string]string) (keyvault.KeyBundle, error) {
	if crv == "" {
		return keyvault.KeyBundle{}, errors.InvalidParameterError("key curve name cannot be empty")
	}
	if kty == "" {
		return keyvault.KeyBundle{}, errors.InvalidParameterError("key type cannot be empty")
	}

	return c.client.CreateKey(ctx, c.cfg.Endpoint, keyName, keyvault.KeyCreateParameters{
		Kty:           kty,
		Curve:         crv,
		Tags:          common.Tomapstrptr(tags),
		// KeyOps:        &ops,
		KeyAttributes: attr,
	})
}

func (c *AKVClient) ImportKey(ctx context.Context, keyName string, k *keyvault.JSONWebKey, tags map[string]string) (keyvault.KeyBundle, error) {
	if k.Crv == "" {
		return keyvault.KeyBundle{}, errors.InvalidParameterError("key curve name cannot be empty")
	}
	if k.Kty == "" {
		return keyvault.KeyBundle{}, errors.InvalidParameterError("key type cannot be empty")
	}

	return c.client.ImportKey(ctx, c.cfg.Endpoint, keyName, keyvault.KeyImportParameters{
		Key:  k,
		Tags: common.Tomapstrptr(tags),
	})
}

func (c *AKVClient) GetKey(ctx context.Context, keyName string, version string) (keyvault.KeyBundle, error) {
	return c.client.GetKey(ctx, c.cfg.Endpoint, keyName, version)
}

func (c *AKVClient) GetKeys(ctx context.Context, maxResults int32) ([]keyvault.KeyItem, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}
	res, err := c.client.GetKeys(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, err
	}

	if len(res.Values()) == 0 {
		return []keyvault.KeyItem{}, nil
	}

	return res.Values(), nil
}

func (c *AKVClient) UpdateKey(ctx context.Context, keyName string, version string, attr *keyvault.KeyAttributes,
	ops []keyvault.JSONWebKeyOperation, tags map[string]string) (keyvault.KeyBundle, error) {
	return c.client.UpdateKey(ctx, c.cfg.Endpoint, keyName, version, keyvault.KeyUpdateParameters{
		KeyAttributes: attr,
		Tags:          common.Tomapstrptr(tags),
		KeyOps:        &ops,
	})
}

func (c *AKVClient) DeleteKey(ctx context.Context, keyName string) (result keyvault.DeletedKeyBundle, err error) {
	return c.client.DeleteKey(ctx, c.cfg.Endpoint, keyName)
}

func (c *AKVClient) GetDeletedKey(ctx context.Context, keyName string) (keyvault.DeletedKeyBundle, error) {
	return c.client.GetDeletedKey(ctx, c.cfg.Endpoint, keyName)
}

func (c *AKVClient) GetDeletedKeys(ctx context.Context, maxResults int32) ([]keyvault.DeletedKeyItem, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}
	res, err := c.client.GetDeletedKeys(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, err
	}
	if len(res.Values()) == 0 {
		return []keyvault.DeletedKeyItem{}, nil
	}

	return res.Values(), nil
}

func (c *AKVClient) PurgeDeletedKey(ctx context.Context, keyName string) (bool, error) {
	res, err := c.client.PurgeDeletedKey(ctx, c.cfg.Endpoint, keyName)
	if err != nil {
		return false, err
	}
	return res.StatusCode == http.StatusNoContent, nil
}

func (c *AKVClient) RecoverDeletedKey(ctx context.Context, keyName string) (keyvault.KeyBundle, error) {
	return c.client.RecoverDeletedKey(ctx, c.cfg.Endpoint, keyName)
}

func (c *AKVClient) Sign(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeySignatureAlgorithm, payload string) (string, error) {
	if alg == "" {
		return "", errors.InvalidParameterError("key signature algorithm cannot be empty")
	}

	res, err := c.client.Sign(ctx, c.cfg.Endpoint, keyName, version, keyvault.KeySignParameters{
		Value: &payload,
		Algorithm: alg,
	})
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to sign")
	}
	return *res.Result, nil
}

func (c *AKVClient) Encrypt(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, payload string) (string, error) {
	if alg == "" {
		return "", errors.InvalidParameterError("key signature algorithm cannot be empty")
	}

	res, err := c.client.Encrypt(ctx, c.cfg.Endpoint, keyName, version, keyvault.KeyOperationsParameters{
		Value: &payload,
		Algorithm: alg,
	})
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to encrypt")
	}
	return *res.Result, nil
}

func (c *AKVClient) Decrypt(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, value string) (string, error) {
	if alg == "" {
		return "", errors.InvalidParameterError("key signature algorithm cannot be empty")
	}

	res, err := c.client.Decrypt(ctx, c.cfg.Endpoint, keyName, version, keyvault.KeyOperationsParameters{
		Value: &value,
		Algorithm: alg,
	})
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to decrypt")
	}

	return *res.Result, nil
}
