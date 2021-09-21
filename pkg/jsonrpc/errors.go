package jsonrpc

import (
	"fmt"
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/http/proxy"
)

func Error(err error) *ErrorMsg {
	if errMsg, ok := err.(*ErrorMsg); ok {
		return errMsg
	}

	return InternalError(err)
}

func NotSupportedVersionError(version string) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32600,
		Message: fmt.Sprintf("JSON-RPC version %q not supported", version),
	}
}

func InvalidRequest(err error) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32600,
		Message: "Invalid Request",
		Data: map[string]interface{}{
			"message": err.Error(),
		},
	}
}

func ParseError(err error) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32700,
		Message: "Parse error",
		Data: map[string]interface{}{
			"message": err.Error(),
		},
	}
}

func InvalidMethodError(method string) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32601,
		Message: fmt.Sprintf("Invalid method %q", method),
	}
}

func NotImplementedMethodError(method string) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32601,
		Message: fmt.Sprintf("Method %q not implemented", method),
	}
}

func MethodNotFoundError() *ErrorMsg {
	return &ErrorMsg{
		Code:    -32601,
		Message: "Method not found",
	}
}

func InvalidParamsError(err error) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32602,
		Message: "Invalid params",
		Data: map[string]interface{}{
			"message": err.Error(),
		},
	}
}

func InternalError(err error) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32603,
		Message: "Internal error",
		Data: map[string]interface{}{
			"message": err.Error(),
		},
	}
}

func DownstreamError(err error) *ErrorMsg {
	if errMsg, ok := err.(*ErrorMsg); ok {
		return errMsg
	}

	code := proxy.StatusCodeFromRoundTripError(err)
	text := proxy.StatusText(code)

	return &ErrorMsg{
		Code:    -32000,
		Message: "Downstream error",
		Data: map[string]interface{}{
			"message": text,
			"status":  code,
		},
	}
}

func InvalidDownstreamHTTPStatusError(code int) *ErrorMsg {
	text := http.StatusText(code)
	return &ErrorMsg{
		Code:    -32001,
		Message: "Invalid downstream HTTP status",
		Data: map[string]interface{}{
			"message": text,
			"status":  code,
		},
	}
}

func InvalidDownstreamResponse(err error) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32003,
		Message: "Invalid downstream JSON-RPC response",
		Data: map[string]interface{}{
			"message": err.Error(),
		},
	}
}
