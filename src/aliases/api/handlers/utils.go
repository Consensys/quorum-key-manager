package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func jsonWrite(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8;")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	return json.NewEncoder(w).Encode(data)
}

func WriteHTTPErrorResponse(rw http.ResponseWriter, err error) {
	var writeErr error
	switch {
	case errors.IsAlreadyExistsError(err):
		writeErr = writeErrorResponse(rw, http.StatusConflict, err)
	case errors.IsNotFoundError(err):
		writeErr = writeErrorResponse(rw, http.StatusNotFound, err)
	case errors.IsUnauthorizedError(err):
		writeErr = writeErrorResponse(rw, http.StatusUnauthorized, err)
	case errors.IsInvalidFormatError(err):
		writeErr = writeErrorResponse(rw, http.StatusBadRequest, err)
	case errors.IsInvalidParameterError(err), errors.IsEncodingError(err):
		writeErr = writeErrorResponse(rw, http.StatusUnprocessableEntity, err)
	case errors.IsNotImplementedError(err), errors.IsNotSupportedError(err):
		writeErr = writeErrorResponse(rw, http.StatusNotImplemented, err)
	default:
		writeErr = writeErrorResponse(rw, http.StatusInternalServerError, fmt.Errorf(internalErrMsg))
	}
	if writeErr != nil {
		// TODO the: use logger
		log.Printf("error writing the original error: %v: %v", writeErr, err)
		http.Error(rw, writeErr.Error(), http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, status int, err error) error {
	msg, e := json.Marshal(ErrorResponse{Message: err.Error(), Code: errors.FromError(err).GetCode()})
	if e != nil {
		return e
	}

	// the: should we move that to a middleware?
	w.Header().Set("Content-Type", "application/json; charset=UTF-8;")
	// the: should we use that in every API response?
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, err = w.Write(msg)
	return err
}

const (
	internalErrMsg = "internal server error. Please ask an admin for help or try again later"
)

// ErrorResponse is the standard API error response.
// the: should we create a common lib? What format? Should every message have a potentially
// empty error info?
type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Code    string `json:"code,omitempty" example:"IR001"`
}
