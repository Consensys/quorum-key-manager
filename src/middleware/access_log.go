package middleware

import (
	"net/http"

	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/gorilla/handlers"
)

func AccessLog(cfg *log.Config) func(handlers http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logger := log.NewLogger(cfg).SetComponent("accesslog")
		return handlers.LoggingHandler(logger, h)
	}
}
