package src

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/app"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/server"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	manifests2 "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/middleware"
	nodes2 "github.com/ConsenSysQuorum/quorum-key-manager/src/nodes"
	stores2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores"
)

type Config struct {
	HTTP      *server.Config
	Logger    *log.Config
	Manifests *manager.Config
}

func New(cfg *Config, logger *log.Logger) (*app.App, error) {
	// Create app
	a := app.New(&app.Config{HTTP: cfg.HTTP}, logger)

	// Register Service Configuration
	err := a.RegisterServiceConfig(cfg.Manifests)
	if err != nil {
		return nil, err
	}

	// Register Services
	err = manifests2.RegisterService(a)
	if err != nil {
		return nil, err
	}

	err = stores2.RegisterService(a)
	if err != nil {
		return nil, err
	}

	err = nodes2.RegisterService(a)
	if err != nil {
		return nil, err
	}

	// Set Middleware
	err = a.SetMiddleware(middleware.AccessLog(cfg.Logger))
	if err != nil {
		return nil, err
	}

	return a, nil
}
