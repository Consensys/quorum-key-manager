package jsonrpc

import "fmt"

func Error(err error) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32000,
		Message: err.Error(),
	}
}

func NotSupporteVersionError(version string) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32600,
		Message: fmt.Sprintf("JSON-RPC version %q not supported", version),
	}
}

func ParseError(err error) *ErrorMsg {
	return &ErrorMsg{
		Code:    -32700,
		Message: fmt.Sprintf("Could not parse JSON-RPC request: %v", err),
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
		Message: fmt.Sprintf("Invalid params: %v", err),
	}
}
