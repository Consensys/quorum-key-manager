package types

type CreateSecretStoreRequest struct {
	Vault string `json:"vault" validate:"required" yaml:"vault" example:"hashicorp-kv-v2"`
}

type CreateKeyStoreRequest struct {
	SecretStore string                 `json:"secretStore,omitempty" yaml:"secret_store,omitempty" example:"my-secret-store"`
	Vault       string                 `json:"vault,omitempty" yaml:"vault,omitempty" example:"hashicorp-quorum"`
	Properties  map[string]interface{} `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type CreateEthereumStoreRequest struct {
	KeyStore string `json:"keyStore" yaml:"key_store" validate:"required" example:"my-key-store"`
}
