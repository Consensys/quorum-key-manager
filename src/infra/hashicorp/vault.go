package hashicorp

import (
	hashicorp "github.com/hashicorp/vault/api"
)

//go:generate mockgen -source=vault.go -destination=mocks/vault.go -package=mocks

type VaultClient interface {
	Read(path string) (*hashicorp.Secret, error)
	Write(path string, data map[string]interface{}) (*hashicorp.Secret, error)
	List(path string) (*hashicorp.Secret, error)
	Client() *hashicorp.Client
	HealthCheck() error
}
