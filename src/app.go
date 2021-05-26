package src

import (
	"context"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/server"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services"
	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/types"
	nodemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/nodes/manager"
)

const Component = "app"

// App is the main Key Manager application object
type App struct {
	cfg *Config

	// server processing entrying HTTP request
	server *http.Server

	// backend managing core business components
	backend core.Backend

	// logger logger object
	logger *log.Logger

	// Manifest local filesystem loader
	mnfstsLoader *manifest.LocalLoader

	mnfstsMsgs <-chan []manifest.Message
}

func New(cfg *Config, logger *log.Logger) (*App, error) {
	backend := core.New()

	httpServer := server.New(cfg.HTTP)
	httpServer.Handler = api.New(backend)

	mnfstsLoader, err := manifest.NewLocalLoader(cfg.ManifestPath)
	if err != nil {
		return nil, err
	}

	msgs := make(chan []manifest.Message, 1)
	// @TODO Implement unsubscribe and error handling
	_, err = mnfstsLoader.Subscribe(msgs)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:          cfg,
		server:       httpServer,
		backend:      backend,
		logger:       logger.SetComponent(Component),
		mnfstsLoader: mnfstsLoader,
		mnfstsMsgs:   msgs,
	}, nil
}

func (app *App) startServer(ctx context.Context) error {
	var cerr = make(chan error, 1)
	defer close(cerr)

	go func() {
		app.logger.WithField("addr", app.server.Addr).Info("started server")
		cerr <- app.server.ListenAndServe()
	}()

	select {
	case err := <-cerr:
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (app *App) stopServer(ctx context.Context) error {
	app.logger.Info("shutting down server")
	return app.server.Shutdown(ctx)
}

func (app *App) closeServer() error {
	return app.server.Close()
}

func (app App) Start(ctx context.Context) error {
	cerr := make(chan error, 1)
	defer close(cerr)

	go func(cerr chan error) {
		for _, msg := range <-app.mnfstsMsgs {
			if msg.Err != nil {
				app.logger.WithError(msg.Err).Warn("failed to read manifest")
				continue
			}

			switch msg.Manifest.Kind {
			case nodemanager.NodeKind:
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
		cerr <- app.startServer(ctx)
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
	err := app.stopServer(ctx)
	return err
}

func (app App) Close() error {
	return app.closeServer()
}

func (app App) Error() error {
	return nil
}
