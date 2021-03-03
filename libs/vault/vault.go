package vault

import (
	hashicorp "github.com/hashicorp/vault/api"
)

type HashicorpVaultClient interface {
	Read(path string) (*hashicorp.Secret, error)
	Write(path string, data map[string]interface{}) (*hashicorp.Secret, error)
	List(path string) (*hashicorp.Secret, error)
	Update(path string, data map[string]interface{}) (*hashicorp.Secret, error)
}
