package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/nodes"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
)

type NodesHandler struct {
	nodes    nodes.Nodes
	userInfo *auth.UserInfo
}

func NewNodesHandler(nodesService nodes.Nodes) *NodesHandler {
	return &NodesHandler{
		nodes:    nodesService,
		userInfo: auth.NewWildcardUser(), // This handler always use the wildcard user because it's a manifest handler
	}
}

func (h *NodesHandler) Register(ctx context.Context, mnfs []entities2.Manifest) error {
	for _, mnf := range mnfs {
		err := h.Create(ctx, mnf.Name, mnf.AllowedTenants, mnf.Specs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *NodesHandler) Create(ctx context.Context, name string, allowedTenants []string, specs interface{}) error {
	config := &proxynode.Config{}
	err := json.UnmarshalYAML(specs, config)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.nodes.Create(ctx, name, config.SetDefault(), allowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
