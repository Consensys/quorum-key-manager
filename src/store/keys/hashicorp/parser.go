package hashicorp

import (
	"encoding/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/hashicorp/vault/api"
	"net/http"
)

func parseResponse(data map[string]interface{}, resp interface{}) error {
	jsonbody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonbody, &resp); err != nil {
		return err
	}

	return nil
}

func parseErrorResponse(err error) error {
	httpError, _ := err.(*api.ResponseError)

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
