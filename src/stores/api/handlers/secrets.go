package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/gorilla/mux"
)

type SecretsHandler struct {
	stores stores.Manager
}

// NewSecretsHandler creates a http.Handler to be served on /secrets
func NewSecretsHandler(s stores.Manager) *SecretsHandler {
	return &SecretsHandler{
		stores: s,
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
// @Param id path string true "Secret Identifier"
// @Param storeName path string true "Store Identifier"
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

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
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
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Secret Identifier"
// @Param version query string false "secret version"
// @Param deleted query bool false "filter by deleted accounts"
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

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
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
// @Description List of secret IDs stored in the selected Store
// @Tags Secrets
// @Accept json
// @Produce json
// @Param deleted query bool false "filter by deleted accounts"
// @Param storeName path string true "Store Identifier"
// @Success 200 {array} []types.SecretResponse "List of Secret IDs"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets [get]
func (h *SecretsHandler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	var ids []string
	getDeleted := request.URL.Query().Get("deleted")
	limit, offset, err := getLimitOffset(request)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	if getDeleted == "" {
		ids, err = secretStore.List(ctx, limit, offset)
	} else {
		ids, err = secretStore.ListDeleted(ctx, limit, offset)
	}
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(ids)
}

// @Summary Delete a secret by id
// @Description Soft delete secret by id. It can be recovered
// @Tags Secrets
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Secret Identifier"
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

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
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
// @Param storeName path string true "Secret Identifier"
// @Param id path string true "Key identifier"
// @Success 204 "Destroyed successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Secret not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/secrets/{id}/destroy [delete]
func (h *SecretsHandler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	id := mux.Vars(request)["id"]
	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
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
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Secret identifier"
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

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	secretStore, err := h.stores.GetSecretStore(ctx, StoreNameFromContext(ctx), userInfo)
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

func getLimitOffset(request *http.Request) (int, int, error) {
	limit := request.URL.Query().Get("limit")
	page := request.URL.Query().Get("page")
	if page != "" && limit == "" {
		limit = "100"
	}

	if limit == "" {
		return 0, 0, nil
	}

	iLimit, err := strconv.Atoi(limit)
	if err != nil {
		return 0, 0, errors.InvalidParameterError("invalid limit value")
	}

	iPage := 0
	iOffset := 0
	if page != "" {
		iPage, err = strconv.Atoi(page)
		if err != nil {
			return 0, 0, errors.InvalidParameterError("invalid page value")
		}

		iOffset = iPage * iLimit
	}

	return iLimit, iOffset, nil
}
