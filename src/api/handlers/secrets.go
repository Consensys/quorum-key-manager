package handlers

import (
	"encoding/json"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	jsonutils "github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"

	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/formatters"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

type SecretsHandler struct {
	backend core.Backend
}

// New creates a http.Handler to be served on /secrets
func NewSecretsHandler(backend core.Backend) *mux.Router {
	h := &SecretsHandler{
		backend: backend,
	}

	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path("/").HandlerFunc(h.set)
	router.Methods(http.MethodGet).Path("/").HandlerFunc(h.list)
	router.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.getOne)

	return router
}

func (h *SecretsHandler) set(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	setSecretRequest := &types.SetSecretRequest{}
	err := jsonutils.UnmarshalBody(request.Body, setSecretRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	secretStore, err := h.backend.StoreManager().GetSecretStore(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	secret, err := secretStore.Set(ctx, setSecretRequest.ID, setSecretRequest.Value, &entities.Attributes{
		Tags: setSecretRequest.Tags,
	})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatSecretResponse(secret))
}

func (h *SecretsHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]
	version := request.URL.Query().Get("version")

	secretStore, err := h.backend.StoreManager().GetSecretStore(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	secret, err := secretStore.Get(ctx, id, version)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatSecretResponse(secret))
}

func (h *SecretsHandler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	secretStore, err := h.backend.StoreManager().GetSecretStore(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	ids, err := secretStore.List(ctx)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(ids)
}
