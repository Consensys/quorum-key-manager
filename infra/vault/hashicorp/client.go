package hashicorp

import (
	hashicorp "github.com/hashicorp/vault/api"
)

type hashicorpVaultClient struct {
	client *hashicorp.Logical
}

func New(client *hashicorp.Client) *hashicorpVaultClient {
	return &hashicorpVaultClient{
		client.Logical(),
	}
}

func (c *hashicorpVaultClient) Read(path string) (*hashicorp.Secret, error) {
	return c.client.Read(path)
}

func (c *hashicorpVaultClient) Write(path string, data map[string]interface{}) (*hashicorp.Secret, error) {
	return c.client.Write(path, data)
}

func (c *hashicorpVaultClient) List(path string) (*hashicorp.Secret, error) {
	return c.client.List(path)
}
