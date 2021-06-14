// nolint
package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func parseResponse(response *http.Response, resp interface{}) error {
	if response.StatusCode == http.StatusAccepted || response.StatusCode == http.StatusOK {
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

	if response.StatusCode == http.StatusNotFound {
		return &ResponseError{
			StatusCode: response.StatusCode,
			Message:    string(respMsg),
		}
	}

	errResp := &ErrorResponse{}
	if err = json.Unmarshal(respMsg, &errResp); err == nil {
		return &ResponseError{
			StatusCode: response.StatusCode,
			Message:    errResp.Message,
			ErrorCode:  errResp.Code,
		}
	}

	return fmt.Errorf(string(respMsg))
}

func parseStringResponse(response *http.Response) (string, error) {
	if response.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			return "", err
		}

		return "", &ResponseError{
			StatusCode: response.StatusCode,
			Message:    errResp.Message,
			ErrorCode:  errResp.Code,
		}
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func parseEmptyBodyResponse(response *http.Response) error {
	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusAccepted {
		errResp := ErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			return err
		}

		return &ResponseError{
			StatusCode: response.StatusCode,
			Message:    errResp.Message,
			ErrorCode:  errResp.Code,
		}
	}

	return nil
}

func closeResponse(response *http.Response) {
	if deferErr := response.Body.Close(); deferErr != nil {
		return
	}
}
