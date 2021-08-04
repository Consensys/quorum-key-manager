package handlers

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	"github.com/gorilla/mux"
)

type AliasHandler struct {
	alias aliases.Alias
}

func NewAliasHandler(alias aliases.Alias) *AliasHandler {
	return &AliasHandler{
		alias: alias,
	}
}

func (h *AliasHandler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/registry/{name}/alias").HandlerFunc(h.create)
	// TODO the: register remaining handling methods
}

func (h *AliasHandler) create(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := req.Context()

	createAliasRequest := new(types.CreateAliasRequest)
	keyStore, err := h.alias.CreateAlias(
		ctx,
		// TODO the: retrieve registry name from gorilla mux variable,
		// TOOO the: retrieve alias key and value from createAliasRequest
	)
	if err != nil {
		// TODO the: WriteHTTPErrorResponse
		// WriteHTTPErrorResponse(rw, err)
		return
	}

	
}

// TODO the: implement remaining handling methods
