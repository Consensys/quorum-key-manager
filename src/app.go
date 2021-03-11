package src

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
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

func New(ctx context.Context, cfg *Config) *App {
	bckend := core.New()
	httpServer := http.NewServer(ctx, cfg.HTTP, api.New(bckend))

	return &App{
		cfg:        cfg,
		httpServer: httpServer,
		backend:    bckend,
		logger:     log.FromContext(ctx).SetComponent(Component),
	}
}

func (a App) Start(ctx context.Context) error {
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
