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

func (c *hashicorpVaultClient) Update(path string, data map[string]interface{}) (*hashicorp.Secret, error) {
	// Update is the same call as write and it automatically create a new version
	return c.client.Write(path, data)
}
