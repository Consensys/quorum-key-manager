package client

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/consensys/quorum-key-manager/pkg/errors"
)

const (
	PurgeDeletedKeyMethod = "PurgeDeletedKey"
)

func parseErrorResponse(err error) error {
	aerr, ok := err.(autorest.DetailedError)
	if !ok {
		return errors.AKVError("%v", err)
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
	case http.StatusTooManyRequests:
		return errors.TooManyRequestError(aerr.Original.Error())
	case http.StatusConflict:
		if aerr.Method == PurgeDeletedKeyMethod {
			return errors.StatusConflictError(aerr.Original.Error())
		}
		return errors.AlreadyExistsError(aerr.Original.Error())
	default:
		return errors.AKVError(aerr.Original.Error())
	}
}
