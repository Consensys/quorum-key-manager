package client

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/hashicorp/vault/api"
)

func parseErrorResponse(err error) error {
	httpError, ok := err.(*api.ResponseError)
	if !ok {
		return errors.HashicorpVaultError("failed to connect to Hashicorp store")
	}

	switch httpError.StatusCode {
	case http.StatusNotFound:
		return errors.NotFoundError(httpError.Errors[0])
	case http.StatusBadRequest:
		return errors.InvalidFormatError(httpError.Errors[0])
	case http.StatusUnprocessableEntity:
		return errors.InvalidParameterError(httpError.Errors[0])
	case http.StatusConflict:
		return errors.AlreadyExistsError(httpError.Errors[0])
	default:
		return errors.HashicorpVaultError(httpError.Errors[0])
	}
}
