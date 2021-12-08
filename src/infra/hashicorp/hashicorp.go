package hashicorp

import (
	hashicorp "github.com/hashicorp/vault/api"
)

//go:generate mockgen -source=hashicorp.go -destination=mocks/hashicorp.go -package=mocks

type Client interface {
	Kvv2Client
	PluginClient
	SetToken(token string)
	UnwrapToken(token string) (*hashicorp.Secret, error)
	Mount(path string, mountInfo *hashicorp.MountInput) error
	HealthCheck() error
}

type Kvv2Client interface {
	ReadData(id string, data map[string][]string) (*hashicorp.Secret, error)
	ReadMetadata(id string) (*hashicorp.Secret, error)
	SetSecret(id string, data map[string]interface{}) (*hashicorp.Secret, error)
	ListSecrets() (*hashicorp.Secret, error)
	DeleteSecret(id string, data map[string][]string) error
	RestoreSecret(id string, data map[string][]string) error
	DestroySecret(id string, data map[string][]string) error
}

type PluginClient interface {
	GetKey(id string) (*hashicorp.Secret, error)
	CreateKey(data map[string]interface{}) (*hashicorp.Secret, error)
	ImportKey(data map[string]interface{}) (*hashicorp.Secret, error)
	ListKeys() (*hashicorp.Secret, error)
	UpdateKey(id string, data map[string]interface{}) (*hashicorp.Secret, error)
	DestroyKey(id string) error
	Sign(id string, data []byte) (*hashicorp.Secret, error)
}
