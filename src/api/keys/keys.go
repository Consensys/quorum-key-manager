package keysapi

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

const RoutePrefix = "/keys"

type handler struct {
	bckend core.Backend
}

// New creates a http.Handler to be served on /secrets
func New(bckend core.Backend) http.Handler {
	h := &handler{
		bckend: bckend,
	}
	
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path(RoutePrefix).HandlerFunc(h.testRoute)
	
	return router
}

// @Summary Test keys route
// @Description Test keys route
// @Success 200 string
// @Router /keys [get]
func (c *handler) testRoute(rw http.ResponseWriter, _ *http.Request) {
	rw.Write([]byte("OK"))
}
