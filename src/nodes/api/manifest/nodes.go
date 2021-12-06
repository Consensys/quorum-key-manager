package manifest

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/json"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/nodes"
	"github.com/consensys/quorum-key-manager/src/nodes/api/types"
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

func (h *NodesHandler) Create(ctx context.Context, name string, specs interface{}) error {
	createReq := &types.CreateNodeRequest{}
	err := json.UnmarshalYAML(specs, createReq)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.nodes.Create(ctx, name, createReq.Config.SetDefault(), createReq.AllowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
