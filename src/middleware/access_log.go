package middleware

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"github.com/consensysquorum/quorum-key-manager/pkg/log/zap"
	"net/http"

	"github.com/gorilla/handlers"
)

func AccessLog(cfg *log.Config) func(handlers http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logger, _ := zap.NewLogger(cfg)
		logger.SetComponent("accesslog")
		return handlers.LoggingHandler(logger, h)
	}
}
