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
		return errors.AKVConnectionError("%v", err)
	}

	if rerr, ok := aerr.Original.(*azure.RequestError); ok && rerr.ServiceError.Code == "NotSupported" {
		return errors.NotSupportedError("%v", rerr)
	}

	switch aerr.StatusCode.(int) {
	case http.StatusNotFound:
		return errors.NotFoundError(aerr.Original.Error())
	case http.StatusBadRequest:
		return errors.InvalidFormatError(aerr.Original.Error())
	case http.StatusUnprocessableEntity:
		return errors.InvalidParameterError(aerr.Original.Error())
	case http.StatusConflict:
		return errors.AlreadyExistsError(aerr.Original.Error())
	default:
		return errors.AKVConnectionError(aerr.Original.Error())
	}
}
