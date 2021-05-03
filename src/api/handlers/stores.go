package handlers

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

const StoreIDHeader = "X-Store-Id"

type StoresHandler struct {
	backend core.Backend
}

// New creates a http.Handler to be served on /stores
func NewStoresHandler(backend core.Backend) *mux.Router {
	h := &StoresHandler{
		backend: backend,
	}

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/test").HandlerFunc(h.testRoute)

	return router
}

func (c *StoresHandler) testRoute(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte("OK"))
}

func getStoreName(request *http.Request) string {
	return request.Header.Get(StoreIDHeader)
}
