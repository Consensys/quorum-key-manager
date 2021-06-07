package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	jsonutils "github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/formatters"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	storesmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"github.com/gorilla/mux"
)

type Eth1Handler struct {
	stores storesmanager.Manager
}

// New creates a http.Handler to be served on /accounts
func NewAccountsHandler(s storesmanager.Manager) *Eth1Handler {
	return &Eth1Handler{
		stores: s,
	}
}

func (h *Eth1Handler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("").HandlerFunc(h.create)
	r.Methods(http.MethodPost).Path("/import").HandlerFunc(h.importAccount)
	r.Methods(http.MethodPost).Path("/{address}/sign").HandlerFunc(h.sign)
	r.Methods(http.MethodPost).Path("/{address}/sign-transaction").HandlerFunc(h.signTransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-quorum-private-transaction").HandlerFunc(h.signPrivateTransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-eea-transaction").HandlerFunc(h.signEEATransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-typed-data").HandlerFunc(h.signTypedData)
	r.Methods(http.MethodPost).Path("/{address}/restore").HandlerFunc(h.restore)
	r.Methods(http.MethodPost).Path("/ec-recover").HandlerFunc(h.ecRecover)
	r.Methods(http.MethodPost).Path("/verify-signature").HandlerFunc(h.verifySignature)
	r.Methods(http.MethodPost).Path("/verify-typed-data-signature").HandlerFunc(h.verifyTypedDataSignature)

	r.Methods(http.MethodPatch).Path("/{address}").HandlerFunc(h.update)

	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodGet).Path("/{address}").HandlerFunc(h.getOne)

	r.Methods(http.MethodDelete).Path("/{address}").HandlerFunc(h.delete)
	r.Methods(http.MethodDelete).Path("/{address}/destroy").HandlerFunc(h.destroy)
}

// @Summary Create ethereum account
// @Description Creates a new ECDSA Secp256k1 key representing an ethereum account
// @Accept  json
// @Produce  json
// @Param storeName path string true "Selected StoreID"
// @Param request body types.CreateEth1AccountRequest true "Create ethereum account request"
// @Success 200 {object} types.Eth1AccountResponse "Created ethereum account"
// @Router /stores/{storeName}/eth1 [post]
func (h *Eth1Handler) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	createReq := &types.CreateEth1AccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, createReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(request.Context()))
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

// @Summary Import ethereum account
// @Description Import an ECDSA Secp256k1 key representing an ethereum account
// @Accept  json
// @Produce  json
// @Param storeName path string true "Selected StoreID"
// @Param request body types.ImportEth1AccountRequest true "Create ethereum account request"
// @Success 200 {object} types.Eth1AccountResponse "Created ethereum account"
// @Router /stores/{storeName}/eth1/import [post]
func (h *Eth1Handler) importAccount(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	importReq := &types.ImportEth1AccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, importReq)
	if err != nil {
		fmt.Println(err)
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	eth1Acc, err := eth1Store.Import(ctx, importReq.ID, importReq.PrivateKey, &entities.Attributes{Tags: importReq.Tags})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEth1AccResponse(eth1Acc))
}

// @Summary Update ethereum account
// @Description Update ethereum account metadata
// @Accept  json
// @Produce  json
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.UpdateEth1AccountRequest true "Update ethereum account metadata request"
// @Success 200 {object} types.Eth1AccountResponse "Update ethereum account"
// @Router /stores/{storeName}/eth1/{address} [patch]
func (h *Eth1Handler) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	updateReq := &types.UpdateEth1AccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, updateReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
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

// @Summary Sign payload with ethereum account 
// @Description Sign random hex payload using selected ethereum account 
// @Accept json
// @Produce plain
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.SignHexPayloadRequest true "Sign payload request"
// @Success 200 {string} string "Signed payload data"
// @Router /stores/{storeName}/eth1/{address}/sign [post]
func (h *Eth1Handler) sign(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signPayloadReq := &types.SignHexPayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPayloadReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := eth1Store.Sign(ctx, getAddress(request), signPayloadReq.Data)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Sign typed data
// @Description Sign typed data, following the EIP-712 standard, using selected ethereum account 
// @Accept json
// @Produce plain
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.SignTypedDataRequest true "Sign typed data request"
// @Success 200 {string} string "Signed payload data"
// @Router /stores/{storeName}/eth1/{address}/sign-typed-data [post]
func (h *Eth1Handler) signTypedData(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signTypedDataReq := &types.SignTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signTypedDataReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
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

// @Summary Sign ethereum transaction
// @Description Sign ethereum transaction using selected ethereum account 
// @Accept json
// @Produce plain
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.SignETHTransactionRequest true "Sign ETH transaction request"
// @Success 200 {string} string "Signed transaction data"
// @Router /stores/{storeName}/eth1/{address}/sign-typed-data [post]
func (h *Eth1Handler) signTransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signTransactionReq := &types.SignETHTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signTransactionReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := eth1Store.SignTransaction(ctx, getAddress(request), signTransactionReq.ChainID.ToInt(), formatters.FormatTransaction(signTransactionReq))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Sign EEA transaction
// @Description Sign EEA transaction using selected ethereum account 
// @Accept json
// @Produce plain
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.SignEEATransactionRequest true "Sign EEA transaction request"
// @Success 200 {string} string "Signed EEA transaction data"
// @Router /stores/{storeName}/eth1/{address}/sign-eea-transaction [post]
func (h *Eth1Handler) signEEATransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signEEAReq := &types.SignEEATransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signEEAReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	tx, privateArgs := formatters.FormatEEATransaction(signEEAReq)
	signature, err := eth1Store.SignEEA(ctx, getAddress(request), signEEAReq.ChainID.ToInt(), tx, privateArgs)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Sign Quorum private transaction
// @Description Sign Quorum private transaction using selected ethereum account 
// @Accept json
// @Produce plain
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.SignQuorumPrivateTransactionRequest true "Sign Quorum transaction request"
// @Success 200 {string} string "Signed EEA transaction data"
// @Router /stores/{storeName}/eth1/{address}/sign-quorum-private-transaction [post]
func (h *Eth1Handler) signPrivateTransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signPrivateReq := &types.SignQuorumPrivateTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPrivateReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
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

// @Summary Get ethereum account
// @Description Fetch ethereum account information by address
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Success 200 {object} types.Eth1AccountResponse "Ethereum account object"
// @Router /stores/{storeName}/eth1/{address} [get]
func (h *Eth1Handler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
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

// @Summary List ethereum accounts
// @Description List addresses of ethereum account
// @Accept json
// @Produce json
// @Param storeName path string true "Selected StoreID"
// @Success 200 {array} []types.Eth1AccountResponse "Ethereum account list"
// @Router /stores/{storeName}/eth1 [get]
func (h *Eth1Handler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
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

// @Summary Delete ethereum account
// @Description Soft delete ethereum account, can be recovered
// @Accept json
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Success 200 {bool} bool
// @Router /stores/{storeName}/eth1/{address} [delete]
func (h *Eth1Handler) delete(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
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

// @Summary Destroy ethereum account
// @Description Hard delete ethereum account, cannot be recovered
// @Accept json
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Success 200 {bool} bool
// @Router /stores/{storeName}/eth1/{address}/destroy [delete]
func (h *Eth1Handler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
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

// @Summary Restore ethereum account
// @Description Recover a soft-deleted ethereum account
// @Accept json
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Success 200 {bool} bool
// @Router /stores/{storeName}/eth1/{address}/restore [post]
func (h *Eth1Handler) restore(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
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

// @Summary EC Recover
// @Description Recover ethereum transaction sender
// @Accept json
// @Produce plain
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.ECRecoverRequest true "Ethereum recover request"
// @Success 200 {string} string "Signed EEA transaction data"
// @Router /stores/{storeName}/eth1/ec-recover [post]
func (h *Eth1Handler) ecRecover(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	ecRecoverReq := &types.ECRecoverRequest{}
	err := jsonutils.UnmarshalBody(request.Body, ecRecoverReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	address, err := eth1Store.ECRevocer(ctx, ecRecoverReq.Data, ecRecoverReq.Signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(address))
}

// @Summary Verify signature
// @Description Verify signature of an ethereum signing
// @Accept json
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.VerifyEth1SignatureRequest true "Ethereum signature verify request"
// @Success 200 {string} string "Verification confirmed"
// @Failure 422 {string} string "Invalid verification"
// @Router /stores/{storeName}/eth1/verify-signature [post]
func (h *Eth1Handler) verifySignature(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	verifyReq := &types.VerifyEth1SignatureRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	err = eth1Store.Verify(ctx, verifyReq.Address.Hex(), verifyReq.Data, verifyReq.Signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Verify typed data signature
// @Description Verify signature of an ethereum type data signing
// @Accept json
// @Param storeName path string true "Selected StoreID"
// @Param address path string true "Ethereum address"
// @Param request body types.VerifyTypedDataRequest true "Ethereum signature verify request"
// @Success 200 {string} string "Verification confirmed"
// @Failure 422 {string} string "Invalid verification"
// @Router /stores/{storeName}/eth1/verify-signature [post]
func (h *Eth1Handler) verifyTypedDataSignature(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	verifyReq := &types.VerifyTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	typedData := formatters.FormatSignTypedDataRequest(&verifyReq.TypedData)
	err = eth1Store.VerifyTypedData(ctx, verifyReq.Address.Hex(), typedData, verifyReq.Signature)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getAddress(request *http.Request) string {
	return mux.Vars(request)["address"]
}
