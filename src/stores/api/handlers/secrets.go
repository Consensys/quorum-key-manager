package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	storesmanager "github.com/consensys/quorum-key-manager/src/stores/manager"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
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
	r.Methods(http.MethodGet).Path("/deleted").HandlerFunc(h.listDeleted)
	r.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.getOne)
	r.Methods(http.MethodGet).Path("/{id}/deleted").HandlerFunc(h.getDeletedOne)
	r.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(h.delete)
	r.Methods(http.MethodDelete).Path("/{id}/destroy").HandlerFunc(h.destroy)
	r.Methods(http.MethodPut).Path("/{id}/restore").HandlerFunc(h.restore)
}

// @Summary Creates a secret
// @Description Create new secret on selected Store
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

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
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

// @Summary Gets a secret by id
// @Description Retrieves secret information by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Secret Identifier"
// @Success 200 {object} types.SecretResponse "Secret object"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id} [get]
func (h *SecretsHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]
	version := request.URL.Query().Get("version")

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
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

// @Summary Gets a deleted secret by id
// @Description Retrieves deleted secret information by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Secret Identifier"
// @Success 200 {object} types.SecretResponse "Secret object"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id}/deleted [get]
func (h *SecretsHandler) getDeletedOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	secret, err := secretStore.GetDeleted(ctx, id)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatSecretResponse(secret))
}

// @Summary List secrets
// @Description List of secret IDs stored in the selected Store
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

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
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

// @Summary List soft-deleted secrets
// @Description List of deleted secret IDs stored in the selected Store
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Success 200 {array} []types.SecretResponse "List of Secret IDs"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/deleted [get]
func (h *SecretsHandler) listDeleted(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	ids, err := secretStore.ListDeleted(ctx)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(ids)
}

// @Summary Deletes a secret by id
// @Description Soft delete secret by id. It can be recovered
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Secret Identifier"
// @Success 204 "Deleted successfully"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id} [delete]
func (h *SecretsHandler) delete(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = secretStore.Delete(ctx, id)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Destroys a secret by ID
// @Description Permanently deletes a secret by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Secret Identifier"
// @Param id path string true "Key identifier"
// @Success 204 "Destroyed successfully"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id}/destroy [delete]
func (h *SecretsHandler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = keyStore.Destroy(ctx, getID(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Restores a soft-deleted secret
// @Description Restores a previously soft-deleted secret by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Secret identifier"
// @Success 204 "Restored successfully"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id}/restore [put]
func (h *SecretsHandler) restore(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = keyStore.Undelete(ctx, getID(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
