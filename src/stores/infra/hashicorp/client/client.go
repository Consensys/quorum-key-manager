package client

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	hashicorp "github.com/hashicorp/vault/api"
)

type HashicorpVaultClient struct {
	client *hashicorp.Client
}

func NewClient(cfg *Config) (*HashicorpVaultClient, error) {
	client, err := hashicorp.NewClient(cfg.ToHashicorpConfig())
	if err != nil {
		return nil, err
	}

	client.SetNamespace(cfg.Namespace)

	return &HashicorpVaultClient{client}, nil
}

func (c *HashicorpVaultClient) Read(path string, data map[string][]string) (*hashicorp.Secret, error) {
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

func (c *HashicorpVaultClient) Write(path string, data map[string]interface{}) (*hashicorp.Secret, error) {
	secret, err := c.client.Logical().Write(path, data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) Delete(path string) error {
	_, err := c.client.Logical().Delete(path)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *HashicorpVaultClient) List(path string) (*hashicorp.Secret, error) {
	secret, err := c.client.Logical().List(path)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) SetToken(token string) {
	c.client.SetToken(token)
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

func (c *HashicorpVaultClient) Client() *hashicorp.Client {
	return c.client
}
