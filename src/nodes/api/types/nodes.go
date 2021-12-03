package types

import proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"

type CreateNodeRequest struct {
	AllowedTenants []string         `json:"allowedTenants,omitempty" yaml:"allowed_tenants,omitempty" example:"tenant1,tenant2"`
	Config         proxynode.Config `json:"config" yaml:"config" validate:"required"`
}
