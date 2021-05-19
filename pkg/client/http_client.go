package client

import (
	"fmt"
	"net/http"
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
