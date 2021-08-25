package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/gorilla/mux"
)

type KeysHandler struct {
	stores stores.Manager
}

// NewKeysHandler creates a http.Handler to be served on /keys
func NewKeysHandler(s stores.Manager) *KeysHandler {
	return &KeysHandler{
		stores: s,
	}
}

func (h *KeysHandler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/{id}/import").HandlerFunc(h.importKey)
	r.Methods(http.MethodPost).Path("/{id}/sign").HandlerFunc(h.sign)
	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.getOne)
	r.Methods(http.MethodPatch).Path("/{id}").HandlerFunc(h.update)
	r.Methods(http.MethodPut).Path("/{id}/restore").HandlerFunc(h.restore)
	r.Methods(http.MethodPost).Path("/verify-signature").HandlerFunc(h.verifySignature)
	r.Methods(http.MethodPost).Path("/{id}").HandlerFunc(h.create)

	r.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(h.delete)
	r.Methods(http.MethodDelete).Path("/{id}/destroy").HandlerFunc(h.destroy)
}

// @Summary Create key
// @Description Create a private Key using the specified Curve and Signing algorithm
// @Tags Keys
// @Accept json
// @Produce json
// @Param id path string true "Key Identifier"
// @Param storeName path string true "Store Identifier"
// @Param request body types.CreateKeyRequest true "Create key request"
// @Success 200 {object} types.KeyResponse "Key data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/{id} [post]
func (h *KeysHandler) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	createKeyRequest := &types.CreateKeyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, createKeyRequest)
	if err != nil && err.Error() != "EOF" {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userCtx := authenticator.UserContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userCtx.UserInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	key, err := keyStore.Create(
		ctx,
		getID(request),
		&entities.Algorithm{
			Type:          entities.KeyType(createKeyRequest.SigningAlgorithm),
			EllipticCurve: entities.Curve(createKeyRequest.Curve),
		},
		&entities.Attributes{
			Tags: createKeyRequest.Tags,
		})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatKeyResponse(key))
}

// @Summary Import Key
// @Description Import a private Key using the specified Curve and Signing algorithm
// @Tags Keys
// @Accept json
// @Produce json
// @Param id path string true "Key Identifier"
// @Param storeName path string true "Store Identifier"
// @Param request body types.ImportKeyRequest true "Create key request"
// @Success 200 {object} types.KeyResponse "Key data"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/{id}/import [post]
func (h *KeysHandler) importKey(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	importKeyRequest := &types.ImportKeyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, importKeyRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	key, err := keyStore.Import(
		ctx,
		getID(request),
		importKeyRequest.PrivateKey,
		&entities.Algorithm{
			Type:          entities.KeyType(importKeyRequest.SigningAlgorithm),
			EllipticCurve: entities.Curve(importKeyRequest.Curve),
		},
		&entities.Attributes{
			Tags: importKeyRequest.Tags,
		})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatKeyResponse(key))
}

// @Summary Sign random payload
// @Description Sign a random payload using the selected key
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Key identifier"
// @Param request body types.SignBase64PayloadRequest true "Signing request"
// @Success 200 {string} {string}"signature in base64"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store/Key not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/{id}/sign [post]
func (h *KeysHandler) sign(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	signPayloadRequest := &types.SignBase64PayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPayloadRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := keyStore.Sign(ctx, getID(request), signPayloadRequest.Data, nil)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(base64.StdEncoding.EncodeToString(signature)))
}

// @Summary Get key by ID
// @Description Retrieve a key object by identifier
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Key identifier"
// @Success 200 {object} types.KeyResponse "Key data"
// @Failure 404 {object} ErrorResponse "Store/Key not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/{id} [get]
func (h *KeysHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	getDeleted := request.URL.Query().Get("deleted")
	var key *entities.Key
	if getDeleted == "" {
		key, err = keyStore.Get(ctx, getID(request))
	} else {
		key, err = keyStore.GetDeleted(ctx, getID(request))
	}
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatKeyResponse(key))
}

// @Summary Update a key
// @Description Update the tags of a key by ID
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Key identifier"
// @Success 200 {object} types.KeyResponse "Key data"
// @Failure 404 {object} ErrorResponse "Store/Key not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/{id} [patch]
func (h *KeysHandler) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	updateRequest := &types.UpdateKeyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, updateRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	key, err := keyStore.Update(ctx, getID(request), &entities.Attributes{
		Tags: updateRequest.Tags,
	})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatKeyResponse(key))
}

// @Summary Restore a soft-deleted key
// @Description Restore a previously soft-deleted key by ID
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Key identifier"
// @Success 204 "Restored successfully"
// @Failure 404 {object} ErrorResponse "Store/Key not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/{id}/restore [put]
func (h *KeysHandler) restore(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = keyStore.Restore(ctx, getID(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary List Key ids
// @Description List identifiers of keys store on selected Store
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Success 200 {array} []types.KeyResponse "List of key ids"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys [get]
func (h *KeysHandler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	getDeleted := request.URL.Query().Get("deleted")
	var ids []string
	if getDeleted == "" {
		ids, err = keyStore.List(ctx)
	} else {
		ids, err = keyStore.ListDeleted(ctx)
	}
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(ids)
}

// @Summary Soft-delete Key
// @Description Delete a Key by ID. The key can be recovered
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Key identifier"
// @Success 204 "Deleted successfully"
// @Failure 404 {object} ErrorResponse "Store/Key not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/{id} [delete]
func (h *KeysHandler) delete(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = keyStore.Delete(ctx, getID(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Destroy a Key
// @Description Permanently delete a Key by ID
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Key identifier"
// @Success 204 "Destroyed successfully"
// @Failure 404 {object} ErrorResponse "Store/Key not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/{id}/destroy [delete]
func (h *KeysHandler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
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

// @Summary Verify key signature
// @Description Verify if signature data was signed by a specific key
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param id path string true "Key identifier"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/keys/verify-signature [post]
func (h *KeysHandler) verifySignature(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	verifyReq := &types.VerifyKeySignatureRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = keyStore.Verify(ctx, verifyReq.PublicKey, verifyReq.Data, verifyReq.Signature, &entities.Algorithm{
		Type:          entities.KeyType(verifyReq.SigningAlgorithm),
		EllipticCurve: entities.Curve(verifyReq.Curve),
	})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getID(request *http.Request) string {
	return mux.Vars(request)["id"]
}
