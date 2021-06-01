package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
)

const jsonRPCPath = "nodes"

type JSONRPCMessage struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
}

func (c *HTTPClient) Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error) {
	reqURL := fmt.Sprintf("%s/%s/%s", c.config.URL, jsonRPCPath, nodeID)
	req := &JSONRPCMessage{
		ID:      strconv.AppendUint(nil, uint64(10), 10),
		Method:  method,
		Version: "2.0",
	}
	if args != nil {
		var err error
		if req.Params, err = json.Marshal(args); err != nil {
			return nil, err
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
