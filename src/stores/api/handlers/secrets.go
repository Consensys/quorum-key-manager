package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensysquorum/quorum-key-manager/pkg/json"
	"github.com/consensysquorum/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensysquorum/quorum-key-manager/src/stores/api/types"
	storesmanager "github.com/consensysquorum/quorum-key-manager/src/stores/manager"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	"github.com/gorilla/mux"
)

type SecretsHandler struct {
	stores storesmanager.Manager
}

// NewSecretsHandler creates a http.Handler to be served on /secrets
func NewSecretsHandler(s storesmanager.Manager) *SecretsHandler {
	return &SecretsHandler{
		stores: s,
	}
}

func (h *SecretsHandler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/{id}").HandlerFunc(h.set)
	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.getOne)
}

// @Summary Create Secret
// @Description Create new Secret on selected Store
// @Tags Secrets
// @Accept json
// @Produce json
// @Param id path string true "Secret Identifier"
// @Param storeName path string true "Store Identifier"
// @Param request body types.SetSecretRequest true "Create Secret request"
// @Success 200 {object} types.SecretResponse "Secret data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id} [post]
func (h *SecretsHandler) set(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	setSecretRequest := &types.SetSecretRequest{}
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

	secret, err := secretStore.Set(ctx, mux.Vars(request)["id"], setSecretRequest.Value, &entities.Attributes{
		Tags: setSecretRequest.Tags,
	})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatSecretResponse(secret))
}

// @Summary Get secret by id
// @Description Retrieve secret information by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Secret ID"
// @Success 200 {object} types.SecretResponse "Secret object"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id} [get]
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

	_ = json.NewEncoder(rw).Encode(formatters.FormatSecretResponse(secret))
}

// @Summary List Secrets
// @Description List of Secret IDs stored in the selected Store
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Success 200 {array} []types.SecretResponse "List of Secret IDs"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets [get]
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
