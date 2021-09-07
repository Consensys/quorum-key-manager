package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/gorilla/mux"
)

const (
	QKMKeyIDPrefix = "qkm-"
)

type EthHandler struct {
	stores stores.Manager
}

// NewAccountsHandler creates a http.Handler to be served on /accounts
func NewAccountsHandler(s stores.Manager) *EthHandler {
	return &EthHandler{
		stores: s,
	}
}

func (h *EthHandler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("").HandlerFunc(h.create)
	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodPost).Path("/import").HandlerFunc(h.importAccount)
	r.Methods(http.MethodPost).Path("/ec-recover").HandlerFunc(h.ecRecover)
	r.Methods(http.MethodPost).Path("/verify").HandlerFunc(h.verify)
	r.Methods(http.MethodPost).Path("/verify-message").HandlerFunc(h.verifyMessage)
	r.Methods(http.MethodPost).Path("/verify-typed-data").HandlerFunc(h.verifyTypedData)
	r.Methods(http.MethodPost).Path("/{address}/sign-transaction").HandlerFunc(h.signTransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-quorum-private-transaction").HandlerFunc(h.signPrivateTransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-eea-transaction").HandlerFunc(h.signEEATransaction)
	r.Methods(http.MethodPost).Path("/{address}/sign-typed-data").HandlerFunc(h.signTypedData)
	r.Methods(http.MethodPost).Path("/{address}/sign-message").HandlerFunc(h.signMessage)
	r.Methods(http.MethodPut).Path("/{address}/restore").HandlerFunc(h.restore)
	r.Methods(http.MethodPatch).Path("/{address}").HandlerFunc(h.update)
	r.Methods(http.MethodGet).Path("/{address}").HandlerFunc(h.getOne)
	r.Methods(http.MethodDelete).Path("/{address}").HandlerFunc(h.delete)
	r.Methods(http.MethodDelete).Path("/{address}/destroy").HandlerFunc(h.destroy)
}

// @Summary Create Ethereum Account
// @Description Create a new ECDSA Secp256k1 key representing an Ethereum Account
// @Tags Ethereum Account
// @Accept  json
// @Produce  json
// @Param storeName path string true "Store Identifier"
// @Param request body types.CreateEthAccountRequest true "Create Ethereum Account request"
// @Success 200 {object} types.EthAccountResponse "Created Ethereum Account"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum [post]
func (h *EthHandler) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	createReq := &types.CreateEthAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, createReq)
	if err != nil && err.Error() != "EOF" {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(request.Context()), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	var keyID string
	if createReq.KeyID != "" {
		keyID = createReq.KeyID
	} else {
		keyID = generateRandomKeyID()
	}

	ethAcc, err := ethStore.Create(ctx, keyID, &entities.Attributes{Tags: createReq.Tags})
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEthAccResponse(ethAcc))
}

// @Summary Import Ethereum Account
// @Description Import an ECDSA Secp256k1 key representing an Ethereum account
// @Accept  json
// @Produce  json
// @Tags Ethereum Account
// @Param storeName path string true "Store Identifier"
// @Param request body types.ImportEthAccountRequest true "Create Ethereum Account request"
// @Success 200 {object} types.EthAccountResponse "Created Ethereum Account"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/import [post]
func (h *EthHandler) importAccount(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	importReq := &types.ImportEthAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, importReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	var keyID string
	if importReq.KeyID != "" {
		keyID = importReq.KeyID
	} else {
		keyID = generateRandomKeyID()
	}

	ethAcc, err := ethStore.Import(ctx, keyID, importReq.PrivateKey, &entities.Attributes{Tags: importReq.Tags})
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEthAccResponse(ethAcc))
}

// @Summary Update Ethereum Account
// @Description Update Ethereum Account metadata
// @Accept  json
// @Produce  json
// @Tags Ethereum Account
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.UpdateEthAccountRequest true "Update Ethereum Account metadata request"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Success 200 {object} types.EthAccountResponse "Update Ethereum Account"
// @Router /stores/{storeName}/ethereum/{address} [patch]
func (h *EthHandler) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	updateReq := &types.UpdateEthAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, updateReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	ethAcc, err := ethStore.Update(ctx, getAddress(request), &entities.Attributes{Tags: updateReq.Tags})
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEthAccResponse(ethAcc))
}

// @Summary Sign a message
// @Description Sign a message using an existing Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignMessageRequest true "Sign message request"
// @Success 200 {string} string "Signed payload signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/sign-message [post]
func (h *EthHandler) signMessage(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signPayloadReq := &types.SignMessageRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPayloadReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := ethStore.SignMessage(ctx, getAddress(request), signPayloadReq.Message)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
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
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/sign-typed-data [post]
func (h *EthHandler) signTypedData(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signTypedDataReq := &types.SignTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signTypedDataReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	typedData := formatters.FormatSignTypedDataRequest(signTypedDataReq)
	signature, err := ethStore.SignTypedData(ctx, getAddress(request), typedData)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Sign Ethereum transaction
// @Description Sign an Ethereum transaction using the selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignETHTransactionRequest true "Sign Ethereum transaction request"
// @Success 200 {string} string "Signed transaction signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/sign-transaction [post]
func (h *EthHandler) signTransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signTransactionReq := &types.SignETHTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signTransactionReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := ethStore.SignTransaction(ctx, getAddress(request), signTransactionReq.ChainID.ToInt(), formatters.FormatTransaction(signTransactionReq))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Sign EEA transaction
// @Description Sign an EEA transaction using the selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignEEATransactionRequest true "Sign EEA transaction request"
// @Success 200 {string} string "Signed EEA transaction signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/sign-eea-transaction [post]
func (h *EthHandler) signEEATransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signEEAReq := &types.SignEEATransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signEEAReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	tx, privateArgs := formatters.FormatEEATransaction(signEEAReq)
	signature, err := ethStore.SignEEA(ctx, getAddress(request), signEEAReq.ChainID.ToInt(), tx, privateArgs)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Sign Quorum private transaction
// @Description Sign a Quorum private transaction using the selected Ethereum Account
// @Tags Ethereum Account
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.SignQuorumPrivateTransactionRequest true "Sign Quorum transaction request"
// @Success 200 {string} string "Signed Quorum private transaction signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/sign-quorum-private-transaction [post]
func (h *EthHandler) signPrivateTransaction(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	signPrivateReq := &types.SignQuorumPrivateTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPrivateReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := ethStore.SignPrivate(ctx, getAddress(request), formatters.FormatPrivateTransaction(signPrivateReq))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Get Ethereum Account
// @Description Fetch an Ethereum Account data by its address
// @Tags Ethereum Account
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param deleted query bool false "filter by deleted accounts"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Success 200 {object} types.EthAccountResponse "Ethereum Account data"
// @Router /stores/{storeName}/ethereum/{address} [get]
func (h *EthHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	getDeleted := request.URL.Query().Get("deleted")
	var ethAcc *entities.ETHAccount
	if getDeleted == "" {
		ethAcc, err = ethStore.Get(ctx, getAddress(request))
	} else {
		ethAcc, err = ethStore.GetDeleted(ctx, getAddress(request))
	}
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatEthAccResponse(ethAcc))
}

// @Summary List Ethereum Accounts
// @Description List Ethereum Accounts located in the Store
// @Tags Ethereum Account
// @Accept json
// @Produce json
// @Param storeName path string true "Store Identifier"
// @Param deleted query bool false "filter by deleted accounts"
// @Param chain_uuid query string false "Chain UUID"
// @Success 200 {array} []types.EthAccountResponse "Ethereum Account list"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum [get]
func (h *EthHandler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	getDeleted := request.URL.Query().Get("deleted")
	var addresses []ethcommon.Address
	if getDeleted == "" {
		addresses, err = ethStore.List(ctx)
	} else {
		addresses, err = ethStore.ListDeleted(ctx)
	}
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(addresses)
}

// @Summary Delete Ethereum Account
// @Description Soft delete an Ethereum Account, can be recovered
// @Tags Ethereum Account
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Success 204 "Deleted successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address} [delete]
func (h *EthHandler) delete(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userCtx := authenticator.UserContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userCtx.UserInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = ethStore.Delete(ctx, getAddress(request))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Destroy Ethereum Account
// @Description Hard delete an Ethereum Account, cannot be recovered
// @Tags Ethereum Account
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Success 204 "Destroyed successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/destroy [delete]
func (h *EthHandler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = ethStore.Destroy(ctx, getAddress(request))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
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
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/restore [put]
func (h *EthHandler) restore(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = ethStore.Restore(ctx, getAddress(request))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary EC Recover
// @Description Recover an Ethereum transaction sender from a signature
// @Tags Ethereum Utils
// @Accept json
// @Produce plain
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.ECRecoverRequest true "Ethereum recover request"
// @Success 200 {string} string "Recovered sender address"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/ec-recover [post]
func (h *EthHandler) ecRecover(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	ecRecoverReq := &types.ECRecoverRequest{}
	err := jsonutils.UnmarshalBody(request.Body, ecRecoverReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	address, err := ethStore.ECRecover(ctx, ecRecoverReq.Data, ecRecoverReq.Signature)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(address.Hex()))
}

// @Summary Verify signature
// @Description Verify the signature of an Ethereum signature
// @Tags Ethereum Utils
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.VerifyRequest true "Ethereum signature verify request"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/verify [post]
func (h *EthHandler) verify(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	verifyReq := &types.VerifyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = ethStore.Verify(ctx, verifyReq.Address, verifyReq.Data, verifyReq.Signature)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Verify message signature
// @Description Verify the signature of a message
// @Tags Ethereum Utils
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.VerifyRequest true "Ethereum signature verify request"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/verify-message [post]
func (h *EthHandler) verifyMessage(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	verifyReq := &types.VerifyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = ethStore.VerifyMessage(ctx, verifyReq.Address, verifyReq.Data, verifyReq.Signature)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Verify typed data signature
// @Description Verify the signature of an Ethereum typed data signing
// @Tags Ethereum Utils
// @Accept json
// @Param storeName path string true "Store Identifier"
// @Param address path string true "Ethereum address"
// @Param request body types.VerifyTypedDataRequest true "Ethereum signature verify request"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/verify-typed-data [post]
func (h *EthHandler) verifyTypedData(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	verifyReq := &types.VerifyTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userInfo := authenticator.UserInfoContextFromContext(ctx)
	ethStore, err := h.stores.GetEthStore(ctx, StoreNameFromContext(ctx), userInfo)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	typedData := formatters.FormatSignTypedDataRequest(&verifyReq.TypedData)
	err = ethStore.VerifyTypedData(ctx, verifyReq.Address, typedData, verifyReq.Signature)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func getAddress(request *http.Request) ethcommon.Address {
	return ethcommon.HexToAddress(mux.Vars(request)["address"])
}

func generateRandomKeyID() string {
	return fmt.Sprintf("%s%s", QKMKeyIDPrefix, common.RandString(15))
}
