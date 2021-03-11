package middleware

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/handlers"
)

func New(_ core.Backend) func(handlers http.Handler) http.Handler {
	// TODO: implement the sequence of middlewares to apply before routing
	return func(h http.Handler) http.Handler {
		logger := log.NewLogger(&log.Config{
			Level:     log.InfoLevel,
			Timestamp: false,
		}).SetComponent("accesslog")
		return handlers.LoggingHandler(logger, h)
	}
}
