package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type key string

const (
	RequestHeaderKey key = "request-headers"
)

func GetRequest(ctx context.Context, client *http.Client, reqURL string) (*http.Response, error) {
	return request(ctx, client, reqURL, http.MethodGet, nil)
}

func DeleteRequest(ctx context.Context, client *http.Client, reqURL string) (*http.Response, error) {
	return request(ctx, client, reqURL, http.MethodDelete, nil)
}

func PostRequest(ctx context.Context, client *http.Client, reqURL string, postRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(postRequest)

	return request(ctx, client, reqURL, http.MethodPost, body)
}

func PatchRequest(ctx context.Context, client *http.Client, reqURL string, patchRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(patchRequest)

	return request(ctx, client, reqURL, http.MethodPatch, body)
}

func PutRequest(ctx context.Context, client *http.Client, reqURL string, putRequest interface{}) (*http.Response, error) {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(putRequest)

	return request(ctx, client, reqURL, http.MethodPut, body)
}

func CloseResponse(response *http.Response) {
	if deferErr := response.Body.Close(); deferErr != nil {
		_ = fmt.Errorf("could not close response body")
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
