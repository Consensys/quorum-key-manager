package utils

import (
	"encoding/json"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/sdk/errors"
	"io/ioutil"
	"net/http"
)

const (
	internalErrMsg = "internal server error. Please ask an admin for help or try again later"
)

func ParseResponse(response *http.Response, resp interface{}) error {
	if response.StatusCode == http.StatusOK {
		if resp == nil {
			return nil
		}

		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			return err
		}

		return nil
	}

	// Read body
	respMsg, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if string(respMsg) != "" {
		errResp := errors.ErrorResponse{}
		if err = json.Unmarshal(respMsg, &errResp); err == nil {
			return fmt.Errorf("%v: %s", errResp.Code, errResp.Message)
		}
	}

	return parseResponseError(response.StatusCode, string(respMsg))
}

func parseResponseError(statusCode int, errMsg string) error {
	switch statusCode {
	case http.StatusBadRequest:
		if errMsg == "" {
			errMsg = "invalid request data"
		}
	case http.StatusConflict:
		if errMsg == "" {
			errMsg = "invalid data message"
		}
	case http.StatusNotFound:
		if errMsg == "" {
			errMsg = "cannot find entity"
		}
	case http.StatusUnauthorized:
		if errMsg == "" {
			errMsg = "not authorized"
		}
	case http.StatusUnprocessableEntity:
		if errMsg == "" {
			errMsg = "invalid request format"
		}
	default:
		if errMsg == "" {
			errMsg = internalErrMsg
		}
	}

	return fmt.Errorf("%v: %s", statusCode, errMsg)
}

func ParseStringResponse(response *http.Response) (string, error) {
	if response.StatusCode != http.StatusOK {
		errResp := errors.ErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			return "", err
		}

		return "", formattedError(errResp.Code, errResp.Message)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func ParseEmptyBodyResponse(response *http.Response) error {
	if response.StatusCode != http.StatusNoContent {
		errResp := errors.ErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			return err
		}

		return formattedError(errResp.Code, errResp.Message)
	}

	return nil
}

func formattedError(code uint64, msg string) error {
	return fmt.Errorf("%v: %s", code, msg)
}
