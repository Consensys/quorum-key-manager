package hashicorp

import (
	"encoding/json"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/common/errors"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

func parseResponse(data map[string]interface{}, resp interface{}) error {
	jsonbody, err := json.Marshal(data)
	if err != nil {
		errMessage := "failed to marshal response data"
		log.WithError(err).Error(errMessage)
		return errors.EncodingError(errMessage)
	}

	if err := json.Unmarshal(jsonbody, &resp); err != nil {
		errMessage := "failed to unmarshal response data"
		log.WithError(err).Error(errMessage)
		return errors.EncodingError(errMessage)
	}

	return nil
}

func parseErrorResponse(err error) error {
	httpError, ok := err.(*api.ResponseError)
	if !ok {
		errMessage := "failed to parse error response"
		log.WithError(err).Error(errMessage)
		return errors.HashicorpVaultConnectionError(errMessage)
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
