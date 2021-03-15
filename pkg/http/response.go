package http

import (
	"encoding/json"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
)

var internalErrMsg = "Internal server error. Please ask an admin for help or try again later"
var internalDepErrMsg = "Failed dependency. Please ask an admin for help or try again later"

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    uint64 `json:"code,omitempty" example:"24000"`
}

func WriteErrorResponse(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsNotFoundError(err):
		writeErrorResponse(rw, http.StatusNotFound, err)
	case errors.IsInvalidFormatError(err):
		writeErrorResponse(rw, http.StatusBadRequest, err)
	case errors.IsInvalidParameterError(err), errors.IsEncodingError(err):
		writeErrorResponse(rw, http.StatusUnprocessableEntity, err)
	case errors.IsHashicorpVaultConnectionError(err):
		writeErrorResponse(rw, http.StatusFailedDependency, errors.DependencyFailureError(internalDepErrMsg))
	case err != nil:
		writeErrorResponse(rw, http.StatusInternalServerError, errors.InternalError(internalErrMsg))
	}
}

func writeErrorResponse(rw http.ResponseWriter, status int, err error) {
	msg, e := json.Marshal(ErrorResponse{Message: err.Error(), Code: errors.FromError(err).GetCode()})
	if e != nil {
		http.Error(rw, e.Error(), status)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(status)
	_, _ = rw.Write(msg)
}

func WriteJSONResponse(rw http.ResponseWriter, resp interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(rw).Encode(resp)
}
