package types

import "github.com/consensys/quorum-key-manager/src/stores/entities"

type CreateVaultRequest struct {
	VaultType string      `json:"type" yaml:"type" validate:"required,isVaultType" example:"hashicorp"`
	Params    interface{} `json:"params" yaml:"params" validate:"required"`
}

type CreateHashicorpVaultRequest struct {
	Name   string                   `json:"name" yaml:"name" validate:"required" example:"hashicorp-kv-v2"`
	Config entities.HashicorpConfig `json:"config" yaml:"config" validate:"required"`
}

type CreateAzureVaultRequest struct {
	Name   string               `json:"name" yaml:"name" validate:"required" example:"akv-europe"`
	Config entities.AzureConfig `json:"config" yaml:"config" validate:"required"`
}

type CreateAWSVaultRequest struct {
	Name   string             `json:"name" validate:"required" example:"aws-europe"`
	Config entities.AWSConfig `json:"config" yaml:"config" validate:"required"`
}
