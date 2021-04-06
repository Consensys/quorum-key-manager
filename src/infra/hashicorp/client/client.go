package client

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
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

	client.SetToken(cfg.Token)
	return &HashicorpVaultClient{client}, nil
}

func (c *HashicorpVaultClient) ReadWithData(path string, data map[string][]string) (*hashicorp.Secret, error) {
	return c.client.Logical().ReadWithData(path, data)
}

func (c *HashicorpVaultClient) Read(path string) (*hashicorp.Secret, error) {
	return c.client.Logical().Read(path)
}

func (c *HashicorpVaultClient) Write(path string, data map[string]interface{}) (*hashicorp.Secret, error) {
	return c.client.Logical().Write(path, data)
}

func (c *HashicorpVaultClient) List(path string) (*hashicorp.Secret, error) {
	return c.client.Logical().List(path)
}

func (c *HashicorpVaultClient) HealthCheck() error {
	resp, err := c.client.Sys().Health()
	if err != nil {
		return err
	}

	if !resp.Initialized {
		errMessage := "client is not initialized"
		return errors.HashicorpVaultConnectionError(errMessage)
	}

	return nil
}

func (c *HashicorpVaultClient) Client() *hashicorp.Client {
	return c.client
}
