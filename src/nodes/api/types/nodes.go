package types

import proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"

type CreateNodeRequest struct {
	Name           string           `json:"name" yaml:"name" validate:"required" example:"geth-node"`
	AllowedTenants []string         `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
	Config         proxynode.Config `json:"config" yaml:"config" validate:"required"`
}
