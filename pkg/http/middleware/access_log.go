package middleware

import (
	"io"
	"net/http"

	"github.com/gorilla/handlers"
)

func AccessLog(logger io.Writer) func(handlers http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(logger, h)
	}
}
