package app

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	postgresinfra "github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/consensys/quorum-key-manager/src/stores/api/http"
	"github.com/consensys/quorum-key-manager/src/stores/api/manifest"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/stores"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/utils"
	"github.com/consensys/quorum-key-manager/src/stores/connectors/vaults"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
)

func RegisterService(
	ctx context.Context,
	a *app.App,
	logger log.Logger,
	postgresClient postgresinfra.Client,
	manifests map[string][]entities.Manifest,
) (*stores.Connector, error) {
	// Data layer
	db := postgres.New(logger, postgresClient)

	// Business layer
	authManager := new(auth.Manager)
	err := a.Service(authManager)
	if err != nil {
		return nil, err
	}

	vaultsConnector := vaults.NewConnector(logger)
	storesConnector := stores.NewConnector(*authManager, db, vaultsConnector, logger)
	utilsConnector := utils.NewConnector(logger)

	// Service layer
	router := a.Router()
	http.NewStoresHandler(storesConnector).Register(router)
	http.NewUtilsHandler(utilsConnector).Register(router)

	manifestVaultHandler := manifest.NewVaultsHandler(vaultsConnector) // Manifest reading is synchronous, similar to a config file
	for _, mnf := range manifests[entities.VaultKind] {
		err = manifestVaultHandler.Register(ctx, mnf)
		if err != nil {
			return nil, err
		}
	}

	manifestStoreHandler := manifest.NewStoresHandler(storesConnector) // Manifest reading is synchronous, similar to a config file
	for _, mnf := range manifests[entities.StoreKind] {
		err = manifestStoreHandler.Register(ctx, mnf)
		if err != nil {
			return nil, err
		}
	}

	return storesConnector, nil
}
