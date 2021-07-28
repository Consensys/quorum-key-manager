package client

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/pkg/common"
)

func (c *AKVClient) SetSecret(ctx context.Context, secretName, value string, tags map[string]string) (keyvault.SecretBundle, error) {
	result, err := c.client.SetSecret(ctx, c.cfg.Endpoint, secretName, keyvault.SecretSetParameters{
		Value: &value,
		Tags:  common.Tomapstrptr(tags),
	})
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetSecret(ctx context.Context, secretName, secretVersion string) (keyvault.SecretBundle, error) {
	result, err := c.client.GetSecret(ctx, c.cfg.Endpoint, secretName, secretVersion)
	if err != nil {
		return result, parseErrorResponse(err)
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
		return nil, parseErrorResponse(err)
	}

	items := []keyvault.SecretItem{}
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

func (c *AKVClient) UpdateSecret(ctx context.Context, secretName, secretVersion string, expireAt time.Time) (keyvault.SecretBundle, error) {
	expireAtDate := date.NewUnixTimeFromNanoseconds(expireAt.UnixNano())
	result, err := c.client.UpdateSecret(ctx, c.cfg.Endpoint, secretName, secretVersion, keyvault.SecretUpdateParameters{
		SecretAttributes: &keyvault.SecretAttributes{
			Expires: &expireAtDate,
		},
	})
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) DeleteSecret(ctx context.Context, secretName string) (keyvault.DeletedSecretBundle, error) {
	result, err := c.client.DeleteSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) RecoverSecret(ctx context.Context, secretName string) (keyvault.SecretBundle, error) {
	result, err := c.client.RecoverDeletedSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return result, parseErrorResponse(err)
	}

	return result, nil
}

func (c *AKVClient) GetDeletedSecret(ctx context.Context, secretName string) (keyvault.DeletedSecretBundle, error) {
	result, err := c.client.GetDeletedSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return result, parseErrorResponse(err)
	}
	return result, nil
}

func (c *AKVClient) GetDeletedSecrets(ctx context.Context, maxResults int32) ([]keyvault.DeletedSecretItem, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}

	res, err := c.client.GetDeletedSecrets(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	items := []keyvault.DeletedSecretItem{}
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

func (c *AKVClient) PurgeDeletedSecret(ctx context.Context, secretName string) (bool, error) {
	res, err := c.client.PurgeDeletedSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return false, parseErrorResponse(err)
	}

	return res.StatusCode == http.StatusNoContent, nil
}
