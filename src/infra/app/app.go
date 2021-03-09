package app

import (
	"net"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
)

// App is the main Key Manager application object
type App struct {
	// listener accepting HTTP connection
	listener net.Listener

	// httpServer processing entrying HTTP request
	httpServer *http.Server

	// backend managing core backend components
	backend core.Backend
}

// @TODO Inject http server
func New() *App {
	bckend := core.New()

	server := &http.Server{
		Handler: api.New(bckend),
	}

	return &App{
		httpServer: server,
		backend:    bckend,
	}
}
