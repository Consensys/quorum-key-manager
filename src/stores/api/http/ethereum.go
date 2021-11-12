package http

import (
	"encoding/json"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"net/http"

	"github.com/consensys/quorum-key-manager/src/auth/api/http_middlewares"

	"github.com/consensys/quorum-key-manager/pkg/common"
	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/gorilla/mux"
)

const (
	QKMKeyIDPrefix = "qkm-"
)

type EthHandler struct {
	stores stores.Stores
}

func NewEthHandler(storesConnector stores.Stores) *EthHandler {
	return &EthHandler{
		stores: storesConnector,
	}
}

func (h *EthHandler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("").HandlerFunc(h.create)
	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodPost).Path("/import").HandlerFunc(h.importAccount)
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

// @Summary Create an Ethereum Account
// @Description Create a new ECDSA Secp256k1 key representing an Ethereum Account
// @Tags Ethereum
// @Accept  json
// @Produce  json
// @Param storeName path string true "Store ID"
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

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(request.Context()), http_middlewares.UserInfoFromContext(ctx))
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

// @Summary Import an Ethereum Account
// @Description Import an ECDSA Secp256k1 key representing an Ethereum account
// @Accept  json
// @Produce  json
// @Tags Ethereum
// @Param storeName path string true "Store ID"
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

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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

// @Summary Update an Ethereum Account
// @Description Update an Ethereum Account metadata
// @Accept  json
// @Produce  json
// @Tags Ethereum
// @Param storeName path string true "Store ID"
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

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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

// @Summary Sign a message (EIP-191)
// @Description Sign a message, following EIP-191, using an existing Ethereum Account
// @Tags Ethereum
// @Accept json
// @Produce plain
// @Param storeName path string true "Store ID"
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

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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

// @Summary Sign Typed Data (EIP-712)
// @Description Sign Typed Data, following EIP-712, using identified Ethereum Account
// @Tags Ethereum
// @Accept json
// @Produce plain
// @Param storeName path string true "Store ID"
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

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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
// @Description Sign an Ethereum transaction using the identified Ethereum Account
// @Tags Ethereum
// @Accept json
// @Produce plain
// @Param storeName path string true "Store ID"
// @Param address path string true "Ethereum address"
// @Param request body types.SignETHTransactionRequest true "Sign Ethereum transaction request"
// @Success 200 {string} string "Signed raw transaction"
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

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	tx, err := formatters.FormatTransaction(signTransactionReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := ethStore.SignTransaction(ctx, getAddress(request), signTransactionReq.ChainID.ToInt(), tx)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(hexutil.Encode(signature)))
}

// @Summary Sign EEA transaction
// @Description Sign an EEA transaction using the identified Ethereum Account
// @Tags Ethereum
// @Accept json
// @Produce plain
// @Param storeName path string true "Store ID"
// @Param address path string true "Ethereum address"
// @Param request body types.SignEEATransactionRequest true "Sign EEA transaction request"
// @Success 200 {string} string "Signed raw EEA transaction"
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

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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
// @Description Sign a Quorum private transaction using the identified Ethereum Account
// @Tags Ethereum
// @Accept json
// @Produce plain
// @Param storeName path string true "Store ID"
// @Param address path string true "Ethereum address"
// @Param request body types.SignQuorumPrivateTransactionRequest true "Sign Quorum transaction request"
// @Success 200 {string} string "Signed raw Quorum private transaction"
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

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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

// @Summary Get an Ethereum Account
// @Description Fetch an Ethereum Account data by its address
// @Tags Ethereum
// @Accept json
// @Produce json
// @Param storeName path string true "Store ID"
// @Param address path string true "Ethereum address"
// @Param deleted query bool false "filter by only deleted accounts"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Success 200 {object} types.EthAccountResponse "Ethereum Account data"
// @Router /stores/{storeName}/ethereum/{address} [get]
func (h *EthHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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

// @Summary List Ethereum accounts
// @Description List Ethereum accounts addresses allocated in the targeted Store
// @Tags Ethereum
// @Accept json
// @Produce json
// @Param storeName path string true "Store ID"
// @Param deleted query bool false "filter by only deleted accounts"
// @Param chain_uuid query string false "Chain UUID"
// @Param limit query int false "page size"
// @Param page query int false "page number"
// @Success 200 {array} PageResponse "Ethereum Account list"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum [get]
func (h *EthHandler) list(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	limit, offset, err := getLimitOffset(request)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	getDeleted := request.URL.Query().Get("deleted")
	var addresses []ethcommon.Address
	if getDeleted == "" {
		addresses, err = ethStore.List(ctx, limit, offset)
	} else {
		addresses, err = ethStore.ListDeleted(ctx, limit, offset)
	}
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = http2.WritePagingResponse(rw, request, addresses)
}

// @Summary Delete Ethereum Account
// @Description Soft delete an Ethereum Account, can be recovered
// @Tags Ethereum
// @Accept json
// @Param storeName path string true "Store ID"
// @Param address path string true "Ethereum address"
// @Success 204 "Deleted successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address} [delete]
func (h *EthHandler) delete(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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
// @Tags Ethereum
// @Accept json
// @Param storeName path string true "Store ID"
// @Param address path string true "Ethereum address"
// @Success 204 "Destroyed successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/destroy [delete]
func (h *EthHandler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), http_middlewares.UserInfoFromContext(ctx))
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
// @Tags Ethereum
// @Accept json
// @Param storeName path string true "Store ID"
// @Param address path string true "Ethereum address"
// @Success 204 "Restored successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Store/Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stores/{storeName}/ethereum/{address}/restore [put]
func (h *EthHandler) restore(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userInfo := http_middlewares.UserInfoFromContext(ctx)
	ethStore, err := h.stores.Ethereum(ctx, StoreNameFromContext(ctx), userInfo)
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

func getAddress(request *http.Request) ethcommon.Address {
	return ethcommon.HexToAddress(mux.Vars(request)["address"])
}

func generateRandomKeyID() string {
	return fmt.Sprintf("%s%s", QKMKeyIDPrefix, common.RandString(15))
}
