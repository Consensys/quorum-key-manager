package client

import (
	"net/http"
)

type HTTPClient struct {
	client *http.Client
	config *Config
}

func NewHTTPClient(h *http.Client, c *Config) KeyManagerClient {
	return &HTTPClient{
		client: h,
		config: c,
	}
}
