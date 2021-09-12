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
		return errors.NotFoundError(httpError.Error())
	case http.StatusBadRequest:
		return errors.InvalidFormatError(httpError.Error())
	case http.StatusUnprocessableEntity:
		return errors.InvalidParameterError(httpError.Error())
	case http.StatusConflict:
		return errors.AlreadyExistsError(httpError.Error())
	default:
		return errors.HashicorpVaultError(httpError.Error())
	}
}
