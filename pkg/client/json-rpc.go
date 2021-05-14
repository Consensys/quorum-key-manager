package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
)

const jsonRPCPath = ""

type JsonRPCMessage struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
}

func (c *HTTPClient) Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error) {
	reqURL := fmt.Sprintf("%s/%s", c.config.URL, jsonRPCPath)
	req := &JsonRPCMessage{
		Method:  method,
		Version: "2.0",
	}
	if args != nil {
		var err error
		if req.Params, err = json.Marshal(args); err != nil {
			return nil, errors.EncodingError(err.Error())
		}
	}

	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)

	jsonRPCResp := &jsonrpc.ResponseMsg{}
	err = parseResponse(response, jsonRPCResp)
	if err != nil {
		return nil, err
	}

	return jsonRPCResp, nil
}
