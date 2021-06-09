package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/server"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	gorillamux "github.com/gorilla/mux"
)

const Component = "app"

const (
	initializingState = iota
	runningState
	stoppingState
	closedState
)

var (
	ErrServiceUnknown = errors.New("unknown service")
	ErrConfigUnknown  = errors.New("unknown config")
)

// App is the main Key Manager application object
type App struct {
	cfg *Config

	// logger logger object
	logger *log.Logger

	// server processing entrying HTTP request
	server  *http.Server
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

func New(cfg *Config, logger *log.Logger) *App {
	// Create router and register APIs
	router := gorillamux.NewRouter()

	// Create API server
	apiServer := server.New(fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port), cfg.HTTP)
	apiServer.Handler = router

	// Create Healthz server
	healthzServer := server.New(fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.HealthzPort), cfg.HTTP)
	healthzServer.Handler = NewHealthzHandler()

	return &App{
		cfg:            cfg,
		logger:         logger.SetComponent(Component),
		server:         apiServer,
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
		return fmt.Errorf("can't register middleware on running or stopped app")
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
		return fmt.Errorf("can't register config on running or stopped app")
	}

	if app.hasConfig(cfg) {
		return fmt.Errorf("attempt to register config %T more than once", cfg)
	}

	cfgV := reflect.ValueOf(cfg)

	if !(cfgV.IsValid() && cfgV.Type().Kind() == reflect.Ptr && cfgV.Type().Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("attempt to register config %T which is not a pointer to struct", cfg)
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
		return fmt.Errorf("can not attach service to a non pointer")
	}

	for typ, config := range app.serviceConfigs {
		if typ == cfgV.Type() {
			cfgV.Elem().Set(config.Elem())
			return nil
		}
	}

	return ErrConfigUnknown
}

func (app *App) RegisterService(srv interface{}) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	if app.state != initializingState {
		return fmt.Errorf("can't register service on running or stopped app")
	}

	if rSrv, ok := srv.(common.Runnable); ok {
		app.services = append(app.services, reflect.ValueOf(rSrv))
	} else {
		return fmt.Errorf("register service is not a runnable")
	}

	if hlzSrv, ok := srv.(common.Checkable); ok {
		if healthz, ok2 := app.healthz.Handler.(*healthzHandler); ok2 {
			healthz.AddReadinessCheck(hlzSrv.ID(), hlzSrv.IsLive)
			healthz.AddReadinessCheck(hlzSrv.ID(), hlzSrv.IsReady)
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
		return fmt.Errorf("can not attach service to a non pointer")
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

	return ErrServiceUnknown
}

func (app *App) Router() *gorillamux.Router {
	return app.router
}

func (app *App) startServer() {
	app.logger.Info("starting server...")

	// Wrap handler into middleware
	if app.middleware != nil {
		app.server.Handler = app.middleware(app.server.Handler)
	}

	go func() {
		app.logger.WithField("addr", app.server.Addr).Info("started API server")
		app.errors <- app.server.ListenAndServe()
	}()
	
	go func() {
		app.logger.WithField("addr", app.healthz.Addr).Info("started Health server")
		app.errors <- app.healthz.ListenAndServe()
	}()
}

func (app *App) stopServer(ctx context.Context) error {
	app.logger.Info("shutting down server...")
	err := app.server.Shutdown(ctx)
	if err != nil {
		app.logger.WithError(err).Errorf("server could not shut down")
		return err
	}
	app.logger.Info("server gracefully shut down")
	return nil
}

func (app *App) closeServer() error {
	return app.server.Close()
}

func (app *App) Start(ctx context.Context) error {
	app.logger.Info("starting application...")
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

	return err
}

func (app *App) Stop(ctx context.Context) error {
	app.logger.Info("stopping application...")
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
