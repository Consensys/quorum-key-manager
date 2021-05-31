package handlers

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"

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
	router.Methods(http.MethodPost).Path("/import").HandlerFunc(h.importAccount)
	router.Methods(http.MethodPost).Path("/{address}/sign").HandlerFunc(h.sign)
	router.Methods(http.MethodPost).Path("/{address}/sign-transaction").HandlerFunc(h.signTransaction)
	router.Methods(http.MethodPost).Path("/{address}/sign-quorum-private-transaction").HandlerFunc(h.signPrivateTransaction)
	router.Methods(http.MethodPost).Path("/{address}/sign-eea-transaction").HandlerFunc(h.signEEATransaction)
	router.Methods(http.MethodPost).Path("/{address}/sign-typed-data").HandlerFunc(h.signTypedData)
	router.Methods(http.MethodPost).Path("/{address}/restore").HandlerFunc(h.restore)
	router.Methods(http.MethodPost).Path("/ec-revocer").HandlerFunc(h.ecRecover)
	router.Methods(http.MethodPost).Path("/verify-signature").HandlerFunc(h.verifySignature)
	router.Methods(http.MethodPost).Path("/verify-typed-data-signature").HandlerFunc(h.verifyTypedDataSignature)

	router.Methods(http.MethodPatch).Path("/{address}").HandlerFunc(h.update)

	router.Methods(http.MethodGet).Path("/").HandlerFunc(h.list)
	router.Methods(http.MethodGet).Path("/{address}").HandlerFunc(h.getOne)

	router.Methods(http.MethodDelete).Path("/{address}").HandlerFunc(h.delete)
	router.Methods(http.MethodDelete).Path("/{address}/destroy").HandlerFunc(h.destroy)

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

	eth1Acc, err := eth1Store.Create(ctx, createReq.ID, &entities.Attributes{Tags: createReq.Tags})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEth1AccResponse(eth1Acc))
}

func (h *Eth1Handler) importAccount(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	importReq := &types.ImportEth1AccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, importReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	privKey, err := hexutil.Decode(importReq.PrivateKey)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	eth1Acc, err := eth1Store.Import(ctx, importReq.ID, privKey, &entities.Attributes{Tags: importReq.Tags})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEth1AccResponse(eth1Acc))
}

func (h *Eth1Handler) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	updateReq := &types.UpdateEth1AccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, updateReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	eth1Acc, err := eth1Store.Update(ctx, getAddress(request), &entities.Attributes{Tags: updateReq.Tags})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEth1AccResponse(eth1Acc))
}

func (h *Eth1Handler) sign(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signPayloadReq := &types.SignHexPayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPayloadReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	data, err := hexutil.Decode(signPayloadReq.Data)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := eth1Store.Sign(ctx, getAddress(request), crypto.Keccak256(data))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

func (h *Eth1Handler) signTypedData(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signTypedDataReq := &types.SignTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signTypedDataReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	typedData := formatters.FormatSignTypedDataRequest(signTypedDataReq)
	signature, err := eth1Store.SignTypedData(ctx, getAddress(request), typedData)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

func (h *Eth1Handler) signTransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signTransactionReq := &types.SignETHTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signTransactionReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	chainID, _ := new(big.Int).SetString(signTransactionReq.ChainID, 10)
	signature, err := eth1Store.SignTransaction(ctx, getAddress(request), chainID, formatters.FormatTransaction(signTransactionReq))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

func (h *Eth1Handler) signEEATransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signEEAReq := &types.SignEEATransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signEEAReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	chainID, _ := new(big.Int).SetString(signEEAReq.ChainID, 10)
	tx, privateArgs := formatters.FormatEEATransaction(signEEAReq)
	signature, err := eth1Store.SignEEA(ctx, getAddress(request), chainID, tx, privateArgs)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

func (h *Eth1Handler) signPrivateTransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signPrivateReq := &types.SignQuorumPrivateTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPrivateReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := eth1Store.SignPrivate(ctx, getAddress(request), formatters.FormatPrivateTransaction(signPrivateReq))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

func (h *Eth1Handler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	getDeleted := request.URL.Query().Get("deleted")
	var eth1Acc *entities.ETH1Account
	if getDeleted == "" {
		eth1Acc, err = eth1Store.Get(ctx, getAddress(request))
	} else {
		eth1Acc, err = eth1Store.GetDeleted(ctx, getAddress(request))
	}
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEth1AccResponse(eth1Acc))
}

func (h *Eth1Handler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	getDeleted := request.URL.Query().Get("deleted")
	var addresses []string
	if getDeleted == "" {
		addresses, err = eth1Store.List(ctx)
	} else {
		addresses, err = eth1Store.ListDeleted(ctx)
	}
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(addresses)
}

func (h *Eth1Handler) delete(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = eth1Store.Delete(ctx, getAddress(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (h *Eth1Handler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = eth1Store.Destroy(ctx, getAddress(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (h *Eth1Handler) restore(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = eth1Store.Undelete(ctx, getAddress(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (h *Eth1Handler) ecRecover(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	ecRecoverReq := &types.ECRecoverRequest{}
	err := jsonutils.UnmarshalBody(request.Body, ecRecoverReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	data, err := hexutil.Decode(ecRecoverReq.Data)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	signature, err := hexutil.Decode(ecRecoverReq.Signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	address, err := eth1Store.ECRevocer(ctx, data, signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(address))
}

func (h *Eth1Handler) verifySignature(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	verifyReq := &types.VerifyEth1SignatureRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	data, err := hexutil.Decode(verifyReq.Data)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	signature, err := hexutil.Decode(verifyReq.Signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	err = eth1Store.Verify(ctx, verifyReq.Address, data, signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (h *Eth1Handler) verifyTypedDataSignature(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	verifyReq := &types.VerifyTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.backend.StoreManager().GetEth1Store(ctx, getStoreName(request))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := hexutil.Decode(verifyReq.Signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	typedData := formatters.FormatSignTypedDataRequest(&verifyReq.TypedData)
	err = eth1Store.VerifyTypedData(ctx, getAddress(request), typedData, signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getAddress(request *http.Request) string {
	return mux.Vars(request)["address"]
}
