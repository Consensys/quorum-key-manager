package handlers

import (
	"encoding/json"
	formatters2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/formatters"
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	jsonutils "github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/gorilla/mux"
)

type SecretsHandler struct {
	stores storemanager.Manager
}

// New creates a http.Handler to be served on /secrets
func NewSecretsHandler(s storemanager.Manager) *SecretsHandler {
	return &SecretsHandler{
		stores: s,
	}
}

func (h *SecretsHandler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("").HandlerFunc(h.set)
	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.getOne)
}

func (h *SecretsHandler) set(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	setSecretRequest := &types2.SetSecretRequest{}
	err := jsonutils.UnmarshalBody(request.Body, setSecretRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	secret, err := secretStore.Set(ctx, setSecretRequest.ID, setSecretRequest.Value, &entities2.Attributes{
		Tags: setSecretRequest.Tags,
	})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters2.FormatSecretResponse(secret))
}

func (h *SecretsHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]
	version := request.URL.Query().Get("version")

	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	secret, err := secretStore.Get(ctx, id, version)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters2.FormatSecretResponse(secret))
}

func (h *SecretsHandler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx))
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
