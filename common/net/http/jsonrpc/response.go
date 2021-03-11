package jsonrpc

import (
	"net/http"
)

type Response struct {
	*http.Response

	Msg *ResponseMsg
}


