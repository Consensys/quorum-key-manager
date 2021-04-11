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
	// listener accepting HTTP connection
	// listener net.Listener

	// httpServer processing entrying HTTP request
	httpServer common.Runnable

	// backend managing core backend components
	backend core.Backend

	logger *log.Logger
}

func New(cfg *Config) *App {
	logger := log.NewLogger(cfg.Logger)
	bckend := core.New()
	httpServer := http.NewServer(cfg.HTTP, api.New(bckend), logger)

	return &App{
		cfg:        cfg,
		httpServer: httpServer,
		backend:    bckend,
		logger:     logger.SetComponent(Component),
	}
}

func (a App) Start(ctx context.Context) error {
	var storeMnfsts []*manifest.Manifest
	var nodeMnfsts []*manifest.Manifest

	for _, mnfst := range a.cfg.Manifests {
		switch mnfst.Kind {
		case "Node":
			nodeMnfsts = append(nodeMnfsts, mnfst)
		default:
			storeMnfsts = append(storeMnfsts, mnfst)
		}
	}

	if err := a.backend.StoreManager().Load(ctx, storeMnfsts...); err != nil {
		return err
	}

	if err := a.backend.NodeManager().Load(ctx, nodeMnfsts...); err != nil {
		return err
	}

	a.logger.Info("starting application")
	return a.httpServer.Start(ctx)
}

func (a App) Stop(ctx context.Context) error {
	a.logger.Info("stopping application")
	err := a.httpServer.Stop(ctx)
	return err
}

func (a App) Close() error {
	return a.httpServer.Close()
}

func (a App) Error() error {
	return a.httpServer.Error()
}
