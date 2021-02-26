package app

import (
	"net"
	"net/http"
)

type App struct {
	listener net.Listener

	httpServer *http.Server

	backend backend.Backend
}
