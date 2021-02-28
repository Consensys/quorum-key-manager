package app

import (
	"net"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/backend"
)

// App is the main Key Manager application object
type App struct {
	// listener accepting HTTP connection
	listener net.Listener

	// httpServer processing entrying HTTP request
	httpServer *http.Server

	// backend managing core backend components
	backend backend.Backend
}
