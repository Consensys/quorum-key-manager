package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

const (
	internalErrMsg    = "internal server error. Please ask an admin for help or try again later"
	internalDepErrMsg = "failed dependency. Please ask an admin for help or try again later"
)

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    string `json:"code,omitempty" example:"IR001"`
}

func WriteHTTPErrorResponse(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsAlreadyExistsError(err):
		writeErrorResponse(rw, http.StatusConflict, err)
	case errors.IsNotFoundError(err):
		writeErrorResponse(rw, http.StatusNotFound, err)
	case errors.IsUnauthorizedError(err):
		writeErrorResponse(rw, http.StatusUnauthorized, err)
	case errors.IsForbiddenError(err):
		writeErrorResponse(rw, http.StatusForbidden, err)
	case errors.IsInvalidFormatError(err):
		writeErrorResponse(rw, http.StatusBadRequest, err)
	case errors.IsInvalidParameterError(err), errors.IsEncodingError(err):
		writeErrorResponse(rw, http.StatusUnprocessableEntity, err)
	case errors.IsHashicorpVaultError(err), errors.IsAKVError(err), errors.IsDependencyFailureError(err), errors.IsAWSError(err), errors.IsPostgresError(err):
		writeErrorResponse(rw, http.StatusFailedDependency, errors.DependencyFailureError(internalDepErrMsg))
	case errors.IsNotImplementedError(err), errors.IsNotSupportedError(err):
		writeErrorResponse(rw, http.StatusNotImplemented, err)
	default:
		writeErrorResponse(rw, http.StatusInternalServerError, fmt.Errorf(internalErrMsg))
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
