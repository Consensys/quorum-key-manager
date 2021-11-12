package types

import "github.com/consensys/quorum-key-manager/src/auth/entities"

type CreateRoleRequest struct {
	Name        string                `json:"name" yaml:"name" validate:"required" example:"admin"`
	Permissions []entities.Permission `json:"permissions" yaml:"permissions" validate:"required" example:"*:*"`
}
