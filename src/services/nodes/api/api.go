package api

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/middleware"
	nodesmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/nodes/manager"
	"github.com/gorilla/mux"
)

const nodesPrefix = "/nodes"

type NodesAPI struct {
	nodes nodesmanager.Manager
}

// New creates a http.Handler to be served on JSON-RPC
func New(mngr nodesmanager.Manager) *NodesAPI {
	return &NodesAPI{
		nodes: mngr,
	}
}

func (h *NodesAPI) Register(router *mux.Router) {
	subrouter := router.PathPrefix(nodesPrefix).Subrouter()
	subrouter.Use(middleware.StripPrefix(nodesPrefix))

	subrouter.Methods(http.MethodPost).Path("/{id}").HandlerFunc(h.serveHTTPDownstream)
}

func (h *NodesAPI) serveHTTPDownstream(rw http.ResponseWriter, req *http.Request) {
	nodeID := mux.Vars(req)["id"]
	n, err := h.nodes.Node(req.Context(), nodeID)
	if err != nil {
		http.NotFound(rw, req)
		return
	}

	req2 := req.Clone(req.Context())
	req2.RequestURI = "/"
	n.ServeHTTP(rw, req2)
}
