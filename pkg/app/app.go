package app

import (
	"context"
	"net/http"
	"reflect"
	"sync"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/http/server"
	gorillamux "github.com/gorilla/mux"
)

const (
	initializingState = iota
	runningState
	stoppingState
	closedState

	pointerErrMessage = "cannot attach service to a non pointer"
)

// App is the main Key Manager application object
type App struct {
	cfg *Config

	// logger logger object
	logger log.Logger

	// server processing HTTP requests
	server *http.Server
	// server processing HTTPS requests
	tlsServer *http.Server
	// server processing HTTP, health related requests
	healthz *http.Server
	router  *gorillamux.Router

	// middleware applied before routing
	middleware func(http.Handler) http.Handler

	// Services attached to the app
	mux            sync.Mutex
	services       []reflect.Value
	serviceConfigs map[reflect.Type]reflect.Value

	state  int // tracks state of app
	errors chan error
}

func New(cfg *Config, logger log.Logger) *App {
	// Create router and register APIs
	router := gorillamux.NewRouter()

	// Create API server
	apiServer := server.New(cfg.HTTP)
	apiServer.Handler = router

	apiTLSServer := server.NewTLS(cfg.HTTP)

	apiTLSServer.Handler = router

	// Create Healthz server
	healthzServer := server.NewHealthz(cfg.HTTP)
	healthzServer.Handler = server.NewHealthzHandler()

	return &App{
		cfg:            cfg,
		logger:         logger,
		server:         apiServer,
		tlsServer:      apiTLSServer,
		healthz:        healthzServer,
		errors:         make(chan error),
		router:         router,
		serviceConfigs: make(map[reflect.Type]reflect.Value),
	}
}

func (app *App) SetMiddleware(mid func(http.Handler) http.Handler) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	if app.state != initializingState {
		errMessage := "can't register middleware on running or stopped app"
		app.logger.Error(errMessage)
		return errors.ConfigError(errMessage)
	}

	app.middleware = mid

	return nil
}

// RegisterServiceConfig register a config
// cfg MUST be a pointer to a struct
func (app *App) RegisterServiceConfig(cfg interface{}) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	if app.state != initializingState {
		errMessage := "can't register config on running or stopped app"
		app.logger.Error(errMessage)
		return errors.ConfigError(errMessage)
	}

	if app.hasConfig(cfg) {
		errMessage := "attempt to register config %T more than once"
		app.logger.Error(errMessage)
		return errors.ConfigError(errMessage)
	}

	cfgV := reflect.ValueOf(cfg)

	if !(cfgV.IsValid() && cfgV.Type().Kind() == reflect.Ptr && cfgV.Type().Elem().Kind() == reflect.Struct) {
		errMessage := "failed to extract config"
		app.logger.Error(errMessage, "config", cfg)
		return errors.ConfigError(errMessage)
	}

	app.serviceConfigs[cfgV.Type()] = cfgV

	return nil
}

func (app *App) hasConfig(i interface{}) bool {
	if _, ok := app.serviceConfigs[reflect.TypeOf(i)]; ok {
		return true
	}

	return false
}

// ServiceConfig loads a service configuration into cfg
// It expects a pointer to a struct and then sets its value
func (app *App) ServiceConfig(cfg interface{}) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	cfgV := reflect.ValueOf(cfg)
	if cfgV.Type().Kind() != reflect.Ptr {
		app.logger.Error(pointerErrMessage)
		return errors.ConfigError(pointerErrMessage)
	}

	for typ, config := range app.serviceConfigs {
		if typ == cfgV.Type() {
			if !config.IsZero() {
				cfgV.Elem().Set(config.Elem())
			}
			return nil
		}
	}

	return errors.ConfigError("unknown config")
}

func (app *App) RegisterService(srv interface{}) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	if app.state != initializingState {
		errMessage := "cannot register service on running or stopped app"
		app.logger.Error(errMessage)
		return errors.ConfigError(errMessage)
	}

	if rSrv, ok := srv.(common.Runnable); ok {
		app.services = append(app.services, reflect.ValueOf(rSrv))
	} else {
		errMessage := "registered service is not a runnable"
		app.logger.Error(errMessage)
		return errors.ConfigError(errMessage)
	}

	if hlzSrv, ok := srv.(common.Checkable); ok {
		if healthz, ok2 := app.healthz.Handler.(*server.HealthzHandler); ok2 {
			healthz.AddLivenessCheck(hlzSrv.ID(), hlzSrv.CheckLiveness)
			healthz.AddReadinessCheck(hlzSrv.ID(), hlzSrv.CheckReadiness)
		}
	}

	return nil
}

// Service loads a service into srv
// It expects a pointer and then sets its value
func (app *App) Service(srv interface{}) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	// Finds the service to return
	srvV := reflect.ValueOf(srv)

	if srvV.Type().Kind() != reflect.Ptr {
		app.logger.Error(pointerErrMessage)
		return errors.ConfigError(pointerErrMessage)
	}

	for _, regSrv := range app.services {
		typ := regSrv.Type()
		switch {
		case typ == srvV.Type():
			srvV.Elem().Set(regSrv.Elem())
			return nil
		case typ == srvV.Type().Elem():
			srvV.Elem().Set(regSrv)
			return nil
		case srvV.Elem().Type().Kind() == reflect.Interface && typ.Implements(srvV.Elem().Type()):
			srvV.Elem().Set(regSrv)
			return nil
		}
	}

	return errors.ConfigError("unknown service")
}

func (app *App) Router() *gorillamux.Router {
	return app.router
}

func (app *App) startServer() {
	app.logger.Debug("starting app server...")

	// Wrap handler into middleware
	if app.middleware != nil {
		app.server.Handler = app.middleware(app.server.Handler)
		app.tlsServer.Handler = app.middleware(app.tlsServer.Handler)
	}

	go func() {
		app.logger.Info("starting API server", "addr", app.server.Addr)
		apiErr := app.server.ListenAndServe()
		if apiErr == nil {
			app.logger.Debug("started API server successfully", "addr", app.server.Addr)
		} else {
			app.logger.WithError(apiErr).Info("failed to start API server")
		}
		app.errors <- apiErr
	}()

	go func() {
		app.logger.Info("starting API TLS server", "addr", app.tlsServer.Addr, "cert", app.cfg.HTTP.TLSCert, "key", app.cfg.HTTP.TLSKey)
		tlsErr := app.tlsServer.ListenAndServeTLS(app.cfg.HTTP.TLSCert, app.cfg.HTTP.TLSKey)
		if tlsErr == nil {
			app.logger.Debug("started API TLS server successfully", "addr", app.server.Addr)
		} else {
			app.logger.WithError(tlsErr).Error("failed to start API TLS server")
		}
		app.errors <- tlsErr
	}()

	go func() {
		app.logger.Info("started Health server", "addr", app.healthz.Addr)
		healthErr := app.healthz.ListenAndServe()
		if healthErr == nil {
			app.logger.Debug("started Health server successfully", "addr", app.server.Addr)
		} else {
			app.logger.WithError(healthErr).Error("failed to start Health server")
		}
		app.errors <- healthErr
	}()

	app.logger.Debug("servers (API, TLS and Health) have started")
}

func (app *App) stopServer(ctx context.Context) error {
	app.logger.Debug("shutting down app server...")
	if err := app.healthz.Shutdown(ctx); err != nil {
		app.logger.WithError(err).Error("health server could not shut down")
		return err
	}

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.WithError(err).Error("http api server could not shut down")
		return err
	}

	if err := app.tlsServer.Shutdown(ctx); err != nil {
		app.logger.WithError(err).Error("tls api server could not shut down")
		return err
	}

	app.logger.Info("servers (API, TLS and Health) gracefully shut down")
	return nil
}

func (app *App) closeServer() error {
	return app.server.Close()
}

func (app *App) Start(ctx context.Context) error {
	app.logger.Debug("starting application...")
	app.state = runningState

	app.startServer()

	// Start all registered services.
	var err error
	var started []common.Runnable
	for i := len(app.services) - 1; i >= 0; i-- {
		srv := app.services[i].Interface().(common.Runnable)
		if err = srv.Start(ctx); err != nil {
			break
		}
		started = append(started, srv)
	}

	// If a service failed to start then stop other services
	if err != nil {
		app.state = stoppingState
		for i := len(started) - 1; i >= 0; i-- {
			_ = started[i].Stop(ctx)
		}
		_ = app.stopServer(ctx)
	}

	app.logger.Info("application started")
	return err
}

func (app *App) Stop(ctx context.Context) error {
	app.logger.Debug("stopping application...")
	app.state = stoppingState

	var err error
	for i := len(app.services) - 1; i >= 0; i-- {
		srv := app.services[i].Interface().(common.Runnable)
		if srvErr := srv.Stop(ctx); srvErr != nil {
			err = srvErr
		}
	}

	httpErr := app.stopServer(ctx)
	if httpErr != nil {
		err = httpErr
	}

	app.logger.Info("application stopped")
	return err
}

func (app *App) Errors() <-chan error {
	return app.errors
}

func (app *App) Close() error {
	app.state = closedState
	close(app.errors)
	return app.closeServer()
}
