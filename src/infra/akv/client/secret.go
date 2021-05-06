package client

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
)

func (c *AKVClient) SetSecret(ctx context.Context, secretName, value string, tags map[string]string) (keyvault.SecretBundle, error) {
	result, err := c.client.SetSecret(ctx, c.cfg.Endpoint, secretName, keyvault.SecretSetParameters{
		Value: &value,
		Tags:  common.Tomapstrptr(tags),
	})
	if err != nil {
		return result, ParseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetSecret(ctx context.Context, secretName, secretVersion string) (keyvault.SecretBundle, error) {
	result, err := c.client.GetSecret(ctx, c.cfg.Endpoint, secretName, secretVersion)
	if err != nil {
		return result, ParseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetSecrets(ctx context.Context, maxResults int32) ([]keyvault.SecretItem, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}
	res, err := c.client.GetSecrets(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, ParseErrorResponse(err)
	}

	if len(res.Values()) == 0 {
		return []keyvault.SecretItem{}, nil
	}

	return res.Values(), nil
}

func (c *AKVClient) UpdateSecret(ctx context.Context, secretName, secretVersion string, expireAt time.Time) (keyvault.SecretBundle, error) {
	expireAtDate := date.NewUnixTimeFromNanoseconds(expireAt.UnixNano())
	result, err := c.client.UpdateSecret(ctx, c.cfg.Endpoint, secretName, secretVersion, keyvault.SecretUpdateParameters{
		SecretAttributes: &keyvault.SecretAttributes{
			Expires: &expireAtDate,
		},
	})
	if err != nil {
		return result, ParseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) DeleteSecret(ctx context.Context, secretName string) (keyvault.DeletedSecretBundle, error) {
	result, err := c.client.DeleteSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return result, ParseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetDeletedSecret(ctx context.Context, secretName string) (keyvault.DeletedSecretBundle, error) {
	result, err := c.client.GetDeletedSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return result, ParseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) PurgeDeletedSecret(ctx context.Context, secretName string) (bool, error) {
	res, err := c.client.PurgeDeletedSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return false, ParseErrorResponse(err)
	}

	return res.StatusCode == http.StatusNoContent, nil
}
