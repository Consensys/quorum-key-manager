package types

import (
	"github.com/consensys/quorum-key-manager/src/entities"
)

type CreateHashicorpVaultRequest struct {
	Config         entities.HashicorpConfig `json:"config" yaml:"config" validate:"required"`
	AllowedTenants []string                 `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
}

type CreateAzureVaultRequest struct {
	Config         entities.AzureConfig `json:"config" yaml:"config" validate:"required"`
	AllowedTenants []string             `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
}

type CreateAWSVaultRequest struct {
	Config         entities.AWSConfig `json:"config" yaml:"config" validate:"required"`
	AllowedTenants []string           `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
}
