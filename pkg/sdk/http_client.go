package client

import (
	"net/http"
)

func NewHTTPClient(h *http.Client, c *Config) KeyManagerClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}

type HTTPClient struct {
	client *http.Client
	config *Config
}
