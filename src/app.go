package src

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/app"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/server"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/pkg/log/zap"
	"github.com/consensysquorum/quorum-key-manager/src/manifests"
	manifestsmanager "github.com/consensysquorum/quorum-key-manager/src/manifests/manager"
	"github.com/consensysquorum/quorum-key-manager/src/middleware"
	"github.com/consensysquorum/quorum-key-manager/src/nodes"
	"github.com/consensysquorum/quorum-key-manager/src/stores"
)

type Config struct {
	HTTP      *server.Config
	Logger    *log.Config
	Manifests *manifestsmanager.Config
}

func New(cfg *Config) (*app.App, error) {
	logger, err := zap.NewLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	// Create app
	a := app.New(&app.Config{
		HTTP: cfg.HTTP,
	}, logger)

	// Register Service Configuration
	err := a.RegisterServiceConfig(cfg.Manifests)
	if err != nil {
		return nil, err
	}

	// Register Services
	err = manifests.RegisterService(a)
	if err != nil {
		return nil, err
	}

	err = stores.RegisterService(a)
	if err != nil {
		return nil, err
	}

	err = nodes.RegisterService(a)
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
