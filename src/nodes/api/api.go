package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/consensys/quorum-key-manager/src/auth/api/http_middlewares"

	nodesmanager "github.com/consensys/quorum-key-manager/src/nodes"

	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/gorilla/mux"
)

type NodesAPI struct {
	nodes nodesmanager.Service
}

// New creates a http.Handler to be served on JSON-RPC
func New(mngr nodesmanager.Service) *NodesAPI {
	return &NodesAPI{
		nodes: mngr,
	}
}

func (h *NodesAPI) Register(router *mux.Router) {
	subrouter := router.PathPrefix("/nodes/{nodeName}").Subrouter()
	subrouter.Use(stripNodePrefix)
	subrouter.PathPrefix("").HandlerFunc(h.serveHTTPDownstream)
}

func stripNodePrefix(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Trim prefix
		prefix := fmt.Sprintf("/nodes/%s", mux.Vars(r)["nodeName"])
		p := strings.TrimPrefix(r.URL.Path, prefix)
		if p == "" {
			p = "/"
		}

		rp := strings.TrimPrefix(r.URL.RawPath, prefix)
		if rp == "" {
			rp = "/"
		}

		uri := strings.TrimPrefix(r.RequestURI, prefix)
		if uri == "" {
			uri = "/"
		}

		// Create request to be updated
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		r2.RequestURI = uri

		// Serve next handler
		h.ServeHTTP(w, r2)
	})
}

func (h *NodesAPI) serveHTTPDownstream(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	nodeName := mux.Vars(req)["nodeName"]

	n, err := h.nodes.Node(req.Context(), nodeName, http_middlewares.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	n.ServeHTTP(rw, req)
}
