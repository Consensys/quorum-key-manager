package secretsapi

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

type handler struct {
	bckend core.Backend
}

// New creates a http.Handler to be served on /secrets
func New(bckend core.Backend) http.Handler {
	h := &handler{
		bckend: bckend,
	}

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/test").HandlerFunc(h.testRoute)

	return router
}

func (c *handler) testRoute(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte("OK"))
}
