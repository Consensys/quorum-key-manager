package client

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/hashicorp/vault/api"
)

type HashicorpVaultClient struct {
	client *api.Client
	cfg    *Config
}

var _ hashicorp.VaultClient = &HashicorpVaultClient{}

func NewClient(cfg *Config) (*HashicorpVaultClient, error) {

	clientConfig := cfg.ToHashicorpConfig()
	if clientConfig.Error != nil {
		return nil, clientConfig.Error
	}
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	client.SetNamespace(cfg.Namespace)

	return &HashicorpVaultClient{client, cfg}, nil
}

func (c *HashicorpVaultClient) Read(path string, data map[string][]string) (*api.Secret, error) {
	if data == nil {
		secret, err := c.client.Logical().Read(path)
		if err != nil {
			return nil, parseErrorResponse(err)
		}

		return secret, nil
	}

	secret, err := c.client.Logical().ReadWithData(path, data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) Delete(path string, data map[string][]string) error {
	_, err := c.client.Logical().DeleteWithData(path, data)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *HashicorpVaultClient) WritePost(path string, data map[string][]string) error {
	req := c.client.NewRequest("POST", fmt.Sprintf("/v1/%s", path))
	if data != nil {
		if err := req.SetJSONBody(data); err != nil {
			return errors.EncodingError(err.Error())
		}
	}

	resp, err := c.client.RawRequest(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *HashicorpVaultClient) List(path string) (*api.Secret, error) {
	secret, err := c.client.Logical().List(path)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) SetToken(token string) {
	c.client.SetToken(token)
}

func (c *HashicorpVaultClient) UnwrapToken(token string) (*api.Secret, error) {
	secret, err := c.client.Logical().Unwrap(token)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) HealthCheck() error {
	resp, err := c.client.Sys().Health()
	if err != nil {
		return parseErrorResponse(err)
	}

	if !resp.Initialized {
		errMessage := "client is not initialized"
		return errors.HashicorpVaultError(errMessage)
	}

	return nil
}

func (c *HashicorpVaultClient) Mount(path string, mountInfo *api.MountInput) error {
	err := c.client.Sys().Mount(path, mountInfo)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}
