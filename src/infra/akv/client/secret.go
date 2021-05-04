package client

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
)

func (c *AKVClient) SetSecret(ctx context.Context, secretName string, value string, tags map[string]string) (keyvault.SecretBundle, error) {
	return c.client.SetSecret(ctx, c.cfg.Endpoint, secretName, keyvault.SecretSetParameters{
		Value: &value,
		Tags:  common.Tomapstrptr(tags),
	})
}

func (c *AKVClient) GetSecret(ctx context.Context, secretName, secretVersion string) (result keyvault.SecretBundle, err error) {
	return c.client.GetSecret(ctx, c.cfg.Endpoint, secretName, secretVersion)
}

func (c *AKVClient) GetSecrets(ctx context.Context, maxResults int32) ([]keyvault.SecretItem, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}
	res, err := c.client.GetSecrets(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, err
	}

	if len(res.Values()) == 0 {
		return []keyvault.SecretItem{}, nil
	}

	return res.Values(), nil
}

func (c *AKVClient) UpdateSecret(ctx context.Context, secretName, secretVersion string, expireAt time.Time) (result keyvault.SecretBundle, err error) {
	expireAtDate := date.NewUnixTimeFromNanoseconds(expireAt.UnixNano())
	return c.client.UpdateSecret(ctx, c.cfg.Endpoint, secretName, secretVersion, keyvault.SecretUpdateParameters{
		SecretAttributes: &keyvault.SecretAttributes{
			Expires: &expireAtDate,
		},
	})
}

func (c *AKVClient) DeleteSecret(ctx context.Context, secretName string) (result keyvault.DeletedSecretBundle, err error) {
	return c.client.DeleteSecret(ctx, c.cfg.Endpoint, secretName)
}
