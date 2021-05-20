package handlers

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

type JSONRPCHandler struct {
	backend core.Backend
}

// New creates a http.Handler to be served on JSON-RPC
func NewJSONRPCHandler(bcknd core.Backend) *mux.Router {
	h := &JSONRPCHandler{
		backend: bcknd,
	}

	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/{id}").HandlerFunc(h.downstream)
	return router
}

func (h *JSONRPCHandler) downstream(rw http.ResponseWriter, req *http.Request) {
	nodeID := mux.Vars(req)["id"]
	n, err := h.backend.NodeManager().Node(req.Context(), nodeID)
	if err != nil {
		http.NotFound(rw, req)
		return
	}

	req2 := req.Clone(req.Context())
	req2.RequestURI = "/"
	n.ServeHTTP(rw, req2)
}
