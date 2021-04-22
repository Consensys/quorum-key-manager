package client

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"net/http"
)

func ParseErrorResponse(err error) error {
	aerr, _ := err.(autorest.DetailedError)

	switch aerr.StatusCode.(int) {
	case http.StatusNotFound:
		return errors.NotFoundError("%v", aerr.Original)
	case http.StatusBadRequest:
		return errors.InvalidFormatError("%v", aerr.Original)
	case http.StatusUnprocessableEntity:
		return errors.InvalidParameterError("%v", aerr.Original)
	case http.StatusConflict:
		return errors.AlreadyExistsError("%v", aerr.Original)
	default:
		return errors.AKVConnectionError("%v", aerr.Original)
	}
}
