package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/consensysquorum/quorum-key-manager/pkg/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensysquorum/quorum-key-manager/pkg/json"
	"github.com/consensysquorum/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensysquorum/quorum-key-manager/src/stores/api/types"
	storesmanager "github.com/consensysquorum/quorum-key-manager/src/stores/manager"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
	"github.com/gorilla/mux"
)

const (
	QKMKeyIDPrefix = "qkm-"
)

type Eth1Handler struct {
	stores storesmanager.Manager
}

// NewAccountsHandler creates a http.Handler to be served on /accounts
func NewAccountsHandler(s storesmanager.Manager) *Eth1Handler {
	return &Eth1Handler{
		stores: s,
	}
}

func (h *Eth1Handler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("").HandlerFunc(h.create)
	r.Methods(http.MethodPost).Path("/import").HandlerFunc(h.importAccount)
	r.Methods(http.MethodPost).Path("/{address}/sign").HandlerFunc(h.sign)
	r.Methods(http.MethodPost).Path("/{address}/sign-data").HandlerFunc(h.signData)
	r.Methods(http.MethodPost).Path("/{address}/sign-transaction").HandlerFunc(h.signTransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-quorum-private-transaction").HandlerFunc(h.signPrivateTransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-eea-transaction").HandlerFunc(h.signEEATransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-typed-data").HandlerFunc(h.signTypedData)
	r.Methods(http.MethodPut).Path("/{address}/restore").HandlerFunc(h.restore)
	r.Methods(http.MethodPost).Path("/ec-recover").HandlerFunc(h.ecRecover)
	r.Methods(http.MethodPost).Path("/verify-signature").HandlerFunc(h.verifySignature)
	r.Methods(http.MethodPost).Path("/verify-typed-data-signature").HandlerFunc(h.verifyTypedDataSignature)

	r.Methods(http.MethodPatch).Path("/{address}").HandlerFunc(h.update)

	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodGet).Path("/{address}").HandlerFunc(h.getOne)

	r.Methods(http.MethodDelete).Path("/{address}").HandlerFunc(h.delete)
	r.Methods(http.MethodDelete).Path("/{address}/destroy").HandlerFunc(h.destroy)
}

// @Summary Create Ethereum Account
// @Description Creates a new ECDSA Secp256k1 key representing an Ethereum Account
// @Tags Ethereum Account
// @Accept  json
// @Produce  json
// @Param storeName path string true "Store Identifier"
// @Param request body types.CreateEth1AccountRequest true "Create Ethereum Account request"
// @Success 200 {object} types.Eth1AccountResponse "Created Ethereum Account"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

	var keyID string
	if createReq.KeyID != "" {
		keyID = createReq.KeyID
	} else {
		keyID = generateRandomKeyID()
	}

	eth1Acc, err := eth1Store.Create(ctx, keyID, &entities.Attributes{Tags: createReq.Tags})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEth1AccResponse(eth1Acc))
}

// @Summary Import Ethereum Account
// @Description Import an ECDSA Secp256k1 key representing an Ethereum account
// @Accept  json
// @Produce  json
// @Tags Ethereum Account
// @Param storeName path string true "Store Identifier"
// @Param request body types.ImportEth1AccountRequest true "Create Ethereum Account request"
// @Success 200 {object} types.Eth1AccountResponse "Created Ethereum Account"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

	var keyID string
	if importReq.KeyID != "" {
		keyID = importReq.KeyID
	} else {
		keyID = generateRandomKeyID()
	}

	eth1Acc, err := eth1Store.Import(ctx, keyID, importReq.PrivateKey, &entities.Attributes{Tags: importReq.Tags})
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEth1AccResponse(eth1Acc))
}

// @Summary Update Ethereum Account
// @Description Update Ethereum Account metadata
// @Accept  json
// @Produce  json
// @Tags Ethereum Account
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.UpdateEth1AccountRequest true "Update Ethereum Account metadata request"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Success 200 {object} types.Eth1AccountResponse "Update Ethereum Account"
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

// @Summary Sign payload with Ethereum Account
// @Description Sign random hexadecimal payload using selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignHexPayloadRequest true "Sign payload request"
// @Success 200 {string} string "Signed payload signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// @Summary Sign hashed payload with Ethereum Account
// @Description Sign Keccak256 payload using selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignHexPayloadRequest true "Sign payload request"
// @Success 200 {string} string "Signed payload signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/eth1/{address}/sign-data [post]
func (h *Eth1Handler) signData(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signPayloadReq := &types.SignHexPayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPayloadReq)
	if err != nil {
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	if len(signPayloadReq.Data) != 32 {
		errMsg := "expected data size of 32 bytes"
		WriteHTTPErrorResponse(rw, errors.InvalidFormatError(errMsg))
		return
	}

	eth1Store, err := h.stores.GetEth1Store(ctx, StoreNameFromContext(ctx))
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := eth1Store.SignData(ctx, getAddress(request), signPayloadReq.Data)
	if err != nil {
		WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Sign Typed Data
// @Description Sign Typed Data, following the EIP-712 Standard, using selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignTypedDataRequest true "Sign typed data request"
// @Success 200 {string} string "Signed typed data signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// @Summary Sign Ethereum transaction
// @Description Sign Ethereum transaction using selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignETHTransactionRequest true "Sign Ethereum transaction request"
// @Success 200 {string} string "Signed transaction signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/eth1/{address}/sign-transaction [post]
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
// @Description Sign EEA transaction using selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignEEATransactionRequest true "Sign EEA transaction request"
// @Success 200 {string} string "Signed EEA transaction signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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
// @Description Sign Quorum private transaction using selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignQuorumPrivateTransactionRequest true "Sign Quorum transaction request"
// @Success 200 {string} string "Signed Quorum private transaction signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// @Summary Get Ethereum Account
// @Description Fetch Ethereum Account data by its address
// @Tags Ethereum Account
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Success 200 {object} types.Eth1AccountResponse "Ethereum Account data"
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

// @Summary List Ethereum Accounts
// @Description List of addresses of stored Ethereum Accounts
// @Tags Ethereum Account
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Success 200 {array} []types.Eth1AccountResponse "Ethereum Account list"
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

// @Summary Delete Ethereum Account
// @Description Soft delete Ethereum Account, can be recovered
// @Tags Ethereum Account
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Success 204 "Deleted successfully"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// @Summary Destroy Ethereum Account
// @Description Hard delete Ethereum Account, cannot be recovered
// @Tags Ethereum Account
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Success 204 "Destroyed successfully"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// @Summary Restore Ethereum Account
// @Description Recover a soft-deleted Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Success 204 "Restored successfully"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/eth1/{address}/restore [put]
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
// @Description Recover Ethereum transaction sender
// @Tags Ethereum
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.ECRecoverRequest true "Ethereum recover request"
// @Success 200 {string} string "Recovered sender address"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
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
// @Description Verify signature of an Ethereum signature
// @Tags Ethereum
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.VerifyEth1SignatureRequest true "Ethereum signature verify request"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
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
// @Tags Ethereum
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.VerifyTypedDataRequest true "Ethereum signature verify request"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/eth1/verify-typed-data-signature [post]
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
	addr := ethcommon.HexToAddress(mux.Vars(request)["address"])
	return addr.Hex()
}

func generateRandomKeyID() string {
	return fmt.Sprintf("%s-%s", QKMKeyIDPrefix, common.RandString(15))
}
