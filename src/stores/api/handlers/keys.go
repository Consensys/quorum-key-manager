package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	jsonutils "github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/formatters"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	storesmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"

	"github.com/gorilla/mux"
)

type KeysHandler struct {
	stores storesmanager.Manager
}

// New creates a http.Handler to be served on /keys
func NewKeysHandler(s storesmanager.Manager) *KeysHandler {
	return &KeysHandler{
		stores: s,
	}
}

func (h *KeysHandler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("").HandlerFunc(h.create)
	r.Methods(http.MethodPost).Path("/import").HandlerFunc(h.importKey)
	r.Methods(http.MethodPost).Path("/{id}/sign").HandlerFunc(h.sign)
	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.getOne)
	r.Methods(http.MethodDelete).Path("/{id}").HandlerFunc(h.destroy)
	r.Methods(http.MethodPost).Path("/verify-signature").HandlerFunc(h.verifySignature)
}


// @Summary Create key
// @Description Create key pair by selected curve and signing algorithm  
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Param request body types.CreateKeyRequest true "Create key request"
// @Success 200 {object} types.KeyResponse "Key object"
// @Router /stores/{storeName}/keys [post]
func (h *KeysHandler) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	createKeyRequest := &types.CreateKeyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, createKeyRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	key, err := keyStore.Create(
		ctx,
		createKeyRequest.ID,
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

// @Summary Import key
// @Description Import key pair by selected curve and signing algorithm
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Param request body types.ImportKeyRequest true "Create key request"
// @Success 200 {object} types.KeyResponse "Key object"
// @Router /stores/{storeName}/keys/import [post]
func (h *KeysHandler) importKey(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	importKeyRequest := &types.ImportKeyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, importKeyRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	key, err := keyStore.Import(
		ctx,
		importKeyRequest.ID,
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

// @Summary Import key
// @Description Import key pair by selected curve and signing algorithm
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Param request body types.ImportKeyRequest true "Create key request"
// @Success 200 {object} types.KeyResponse "Key object"
// @Router /stores/{storeName}/keys/import [post]
func (h *KeysHandler) sign(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	signPayloadRequest := &types.SignBase64PayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPayloadRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := keyStore.Sign(ctx, mux.Vars(request)["id"], signPayloadRequest.Data)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(base64.URLEncoding.EncodeToString(signature)))
}

// @Summary Get key by ID
// @Description Retrieve key object by identifier
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Param id path string true "Retrieve keyID"
// @Success 200 {object} types.KeyResponse "Key object"
// @Router /stores/{storeName}/keys/{id} [get]
func (h *KeysHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	key, err := keyStore.Get(ctx, mux.Vars(request)["id"])
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatKeyResponse(key))
}

// @Summary List key ids
// @Description List store key ids
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Success 200 {array} []types.KeyResponse "List of key ids"
// @Router /stores/{storeName}/keys [get]
func (h *KeysHandler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	ids, err := keyStore.List(ctx)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(ids)
}

// @Summary Destroy key by ID
// @Description Hard delete selected key
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Param id path string true "KeyID to delete"
// @Success 200
// @Router /stores/{storeName}/keys/{id} [delete]
func (h *KeysHandler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = keyStore.Destroy(ctx, mux.Vars(request)["id"])
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Verify key signature
// @Description Verify whether a signature corresponds to an specific key
// @Tags Keys
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Param id path string true "Retrieve keyID"
// @Success 200 "Successful verification"
// @Success 422 "Invalid verification"
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

	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx))
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
