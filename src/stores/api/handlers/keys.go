package handlers

import (
	"encoding/base64"
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

type KeysHandler struct {
	stores storemanager.Manager
}

// New creates a http.Handler to be served on /keys
func NewKeysHandler(s storemanager.Manager) *KeysHandler {
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
}

func (h *KeysHandler) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	createKeyRequest := &types2.CreateKeyRequest{}
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
		&entities2.Algorithm{
			Type:          entities2.KeyType(createKeyRequest.SigningAlgorithm),
			EllipticCurve: entities2.Curve(createKeyRequest.Curve),
		},
		&entities2.Attributes{
			Tags: createKeyRequest.Tags,
		})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters2.FormatKeyResponse(key))
}

func (h *KeysHandler) importKey(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	importKeyRequest := &types2.ImportKeyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, importKeyRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	privKey, err := base64.URLEncoding.DecodeString(importKeyRequest.PrivateKey)
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
		privKey,
		&entities2.Algorithm{
			Type:          entities2.KeyType(importKeyRequest.SigningAlgorithm),
			EllipticCurve: entities2.Curve(importKeyRequest.Curve),
		},
		&entities2.Attributes{
			Tags: importKeyRequest.Tags,
		})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters2.FormatKeyResponse(key))
}

func (h *KeysHandler) sign(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	signPayloadRequest := &types2.SignBase64PayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPayloadRequest)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	data, err := base64.URLEncoding.DecodeString(signPayloadRequest.Data)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	keyStore, err := h.stores.GetKeyStore(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := keyStore.Sign(ctx, mux.Vars(request)["id"], data)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(base64.URLEncoding.EncodeToString(signature)))
}

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

	_ = json.NewEncoder(rw).Encode(formatters2.FormatKeyResponse(key))
}

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
