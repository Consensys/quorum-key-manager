package client

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/infra/akv/utils"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

func (c *AKVClient) SetSecret(ctx context.Context, secretName, value string, tags map[string]string) (*entities.Secret, error) {
	result, err := c.client.SetSecret(ctx, c.cfg.Endpoint, secretName, keyvault.SecretSetParameters{
		Value: &value,
		Tags:  common.Tomapstrptr(tags),
	})

	if err != nil {
		return nil, utils.ErrorResponse(err)
	}

	return utils.ParseSecretBundle(&result), nil
}

func (c *AKVClient) GetSecret(ctx context.Context, secretName, secretVersion string) (*entities.Secret, error) {
	result, err := c.client.GetSecret(ctx, c.cfg.Endpoint, secretName, secretVersion)
	if err != nil {
		return nil, utils.ErrorResponse(err)
	}
	return utils.ParseSecretBundle(&result), nil
}

func (c *AKVClient) GetSecrets(ctx context.Context, maxResults int32) ([]*entities.Secret, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}

	res, err := c.client.GetSecrets(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, utils.ErrorResponse(err)
	}

	items := []*entities.Secret{}
	for {
		for _, v := range res.Values() {
			if maxResults != 0 && len(items) >= int(maxResults) {
				return items, nil
			}
			items = append(items, utils.ParseSecretItem(&v))
		}

		if !res.NotDone() {
			break
		}

		err := res.NextWithContext(ctx)
		if err != nil {
			return items, err
		}
	}

	return items, nil
}

func (c *AKVClient) UpdateSecret(ctx context.Context, secretName, secretVersion string, expireAt time.Time) (*entities.Secret, error) {
	expireAtDate := date.NewUnixTimeFromNanoseconds(expireAt.UnixNano())
	result, err := c.client.UpdateSecret(ctx, c.cfg.Endpoint, secretName, secretVersion, keyvault.SecretUpdateParameters{
		SecretAttributes: &keyvault.SecretAttributes{
			Expires: &expireAtDate,
		},
	})
	if err != nil {
		return nil, utils.ErrorResponse(err)
	}

	return utils.ParseSecretBundle(&result), nil
}

func (c *AKVClient) DeleteSecret(ctx context.Context, secretName string) (*entities.Secret, error) {
	result, err := c.client.DeleteSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return nil, utils.ErrorResponse(err)
	}
	return utils.ParseDeleteSecretBundle(&result), nil
}

func (c *AKVClient) RecoverSecret(ctx context.Context, secretName string) (*entities.Secret, error) {
	result, err := c.client.RecoverDeletedSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return nil, utils.ErrorResponse(err)
	}

	return utils.ParseSecretBundle(&result), nil
}

func (c *AKVClient) GetDeletedSecret(ctx context.Context, secretName string) (*entities.Secret, error) {
	result, err := c.client.GetDeletedSecret(ctx, c.cfg.Endpoint, secretName)
	if err != nil {
		return nil, utils.ErrorResponse(err)
	}
	return utils.ParseDeleteSecretBundle(&result), nil
}

func (c *AKVClient) GetDeletedSecrets(ctx context.Context, maxResults int32) ([]*entities.Secret, error) {
	maxResultPtr := &maxResults
	if maxResults == 0 {
		maxResultPtr = nil
	}

	res, err := c.client.GetDeletedSecrets(ctx, c.cfg.Endpoint, maxResultPtr)
	if err != nil {
		return nil, utils.ErrorResponse(err)
	}

	items := []*entities.Secret{}
	for {
		for _, v := range res.Values() {
			items = append(items, utils.ParseDeletedSecretItem(&v))
		}

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
		return false, utils.ErrorResponse(err)
	}

	return res.StatusCode == http.StatusNoContent, nil
}
