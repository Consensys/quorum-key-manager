package http

import (
	"encoding/json"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"net/http"

	auth "github.com/consensys/quorum-key-manager/src/auth/api/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/gorilla/mux"
)

type SecretsHandler struct {
	stores stores.Stores
}

func NewSecretsHandler(storesConnector stores.Stores) *SecretsHandler {
	return &SecretsHandler{
		stores: storesConnector,
	}
}

func (h *SecretsHandler) Register(r *mux.Router) {
	r.Methods(http.MethodDelete).Path("/{id}/destroy").HandlerFunc(h.destroy)
	r.Methods(http.MethodPut).Path("/{id}/restore").HandlerFunc(h.restore)
	r.Methods(http.MethodPost).Path("/{id}").HandlerFunc(h.set)
	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.getOne)
	r.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(h.delete)
}

// @Summary Create a secret
// @Description Create new secret on selected Store
// @Tags Secrets
// @Accept json
// @Produce json
// @Param id path string true "Secret ID"
// @Param storeName path string true "Store ID"
// @Param request body types.SetSecretRequest true "Create Secret request"
// @Success 200 {object} types.SecretResponse "Secret data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id} [post]
func (h *SecretsHandler) set(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]
	setSecretRequest := &types.SetSecretRequest{}
	err := jsonutils.UnmarshalBody(request.Body, setSecretRequest)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	secretStore, err := h.stores.Secret(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	secret, err := secretStore.Set(ctx, id, setSecretRequest.Value, &entities.Attributes{
		Tags: setSecretRequest.Tags,
	})
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatSecretResponse(secret))
}

// @Summary Get a secret by id
// @Description Retrieve secret information by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store ID"
// @Param id path string true "Secret ID"
// @Param version query string false "secret version"
// @Param deleted query bool false "filter by only deleted accounts"
// @Success 200 {object} types.SecretResponse "Secret object"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id} [get]
func (h *SecretsHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]

	secretStore, err := h.stores.Secret(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	var secret *entities.Secret
	getDeleted := request.URL.Query().Get("deleted")
	if getDeleted == "" {
		version := request.URL.Query().Get("version")
		secret, err = secretStore.Get(ctx, id, version)
	} else {
		secret, err = secretStore.GetDeleted(ctx, id)
	}

	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatSecretResponse(secret))
}

// @Summary List secrets
// @Description List of secrets ids allocated in the targeted Store
// @Tags Secrets
// @Accept json
// @Produce json
// @Param deleted query bool false "filter by deleted accounts"
// @Param storeName path string true "Store ID"
// @Param limit query int false "page size"
// @Param page query int false "page number"
// @Success 200 {array} PageResponse "List of Secret IDs"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets [get]
func (h *SecretsHandler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	secretStore, err := h.stores.Secret(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	limit, offset, err := getLimitOffset(request)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	var ids []string
	getDeleted := request.URL.Query().Get("deleted")
	if getDeleted == "" {
		ids, err = secretStore.List(ctx, limit, offset)
	} else {
		ids, err = secretStore.ListDeleted(ctx, limit, offset)
	}
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = http2.WritePagingResponse(rw, request, ids)
}

// @Summary Delete a secret by id
// @Description Soft delete secret by id. It can be recovered
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store ID"
// @Param id path string true "Secret ID"
// @Success 204 "Deleted successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id} [delete]
func (h *SecretsHandler) delete(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]

	secretStore, err := h.stores.Secret(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = secretStore.Delete(ctx, id)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Destroy a secret by ID
// @Description Permanently delete a secret by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Secret ID"
// @Param id path string true "Key ID"
// @Success 204 "Destroyed successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id}/destroy [delete]
func (h *SecretsHandler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	id := mux.Vars(request)["id"]
	secretStore, err := h.stores.Secret(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = secretStore.Destroy(ctx, id)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Restore a soft-deleted secret
// @Description Restore a previously soft-deleted secret by ID
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store ID"
// @Param id path string true "Secret ID"
// @Success 204 "Restored successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id}/restore [put]
func (h *SecretsHandler) restore(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	id := mux.Vars(request)["id"]

	secretStore, err := h.stores.Secret(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = secretStore.Restore(ctx, id)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
