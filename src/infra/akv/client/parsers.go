package client

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
)

func ParseErrorResponse(err error) error {
	aerr, ok := err.(autorest.DetailedError)
	if !ok {
		return err
	}

	if rerr, ok := aerr.Original.(*azure.RequestError); ok && rerr.ServiceError.Code == "NotSupported" {
		return errors.NotSupportedError("%v", rerr)
	}

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
