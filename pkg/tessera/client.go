package tessera

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	httpclient "github.com/consensysquorum/quorum-key-manager/pkg/http/client"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/request"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/response"
)

//go:generate mockgen -source=client.go -destination=mock/client.go -package=mock

// Client is a client to Tessera Private Transaction Manager
type Client interface {
	StoreRaw(ctx context.Context, payload []byte, privateFrom string) ([]byte, error)
}

// HTTPClient is a tessera.Client that uses http
type HTTPClient struct {
	client httpclient.Client
}

// NewHTTPClient creates a new HTTPClient
func NewHTTPClient(c httpclient.Client) *HTTPClient {
	return &HTTPClient{
		client: c,
	}
}

type StoreRawRequest struct {
	Payload     string `json:"payload"`
	PrivateFrom string `json:"privateFrom"`
}

type StoreRawResponse struct {
	Key string `json:"key"`
}

func (c *HTTPClient) StoreRaw(ctx context.Context, payload []byte, privateFrom string) ([]byte, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/storeraw", nil)

	err := request.WriteJSON(req, &StoreRawRequest{
		Payload:     base64.StdEncoding.EncodeToString(payload),
		PrivateFrom: privateFrom,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	msg := new(StoreRawResponse)
	err = response.ReadJSON(resp, msg)
	if err != nil {
		return nil, err
	}

	b, err := base64.StdEncoding.DecodeString(msg.Key)
	if err != nil {
		return nil, err
	}

	return b, nil
}

var ErrNotConfigured = fmt.Errorf("tessera not configured")

// NotConfiguredClient is a Tessera Client that always return a tessera not configured error
type NotConfiguredClient struct{}

func (c *NotConfiguredClient) StoreRaw(context.Context, []byte, string) ([]byte, error) {
	return nil, ErrNotConfigured
}
