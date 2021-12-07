package src

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth"
	rolesapi "github.com/consensys/quorum-key-manager/src/auth/api/manifest"
	"github.com/consensys/quorum-key-manager/src/entities"
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/yaml"
	"github.com/consensys/quorum-key-manager/src/nodes"
	nodesapi "github.com/consensys/quorum-key-manager/src/nodes/api/manifest"
	"github.com/consensys/quorum-key-manager/src/stores"
	storesapi "github.com/consensys/quorum-key-manager/src/stores/api/manifest"
	"github.com/consensys/quorum-key-manager/src/vaults"
	vaultsapi "github.com/consensys/quorum-key-manager/src/vaults/api/manifest"
)

func initialize(
	ctx context.Context,
	cfg *manifestreader.Config,
	rolesService auth.Roles,
	vaultsService vaults.Vaults,
	storesService stores.Stores,
	nodesService nodes.Nodes,
) error {
	manifestReader, err := manifestreader.New(cfg)
	if err != nil {
		return err
	}

	manifests, err := manifestReader.Load(ctx)
	if err != nil {
		return err
	}

	// Note that order is important here as stores depend on the existing vaults, do not use a switch!

	err = rolesapi.NewRolesHandler(rolesService).Register(ctx, manifests[entities.RoleKind])
	if err != nil {
		return err
	}

	err = vaultsapi.NewVaultsHandler(vaultsService).Register(ctx, manifests[entities.VaultKind])
	if err != nil {
		return err
	}

	err = storesapi.NewStoresHandler(storesService).Register(ctx, manifests[entities.StoreKind])
	if err != nil {
		return err
	}

	err = nodesapi.NewNodesHandler(nodesService).Register(ctx, manifests[entities.NodeKind])
	if err != nil {
		return err
	}

	return nil
}
