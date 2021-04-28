package client

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/hashicorp/vault/api"
)

func ParseErrorResponse(err error) error {
	httpError, ok := err.(*api.ResponseError)
	if !ok {
		return errors.HashicorpVaultConnectionError("failed to connect to Hashicorp store")
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
		return errors.HashicorpVaultConnectionError(httpError.Error())
	}
}
