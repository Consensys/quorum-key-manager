package handlers

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	jsonutils "github.com/consensys/quorum-key-manager/pkg/json"
	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/gorilla/mux"
)

type UtilsHandler struct {
	utils stores.Utilities
}

func NewUtilsHandler(utils stores.Utilities) *UtilsHandler {
	return &UtilsHandler{
		utils: utils,
	}
}

func (h *UtilsHandler) Register(r *mux.Router) {
	// Register utilities handler on /utilities
	utilsSubrouter := r.PathPrefix("/utilities").Subrouter()

	utilsSubrouter.Methods(http.MethodPost).Path("/keys/verify-signature").HandlerFunc(h.verifySignature)

	utilsSubrouter.Methods(http.MethodPost).Path("/ethereum/ec-recover").HandlerFunc(h.ecRecover)
	utilsSubrouter.Methods(http.MethodPost).Path("/ethereum/verify-message").HandlerFunc(h.verifyMessage)
	utilsSubrouter.Methods(http.MethodPost).Path("/ethereum/verify-typed-data").HandlerFunc(h.verifyTypedData)
}

// @Summary Verify key signature
// @Description Verify if signature data was signed by a specific key
// @Tags Utilities
// @Accept json
// @Produce json
// @Param request body types.VerifyKeySignatureRequest true "Verify signature request"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /keys/verify-signature [post]
func (h *UtilsHandler) verifySignature(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	verifyReq := &types.VerifyKeySignatureRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	err = h.utils.Verify(verifyReq.PublicKey, verifyReq.Data, verifyReq.Signature, &entities.Algorithm{
		Type:          entities.KeyType(verifyReq.SigningAlgorithm),
		EllipticCurve: entities.Curve(verifyReq.Curve),
	})
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary EC Recover
// @Description Recover an Ethereum sender from a signature of the format [R || S || V] where V is 0 or 1
// @Tags Utilities
// @Accept json
// @Produce plain
// @Param request body types.ECRecoverRequest true "Ethereum recover request"
// @Success 200 {string} string "Recovered sender address"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ethereum/ec-recover [post]
func (h *UtilsHandler) ecRecover(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	ecRecoverReq := &types.ECRecoverRequest{}
	err := jsonutils.UnmarshalBody(request.Body, ecRecoverReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	address, err := h.utils.ECRecover(ecRecoverReq.Data, ecRecoverReq.Signature)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(address.Hex()))
}

// @Summary Verify message signature (EIP-191)
// @Description Verify the signature of a message signed using standard format EIP-191
// @Tags Utilities
// @Accept json
// @Param request body types.VerifyRequest true "Ethereum signature verify request"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ethereum/verify-message [post]
func (h *UtilsHandler) verifyMessage(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	verifyReq := &types.VerifyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	err = h.utils.VerifyMessage(verifyReq.Address, verifyReq.Data, verifyReq.Signature)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Verify typed data signature (EIP-712)
// @Description Verify the signature of an Ethereum typed data using format defined at EIP-712
// @Tags Utilities
// @Accept json
// @Param request body types.VerifyTypedDataRequest true "Typed data request to verify"
// @Success 204 "Successful verification"
// @Failure 422 {object} ErrorResponse "Cannot verify signature"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /ethereum/verify-typed-data [post]
func (h *UtilsHandler) verifyTypedData(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	verifyReq := &types.VerifyTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyReq)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	typedData := formatters.FormatSignTypedDataRequest(&verifyReq.TypedData)
	err = h.utils.VerifyTypedData(verifyReq.Address, typedData, verifyReq.Signature)
	if err != nil {
		http2.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
