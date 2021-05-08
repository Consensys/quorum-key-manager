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

	// Manifest local filesystem loader 
	mnfstsLoader *manifest.LocalLoader

	mnfstsMsgs <-chan []*manifest.Message
}

func New(cfg *Config, logger *log.Logger) (*App, error) {
	backend := core.New()
	httpServer := http.NewServer(cfg.HTTP, api.New(backend), logger)

	mnfstsLoader, err := manifest.NewLocalLoader(cfg.ManifestPath)
	if err != nil {
		return nil, err
	}

	msgs := make(chan []*manifest.Message, 1)
	// @TODO Implement unsubscribe and error handling
	_, err = mnfstsLoader.Subscribe(msgs)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:          cfg,
		httpServer:   httpServer,
		backend:      backend,
		logger:       logger.SetComponent(Component),
		mnfstsLoader: mnfstsLoader,
		mnfstsMsgs:   msgs,
	}, nil
}

func (app App) Start(ctx context.Context) error {
	cerr := make(chan error, 1)
	defer close(cerr)

	go func(cerr chan error) {
		// @TODO Improve to read in batches
		for _, msg := range <-app.mnfstsMsgs {
			switch msg.Manifest.Kind {
			case "Node":
				if err := app.backend.NodeManager().Load(log.With(ctx, app.logger), msg.Manifest); err != nil {
					cerr <- err
				}
			default:
				if err := app.backend.StoreManager().Load(log.With(ctx, app.logger), msg.Manifest); err != nil {
					cerr <- err
				}
			}
		}
	}(cerr)

	go func(cerr chan error) {
		app.logger.Info("starting application")
		cerr <- app.httpServer.Start(ctx)
	}(cerr)
	
	if err := app.mnfstsLoader.Start(); err != nil {
		return err
	}

	select {
	case err := <-cerr:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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
