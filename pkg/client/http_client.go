package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPClient struct {
	client *http.Client
	config *Config
}

var _ KeyManagerClient = &HTTPClient{}

func NewHTTPClient(h *http.Client, c *Config) *HTTPClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

func withURLStore(rootURL, storeID string) string {
	return fmt.Sprintf("%s/stores/%s", rootURL, storeID)
}

func listRequest(ctx context.Context, client *http.Client, urlPath string, deleted bool, limit, page uint64) ([]string, error) {
	reqURL, _ := url.Parse(urlPath)
	values := url.Values{}
	if deleted {
		values.Set("deleted", "true")
	}

	if limit != 0 {
		values.Set("limit", fmt.Sprintf("%d", limit))
	}
	if page != 0 {
		values.Set("page", fmt.Sprintf("%d", page))
	}

	reqURL.RawQuery = values.Encode()
	response, err := getRequest(ctx, client, reqURL.String())
	if err != nil {
		return nil, err
	}

	var pageRes pageStringResponse
	defer closeResponse(response)
	err = parseResponse(response, &pageRes)
	if err != nil {
		return nil, err
	}

	return pageRes.Data, nil
}
