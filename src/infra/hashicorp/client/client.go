package client

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	hashicorp2 "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp"
	hashicorp "github.com/hashicorp/vault/api"
)

type hashicorpVaultClient struct {
	client *hashicorp.Client
}

func NewClient(cfg *Config) (hashicorp2.VaultClient, error) {
	client, err := hashicorp.NewClient(cfg.ToHashicorpConfig())
	if err != nil {
		return nil, err
	}

	client.SetToken(cfg.Token)
	return &hashicorpVaultClient{client}, nil
}

func (c *hashicorpVaultClient) Read(path string) (*hashicorp.Secret, error) {
	return c.client.Logical().Read(path)
}

func (c *hashicorpVaultClient) Write(path string, data map[string]interface{}) (*hashicorp.Secret, error) {
	return c.client.Logical().Write(path, data)
}

func (c *hashicorpVaultClient) List(path string) (*hashicorp.Secret, error) {
	return c.client.Logical().List(path)
}

func (c *hashicorpVaultClient) HealthCheck() error {
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

func (c *hashicorpVaultClient) Client() *hashicorp.Client {
	return c.client
}
