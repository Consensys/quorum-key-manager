package src

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth"
	authapi "github.com/consensys/quorum-key-manager/src/auth/api/manifest"
	"github.com/consensys/quorum-key-manager/src/entities"
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/yaml"
	"github.com/consensys/quorum-key-manager/src/nodes"
	nodesapi "github.com/consensys/quorum-key-manager/src/nodes/api/manifest"
	"github.com/consensys/quorum-key-manager/src/stores"
	storesapi "github.com/consensys/quorum-key-manager/src/stores/api/manifest"
	entities2 "github.com/consensys/quorum-key-manager/src/stores/entities"
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

	manifestRolesHandler := authapi.NewRolesHandler(rolesService)
	for _, mnf := range manifests[entities.RoleKind] {
		err = manifestRolesHandler.Create(ctx, mnf.Specs)
		if err != nil {
			return err
		}
	}

	manifestVaultHandler := vaultsapi.NewVaultsHandler(vaultsService)
	for _, mnf := range manifests[entities.VaultKind] {
		switch mnf.ResourceType {
		case entities.HashicorpVaultType:
			return manifestVaultHandler.CreateHashicorp(ctx, mnf.Specs)
		case entities.AzureVaultType:
			return manifestVaultHandler.CreateAzure(ctx, mnf.Specs)
		case entities.AWSVaultType:
			return manifestVaultHandler.CreateAWS(ctx, mnf.Specs)
		default:
			return errors.InvalidFormatError("invalid vault type")
		}
	}

	manifestStoreHandler := storesapi.NewStoresHandler(storesService)
	for _, mnf := range manifests[entities.StoreKind] {
		switch mnf.ResourceType {
		case entities2.SecretStoreType:
			err = manifestStoreHandler.CreateSecret(ctx, mnf.Name, mnf.Specs)
		case entities2.KeyStoreType:
			err = manifestStoreHandler.CreateKey(ctx, mnf.Name, mnf.Specs)
		case entities2.EthereumStoreType:
			err = manifestStoreHandler.CreateEthereum(ctx, mnf.Name, mnf.Specs)
		default:
			err = errors.InvalidFormatError("invalid store type")
		}

		if err != nil {
			return err
		}
	}

	manifestNodesHandler := nodesapi.NewNodesHandler(nodesService)
	for _, mnf := range manifests[entities.NodeKind] {
		err = manifestNodesHandler.Create(ctx, mnf.Specs)
		if err != nil {
			return err
		}
	}

	return nil
}
