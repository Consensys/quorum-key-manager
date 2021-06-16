package middleware

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	"net/http"

	"github.com/gorilla/handlers"
)

func AccessLog(logger log.Logger) func(handlers http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(logger, h)
	}
}
