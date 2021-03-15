package storesapi

import (
	"net/http"

	pkghttp "github.com/ConsenSysQuorum/quorum-key-manager/pkg/http"
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

	return router
}

func (c *handler) listRoute(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	storeKind := req.URL.Query().Get("kind")
	storeList, err := c.bckend.StoreManager().List(ctx, types.Kind(storeKind))
	if err != nil {
		pkghttp.WriteErrorResponse(rw, err)
		return
	}

	pkghttp.WriteJSONResponse(rw, storeList)
}
