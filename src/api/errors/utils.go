package errors

import (
	"encoding/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"net/http"
)

func WriteHTTPErrorResponse(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsAlreadyExistsError(err), errors.IsInvalidStateError(err):
		WriteErrorResponse(rw, http.StatusConflict, err)
	case errors.IsNotFoundError(err):
		WriteErrorResponse(rw, http.StatusNotFound, err)
	case errors.IsInvalidAuthenticationError(err), errors.IsUnauthorizedError(err):
		WriteErrorResponse(rw, http.StatusUnauthorized, err)
	case errors.IsInvalidFormatError(err):
		WriteErrorResponse(rw, http.StatusBadRequest, err)
	case errors.IsInvalidParameterError(err), errors.IsEncodingError(err):
		WriteErrorResponse(rw, http.StatusUnprocessableEntity, err)
	case err != nil:
		errMessage := "Internal server error. Please ask an admin for help or try again later"
		WriteErrorResponse(rw, http.StatusInternalServerError, errors.InternalError(errMessage))
	}
}

func WriteErrorResponse(rw http.ResponseWriter, status int, err error) {
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
