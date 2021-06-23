package jsonrpc

import (
	"net/http"

	httpclient "github.com/consensysquorum/quorum-key-manager/pkg/http/client"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/request"
	"github.com/consensysquorum/quorum-key-manager/pkg/http/response"
)

// HTTPClient is a connector to a jsonrpc server
type HTTPClient struct {
	client httpclient.Client
}

// NewClient creates a new jsonrpc HTTPClient from an HTTP HTTPClient
func NewHTTPClient(c httpclient.Client) *HTTPClient {
	return &HTTPClient{
		client: c,
	}
}

// Do sends an jsonrpc request over the underlying HTTP client and returns a jsonrpc response
func (c *HTTPClient) Do(reqMsg *RequestMsg) (*ResponseMsg, error) {
	err := reqMsg.Validate()
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(reqMsg.Context(), http.MethodPost, "", nil)

	// write request body
	err = request.WriteJSON(req, reqMsg)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, DownstreamError(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, InvalidDownstreamHTTPStatuError(resp.StatusCode)
	}

	// Create response and reads body
	respMsg := new(ResponseMsg)
	err = response.ReadJSON(resp, respMsg)
	if err != nil {
		return nil, InvalidDownstreamResponse(err)
	}

	err = respMsg.Validate()
	if err != nil {
		return nil, InvalidDownstreamResponse(err)
	}

	return respMsg, nil
}
