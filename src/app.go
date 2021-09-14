package src

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/pkg/http/middleware"
	"github.com/consensys/quorum-key-manager/pkg/http/server"
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	"github.com/consensys/quorum-key-manager/src/manifests"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/consensys/quorum-key-manager/src/nodes"
	stores "github.com/consensys/quorum-key-manager/src/stores/app"
	"github.com/justinas/alice"
)

type Config struct {
	HTTP      *server.Config
	Logger    *log.Config
	Manifests *manifestsmanager.Config
	Postgres  *client.Config
	Auth      *auth.Config
}

func New(cfg *Config, logger log.Logger) (*app.App, error) {
	// Create app
	a := app.New(&app.Config{HTTP: cfg.HTTP}, logger.WithComponent("app"))

	// Register Service Configuration
	err := a.RegisterServiceConfig(cfg.Manifests)
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(cfg.Auth)
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(&stores.Config{Postgres: cfg.Postgres})
	if err != nil {
		return nil, err
	}

	// Register Services
	err = manifests.RegisterService(a, logger.WithComponent("manifests"))
	if err != nil {
		return nil, err
	}

	err = auth.RegisterService(a, logger.WithComponent("auth"))
	if err != nil {
		return nil, err
	}

	err = stores.RegisterService(a, logger.WithComponent("stores"))
	if err != nil {
		return nil, err
	}

	err = a.RegisterServiceConfig(&aliases.Config{Postgres: cfg.Postgres})
	if err != nil {
		return nil, err
	}

	err = aliases.RegisterService(a, logger.WithComponent("aliases"))
	if err != nil {
		return nil, err
	}

	err = nodes.RegisterService(a, logger.WithComponent("nodes"))
	if err != nil {
		return nil, err
	}

	// Set Middleware
	authmid, err := auth.Middleware(a, logger.WithComponent("auth-mid"))
	if err != nil {
		return nil, err
	}

	mid := alice.New(middleware.AccessLog(logger.WithComponent("accesslog")), authmid)

	err = a.SetMiddleware(mid.Then)
	if err != nil {
		return nil, err
	}

	return a, nil
}
