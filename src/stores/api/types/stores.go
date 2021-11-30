package types

type CreateSecretStoreRequest struct {
	AllowedTenants []string `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
	Vault          string   `json:"vault" validate:"required" yaml:"vault" example:"hashicorp-kv-v2"`
}

type CreateKeyStoreRequest struct {
	AllowedTenants []string `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
	SecretStore    string   `json:"secretStore,omitempty" yaml:"secret_store,omitempty" example:"my-secret-store"`
	Vault          string   `json:"vault,omitempty" yaml:"vault,omitempty" example:"hashicorp-quorum"`
}

type CreateEthereumStoreRequest struct {
	AllowedTenants []string `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
	KeyStore       string   `json:"keyStore" yaml:"key_store" validate:"required" example:"my-key-store"`
}
