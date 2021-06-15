package middleware

import (
	"net/http"

	"github.com/consensysquorum/quorum-key-manager/pkg/log-old"
	"github.com/gorilla/handlers"
)

func AccessLog(cfg *log_old.Config) func(handlers http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logger := log_old.NewLogger(cfg).SetComponent("accesslog")
		return handlers.LoggingHandler(logger, h)
	}
}
