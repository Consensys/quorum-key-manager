package accesslog

import (
	"github.com/gorilla/handlers"
	"io"
	"net/http"
)

type Middleware struct {
	logger io.Writer
}

func NewMiddleware(accessLogger io.Writer) *Middleware {
	return &Middleware{
		logger: accessLogger,
	}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(m.logger, next)
}
