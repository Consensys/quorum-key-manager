package src

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/pkg/http/server"
	"github.com/consensys/quorum-key-manager/pkg/log"
	"github.com/consensys/quorum-key-manager/src/manifests"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/consensys/quorum-key-manager/src/middleware"
	"github.com/consensys/quorum-key-manager/src/nodes"
	"github.com/consensys/quorum-key-manager/src/stores"
	postgresclient "github.com/consensys/quorum-key-manager/src/stores/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/src/stores/store/database/postgres"
)

type Config struct {
	HTTP      *server.Config
	Logger    *log.Config
	Manifests *manifestsmanager.Config
	Postgres  *postgresclient.Config
}

func New(cfg *Config, logger log.Logger) (*app.App, error) {
	// Create app
	a := app.New(&app.Config{HTTP: cfg.HTTP}, logger.WithComponent("app"))

	// No need to differentiate between DBs as we always use Postgres
	postgresClient, err := postgresclient.NewClient(cfg.Postgres)
	if err != nil {
		return nil, err
	}
	db := postgres.New(logger, postgresClient)

	// Register Service Configuration
	err = a.RegisterServiceConfig(cfg.Manifests)
	if err != nil {
		return nil, err
	}

	// Register Services
	err = manifests.RegisterService(a, logger.WithComponent("manifests"))
	if err != nil {
		return nil, err
	}

	err = stores.RegisterService(a, logger.WithComponent("stores"), db)
	if err != nil {
		return nil, err
	}

	err = nodes.RegisterService(a, logger.WithComponent("nodes"))
	if err != nil {
		return nil, err
	}

	// Set Middleware
	err = a.SetMiddleware(middleware.AccessLog(logger.WithComponent("accesslog")))
	if err != nil {
		return nil, err
	}

	return a, nil
}
