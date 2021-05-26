package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	jsonutils "github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/formatters"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/gorilla/mux"
)

type Eth1Handler struct {
	backend core.Backend
}

// New creates a http.Handler to be served on /accounts
func NewAccountsHandler(backend core.Backend) *mux.Router {
	h := &Eth1Handler{
		backend: backend,
	}

	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/").HandlerFunc(h.create)

	return router
}

func (h *Eth1Handler) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	createReq := &types.CreateEth1AccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, createReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	eth1Acc, err := eth1Store.Create(
		ctx,
		createReq.ID,
		&entities.Attributes{
			// @TODO Fill up attributes based on user request
		})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEth1AccResponse(eth1Acc))
}
