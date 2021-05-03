package src

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/http"
)

const Component = "app"

// App is the main Key Manager application object
type App struct {
	cfg *Config

	// httpServer processing entrying HTTP request
	httpServer common.Runnable

	// backend managing core business components
	backend core.Backend

	// logger logger object
	logger *log.Logger
}

func New(cfg *Config, logger *log.Logger) *App {
	backend := core.New()
	httpServer := http.NewServer(cfg.HTTP, api.New(backend), logger)

	return &App{
		cfg:        cfg,
		httpServer: httpServer,
		backend:    backend,
		logger:     logger.SetComponent(Component),
	}
}

func (app App) Start(ctx context.Context) error {
	var storeMnfsts []*manifest.Manifest
	var nodeMnfsts []*manifest.Manifest

	for _, mnfst := range app.cfg.Manifests {
		switch mnfst.Kind {
		case "Node":
			nodeMnfsts = append(nodeMnfsts, mnfst)
		default:
			storeMnfsts = append(storeMnfsts, mnfst)
		}
	}

	if err := app.backend.StoreManager().Load(log.With(ctx, app.logger), storeMnfsts...); err != nil {
		return err
	}

	if err := app.backend.NodeManager().Load(log.With(ctx, app.logger), nodeMnfsts...); err != nil {
		return err
	}

	app.logger.Info("starting application")
	return app.httpServer.Start(ctx)
}

func (app App) Stop(ctx context.Context) error {
	app.logger.Info("stopping application")
	err := app.httpServer.Stop(ctx)
	return err
}

func (app App) Close() error {
	return app.httpServer.Close()
}

func (app App) Error() error {
	return app.httpServer.Error()
}
