package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"io"
	"io/ioutil"
	"net/http"
)

type key string

const (
	RequestHeaderKey key = "request-headers"
)

func getRequest(ctx context.Context, client *http.Client, reqURL string) (*http.Response, error) {
	return request(ctx, client, reqURL, http.MethodGet, nil)
}

func deleteRequest(ctx context.Context, client *http.Client, reqURL string) (*http.Response, error) {
	return request(ctx, client, reqURL, http.MethodDelete, nil)
}

func postRequest(ctx context.Context, client *http.Client, reqURL string, postRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(postRequest)

	return request(ctx, client, reqURL, http.MethodPost, body)
}

func patchRequest(ctx context.Context, client *http.Client, reqURL string, patchRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(patchRequest)

	return request(ctx, client, reqURL, http.MethodPatch, body)
}

func putRequest(ctx context.Context, client *http.Client, reqURL string, putRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(putRequest)

	return request(ctx, client, reqURL, http.MethodPut, body)
}

func closeResponse(response *http.Response) {
	if deferErr := response.Body.Close(); deferErr != nil {
		return
	}
}

func request(ctx context.Context, client *http.Client, reqURL, method string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(ctx, method, reqURL, body)
	if ctx.Value(RequestHeaderKey) != nil {
		for key, val := range ctx.Value(RequestHeaderKey).(map[string]string) {
			req.Header.Set(key, val)
		}
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func parseResponse(ctx context.Context, response *http.Response, resp interface{}) error {
	if response.StatusCode == http.StatusAccepted || response.StatusCode == http.StatusOK {
		if resp == nil {
			return nil
		}

		if err := json.NewDecoder(response.Body).Decode(resp); err != nil {
			log.FromContext(ctx).WithError(err).Error(invalidResponseBody)
			return errors.ServiceConnectionError(invalidResponseBody)
		}

		return nil
	}

	// Read body
	respMsg, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.FromContext(ctx).WithError(err).Error(cannotReadResponseBody)
		return errors.ServiceConnectionError(cannotReadResponseBody)
	}

	if string(respMsg) != "" {
		errResp := ErrorResponse{}
		if err = json.Unmarshal(respMsg, &errResp); err == nil {
			return errors.Errorf(errResp.Code, errResp.Message)
		}
	}

	return parseResponseError(response.StatusCode, string(respMsg))
}
