package accesslog

import (
	httpinfra "github.com/consensys/quorum-key-manager/src/infra/http/middlewares"
	"github.com/gorilla/handlers"
	"io"
	"net/http"
)

type Middleware struct {
	logger io.Writer
}

var _ httpinfra.Middleware = &Middleware{}

func NewMiddleware(accessLogger io.Writer) *Middleware {
	return &Middleware{
		logger: accessLogger,
	}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(m.logger, next)
}
