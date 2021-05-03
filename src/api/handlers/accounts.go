package handlers

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

type AccountsHandler struct {
	backend core.Backend
}

// New creates a http.Handler to be served on /accounts
func NewAccountsHandler(backend core.Backend) *mux.Router {
	h := &AccountsHandler{
		backend: backend,
	}

	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/test").HandlerFunc(h.testRoute)

	return router
}

func (h *AccountsHandler) testRoute(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte("OK"))
}
