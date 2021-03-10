package src

import (
	"context"
	"net"
	"os"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
)

// App is the main Key Manager application object
type App struct {
	cfg *Config
	// listener accepting HTTP connection
	listener net.Listener

	// httpServer processing entrying HTTP request
	httpServer common.Runnable

	// backend managing core backend components
	backend core.Backend

	logger *log.Logger
}

func New(cfg *Config) common.Runnable {
	bckend := core.New()
	httpServer := api.New(cfg.HTTP, bckend)
	logger := log.NewLogger(cfg.Logger)

	return &App{
		cfg:        cfg,
		httpServer: httpServer,
		backend:    bckend,
		logger:     logger,
	}
}

func (a App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	ctx = log.With(ctx, a.logger)

	sig := common.NewSignalListener(func(sig os.Signal) {
		a.logger.WithField("sig", sig.String()).Warn("signal intercepted")
		cancel()
	})

	defer sig.Close()

	var cerr = make(chan error, 1)
	defer close(cerr)

	go func() {
		cerr <- a.httpServer.Start(ctx)
		cancel()
	}()
	
	select {
		case err := <-cerr:
			a.logger.WithError(err).Error("application exited with errors")
			return err
		case <-ctx.Done():
			a.logger.WithError(ctx.Err()).Info("application exited successfully")
	}

	return nil
}

func (a App) Stop(ctx context.Context) error {
	return a.httpServer.Stop(ctx)
}

func (a App) Close() error {
	panic("implement me")
}

func (a App) Error() error {
	panic("implement me")
}
