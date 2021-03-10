// nolint
package src

import (
	"context"
	"net"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
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
	logger := log.NewLogger()

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

	wg := &sync.WaitGroup{}
	wg.Add(1)
	var gerr error
	go func() {
		err := a.httpServer.Start(ctx)
		if err != nil {
			a.logger.WithError(err).Error("http server failed")
			gerr = errors.CombineErrors(gerr, err)
		}
		cancel()
		wg.Done()
	}()

	// Wait for all Envelope to have complete execution
	wg.Wait()
	if gerr != nil {
		a.logger.WithError(gerr).Error("application exited with errors")
	} else {
		a.logger.Info("application exited")
	}

	return gerr
}

func (a App) Stop(ctx context.Context) error {
	err := a.httpServer.Stop(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a App) Close() error {
	panic("implement me")
}

func (a App) Error() error {
	panic("implement me")
}
