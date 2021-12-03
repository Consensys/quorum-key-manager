package client

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/hashicorp/vault/api"
)

const dataLabel = "data"

type HashicorpVaultClient struct {
	client     *api.Client
	mountPoint string
}

var _ hashicorp.Client = &HashicorpVaultClient{}

func NewClient(cfg *Config) (*HashicorpVaultClient, error) {
	clientConfig, err := cfg.ToHashicorpConfig()
	if err != nil {
		return nil, err
	}
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	client.SetNamespace(cfg.Namespace)

	return &HashicorpVaultClient{client: client, mountPoint: cfg.MountPoint}, nil
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
