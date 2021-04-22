package handlers

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

type SecretsHandler struct {
	backend core.Backend
}

// New creates a http.Handler to be served on /secrets
func NewSecretsHandler(backend core.Backend) *mux.Router {
	h := &SecretsHandler{
		backend: backend,
	}

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/test").HandlerFunc(h.testRoute)

	return router
}

func (c *SecretsHandler) testRoute(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte("OK"))
}
