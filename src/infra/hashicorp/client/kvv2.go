package client

import (
	"fmt"
	"path"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/hashicorp/vault/api"
)

func (c *HashicorpVaultClient) ReadData(id string, data map[string][]string) (*api.Secret, error) {
	secret, err := c.client.Logical().ReadWithData(c.pathData(id), data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) ReadMetadata(id string) (*api.Secret, error) {
	secret, err := c.client.Logical().Read(c.pathMetadata(id))
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) SetSecret(id string, data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(c.pathData(id), map[string]interface{}{
		dataLabel: data,
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) DeleteSecret(id string, data map[string][]string) error {
	_, err := c.client.Logical().DeleteWithData(c.pathData(id), data)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *HashicorpVaultClient) RestoreSecret(id string, data map[string][]string) error {
	err := c.writePost(path.Join(c.mountPoint, "undelete", id), data)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *HashicorpVaultClient) DestroySecret(id string, data map[string][]string) error {
	err := c.writePost(path.Join(c.mountPoint, "destroy", id), data)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *HashicorpVaultClient) ListSecrets() (*api.Secret, error) {
	secret, err := c.client.Logical().List(c.pathMetadata(""))
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) writePost(endpoint string, data map[string][]string) error {
	req := c.client.NewRequest("POST", fmt.Sprintf("/v1/%s", endpoint))
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

func (c *HashicorpVaultClient) pathData(id string) string {
	return path.Join(c.mountPoint, dataLabel, id)
}

func (c *HashicorpVaultClient) pathMetadata(id string) string {
	return path.Join(c.mountPoint, "metadata", id)
}
