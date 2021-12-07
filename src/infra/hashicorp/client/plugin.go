package client

import (
	"encoding/base64"
	"path"

	"github.com/hashicorp/vault/api"
)

func (c *HashicorpVaultClient) GetKey(id string) (*api.Secret, error) {
	secret, err := c.client.Logical().Read(c.pathKeys(id))
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) CreateKey(data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(c.pathKeys(""), data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) ImportKey(data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(c.pathKeys("import"), data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) UpdateKey(id string, data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(c.pathKeys(id), data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) DestroyKey(id string) error {
	_, err := c.client.Logical().Delete(path.Join(c.pathKeys(id), "destroy"))
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *HashicorpVaultClient) ListKeys() (*api.Secret, error) {
	secret, err := c.client.Logical().List(c.pathKeys(""))
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) Sign(id string, data []byte) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(path.Join(c.pathKeys(id), "sign"), map[string]interface{}{
		dataLabel: base64.URLEncoding.EncodeToString(data),
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) pathKeys(suffix string) string {
	return path.Join(c.mountPoint, "keys", suffix)
}
