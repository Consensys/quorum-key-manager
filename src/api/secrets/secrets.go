package secretsapi

import (
	"encoding/json"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/types"
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
	router.Methods(http.MethodGet).Path("/").HandlerFunc(h.listRoute)
	router.Methods(http.MethodGet).Path("/test").HandlerFunc(h.testRoute)

	return router
}

func (c *handler) listRoute(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	storeList, err := c.bckend.StoreManager().List(ctx, types.HashicorpSecrets)
	if err != nil {
		_, _ = rw.Write([]byte(err.Error()))
		return
	}
	
	rawList, _ := json.Marshal(storeList)
	_, _ = rw.Write(rawList)
}

func (c *handler) testRoute(rw http.ResponseWriter, _ *http.Request) {
	_, _ = rw.Write([]byte("OK"))
}
